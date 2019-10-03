// Package dbstorage is a file server handler for retrieving, uploading and storing files as blob in database.
// It uses SQL database for storing files and optional redis database for caching get requests.
package dbstorage

import (
	"crypto/sha256"
	"fmt"
	fs "github.com/gidyon/file-handlers"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"hash"
	"net/http"
	"path"
	"strconv"
	"strings"
)

var (
	urlQueryKeyOwnerID  = "oid"
	urlQueryKeyOwnerTag = "otag"
	urlQueryKeyFormFile = "file"
	urlQueryCacheKey    = "ch"

	maxUploadSize    int64 = 8 * 1024 * 1024 // ~ 8mb
	maxRedisFileSize int64 = 50 * 1024       // ~ 50kb
	redisCaching           = true
	useDB                  = true
)

// SetURLQueryKeyOwnerID sets the URL query key name for passing owner id
func SetURLQueryKeyOwnerID(key string) {
	urlQueryKeyOwnerID = key
}

// SetURLQueryKeyOwnerTag sets the URL query key name for passing owner tag
func SetURLQueryKeyOwnerTag(key string) {
	urlQueryKeyOwnerTag = key
}

// SetURLQueryKeyFormFile sets the URL query key name for passing form file
func SetURLQueryKeyFormFile(key string) {
	urlQueryKeyFormFile = key
}

// SetURLQueryCacheKey sets the URL query key name for passing cache option
func SetURLQueryCacheKey(key string) {
	urlQueryCacheKey = key
}

// SetMaxFileUploadSize sets the maximum upload size for files/files
func SetMaxFileUploadSize(size int) {
	if size < 0 {
		maxUploadSize = int64(size)
	}
}

// SetMaxRedisFileSize sets the maximum size of file that can be stored in redis. Files larger than the size specified will not be cached.
func SetMaxRedisFileSize(size int) {
	if size < 0 {
		maxRedisFileSize = int64(size)
	}
}

// DisableRedisCaching disable redis caching
func DisableRedisCaching() {
	redisCaching = false
}

type fileDBHandler struct {
	redisCaching  bool
	maxUploadSize int64
	redisClient   *redis.Client
	db            *gorm.DB
	hashFn        hash.Hash
}

// NewFileHandler creates a new file server that uses SQL database for storage of files and an optional redis database for caching. The database connection is mandatory for the handler to start. To disable redis caching, you can pass nil to redisClient argument or call API DisableRedisCaching function.
func NewFileHandler(db *gorm.DB, redisClient *redis.Client) (http.Handler, error) {
	if db == nil {
		return nil, errors.New("database connection is required")
	}

	// perform automigration
	db.AutoMigrate(&fs.FileData{}, &fs.FileInfo{})

	if redisClient == nil {
		redisCaching = false
	}

	// Hash function
	hashFn := sha256.New()

	return &fileDBHandler{
		redisCaching:  redisCaching,
		maxUploadSize: maxUploadSize,
		hashFn:        hashFn,
		redisClient:   redisClient,
		db:            db,
	}, nil
}

func (fsDBH *fileDBHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := path.Clean(r.URL.Path)
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	// get hash of file path
	key := fmt.Sprintf("%x", string(fsDBH.hashFn.Sum([]byte(upath))))

	switch r.Method {
	case http.MethodGet:
		fsDBH.getFile(w, r, key)
	case http.MethodPost:
		fsDBH.saveFile(w, r, key, upath)
	case http.MethodPut:
		fsDBH.saveFile(w, r, key, upath)
	case http.MethodDelete:
		fsDBH.deleteFile(w, r, key)
	}
}

// writeResponse write response headers and bytes
func writeResponse(w http.ResponseWriter, r *http.Request, data []byte) {
	// set headers
	w.Header().Set("Content-type", http.DetectContentType(data))
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(data)), 10))
	w.Header().Set("Accept-Ranges", "bytes")

	_, err := w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
