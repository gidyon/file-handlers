package file

import (
	"crypto/sha256"
	"fmt"
	"github.com/gidyon/file-handlers"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"hash"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	urlQueryKeyOwnerID   = "oid"
	urlQueryKeyOwnerTag  = "otag"
	urlQueryKeyDirectory = "dir"
	urlQueryKeyFormFile  = "file"

	maxUploadSize int64 = 8 * 1024 * 1024
	defaultDir          = "."
	useDB               = true
)

// SetURLQueryKeyOwnerID sets the URL query key for passing owner id
func SetURLQueryKeyOwnerID(key string) {
	urlQueryKeyOwnerID = key
}

// SetURLQueryKeyOwnerTag sets the URL query key for passing owner tag
func SetURLQueryKeyOwnerTag(key string) {
	urlQueryKeyOwnerTag = key
}

// SetURLQueryKeyDirectory sets the URL query key for passing directory to use
func SetURLQueryKeyDirectory(key string) {
	urlQueryKeyDirectory = key
}

// SetURLQueryKeyFormFile sets the URL query key for passing form file
func SetURLQueryKeyFormFile(key string) {
	urlQueryKeyFormFile = key
}

// SetMaxUploadSize sets the maximum upload size for files/files
func SetMaxUploadSize(size int) {
	if size < 0 {
		maxUploadSize = int64(size)
	}
}

// SetDefaultDir sets the uploads directory relative to the root
func SetDefaultDir(dir string) {
	defaultDir = filepath.Clean(dir)
}

// DisableDB disables the use of database to store image metadata
func DisableDB() {
	useDB = false
}

// ServerOptions contains options to setup and configure the file server
type ServerOptions struct {
	RootDir         string       // Root directory
	AllowedDirs     []string     // List of directories that is allowed access by server under root
	NotFoundHandler http.Handler // NotFound costom handler
	DB              *gorm.DB     // Database connection for storing file metadata
}

type fsHandler struct {
	root            string
	allowedDirs     []string
	defaultDir      string
	notFoundHandler http.Handler
	hashFn          hash.Hash
	db              *gorm.DB
	useDB           bool
	maxUploadSize   int64
}

// New creates a file server for the given root dir. It stores files metadata on the provided database connection.
// To disable storage of file metadata on the database, pass nil on opt.DB or call DisableDB prior to calling this functin.
// Once this function has been called, subsequent calls to API setters have no effect.
func New(opt *ServerOptions) (http.Handler, error) {
	if opt.RootDir == "" {
		// set root to current directory if its empty
		opt.RootDir = "."
	}

	// clean root
	opt.RootDir = filepath.Clean(opt.RootDir)

	defaultDir := filepath.Join(opt.RootDir, defaultDir)

	// allowed directories
	allowedDirs := make([]string, 0, len(opt.AllowedDirs)+1)
	allowedDirs = append(allowedDirs, defaultDir)
	for _, dir := range opt.AllowedDirs {
		allowedDirs = append(allowedDirs, filepath.Clean(dir))
	}

	testFile := filepath.Join(defaultDir, uuid.New().String()+".txt")

	// create simple file in default directory to check that it exitst
	f, err := os.Create(testFile)
	if err != nil {
		return nil, err
	}

	// delete the file later
	defer os.Remove(testFile)

	// close the file
	defer f.Close()

	if opt.NotFoundHandler == nil {
		// set not found to be http not found
		opt.NotFoundHandler = http.NotFoundHandler()
	}

	if opt.DB == nil {
		useDB = false
	}

	if useDB {
		// perform automigration
		opt.DB.AutoMigrate(&fs.FileData{}, &fs.FileInfo{})
	}

	// hash function
	hashFn := sha256.New()

	return &fsHandler{
		root:            opt.RootDir,
		allowedDirs:     allowedDirs,
		defaultDir:      defaultDir,
		notFoundHandler: opt.NotFoundHandler,
		hashFn:          hashFn,
		db:              opt.DB,
		useDB:           useDB,
		maxUploadSize:   maxUploadSize,
	}, nil
}

func (fsh *fsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// add prefix and clean
	upath := path.Clean(r.URL.Path)
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	// get hash of file path
	key := fmt.Sprintf("%x", string(fsh.hashFn.Sum([]byte(upath))))

	switch r.Method {
	case http.MethodGet:
		fsh.getFile(w, r, key, upath)
	case http.MethodPost:
		fsh.saveFile(w, r, key, upath)
	case http.MethodPut:
		fsh.saveFile(w, r, key, upath)
	case http.MethodDelete:
		fsh.deleteFile(w, r, key, upath)
	}
}

// isDirAllowed checks if a directory is present in the list of allowed directories
func (fsh *fsHandler) isDirAllowed(dir string) bool {
	for _, d := range fsh.allowedDirs {
		if d == dir {
			return true
		}
	}
	return false
}
