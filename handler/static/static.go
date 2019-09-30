// Package static is a fast static file server with support for http2 server push.
// It caches static files in memory so that subsequent requests for the static file will retrieve the file from memory making it very fast than calling os primitives to open and read the files data.
// It can serve single page applications (SPAs) with an improved performance because of path caching on paths reesulting to 404.
// The API allows user to customize how the handler handles NotFound routes.
package static

import (
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// staticFile contains cached data for a static file to be used for writing to http response
type staticFile struct {
	data  []byte      // file data
	finfo os.FileInfo // file info
	ctype string      // file content type
}

// pushOptions contains server push information used for http2 server push
type pushOptions struct {
	path      string
	strict    bool
	pushFiles []string
}

func (po *pushOptions) Key() string {
	return po.path
}

// staticFileServer contains cached data and options to customize the static file server
type staticFileServer struct {
	rootDir         string
	pushSupport     bool
	pushOptions     *http.PushOptions
	mu              *sync.RWMutex // guards files
	files           map[string]*staticFile
	pushContent     map[string]*pushOptions
	notFoundPaths   map[string]int8
	notFoundHandler http.Handler
}

// NewHandler creates a static file server for the given rootDir directory.
// It caches the file in memory so that subsequent calls only write the files data to response
func NewHandler(
	rootDir string,
	URLPathPrefix string,
	notFoundHandler http.Handler,
	pushContent map[string][]string,
) (http.Handler, error) {
	if rootDir == "" {
		// set rootDir to current directory
		rootDir = "."
	}

	// clean rootDir
	rootDir = filepath.Clean(rootDir)

	// clean and update URLPathPrefix
	URLPathPrefix = "/" + strings.TrimPrefix(filepath.Clean(URLPathPrefix), "/")

	if notFoundHandler == nil {
		// set not found to be http default
		notFoundHandler = http.NotFoundHandler()
	}

	// create the server
	sfs := &staticFileServer{
		rootDir:       rootDir,
		mu:            &sync.RWMutex{},
		files:         make(map[string]*staticFile, 0),
		pushContent:   make(map[string]*pushOptions, len(pushContent)),
		notFoundPaths: make(map[string]int8, 0),
		pushOptions: &http.PushOptions{
			Method: http.MethodGet,
			Header: http.Header{
				"pushed-from": []string{"api"},
			},
		},
		pushSupport:     pushContent != nil && len(pushContent) > 0,
		notFoundHandler: notFoundHandler,
	}

	if sfs.pushSupport {
		for ppath, files := range pushContent {
			pushVal := &pushOptions{
				strict:    !strings.HasSuffix(ppath, "*"),
				path:      filepath.Clean(ppath), // removes any trailing /
				pushFiles: make([]string, 0),
			}

			for _, file := range files {
				filePath := filepath.Clean(file)
				pushVal.pushFiles = append(pushVal.pushFiles, filepath.Join(URLPathPrefix, filePath))

				// add the file to static files
				err := sfs.addStaticFile(filepath.Clean(filePath))
				if err != nil {
					return nil, err
				}
			}

			// update static file server pushContent map entry
			sfs.pushContent[pushVal.Key()] = pushVal
		}
	}

	return sfs, nil
}

func (sfs *staticFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Method must be get
	if r.Method != http.MethodGet {
		http.Error(w, "METHOD_NOT_ALLOWED", http.StatusBadRequest)
		return
	}

	fpath := path.Clean(r.URL.Path)
	if !strings.HasPrefix(fpath, "/") {
		// add prefix and clean
		fpath = "/" + fpath
		r.URL.Path = fpath
	}

	// update to render index page
	if fpath == "/" || fpath == "" {
		fpath = "/index.html"
	}

	// Check if file name exist in map
	_, ok := sfs.getStaticFile(fpath)
	if !ok {
		// pushes content to the client and serve index page
		pushAndServe := func() {
			sfs.serverPush(w, "/index.html")
			sfs.notFoundHandler.ServeHTTP(w, r)
		}

		// check if the path is in notFoundPaths so that we serve index page
		if _, ok := sfs.notFoundPaths[fpath]; ok {
			pushAndServe()
			return
		}

		err := sfs.addStaticFile(fpath)
		if os.IsNotExist(err) {
			// we may have a data race but its fine :)
			sfs.addNotFoundPath(fpath)
			pushAndServe()
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// push content to client
	sfs.serverPush(w, fpath)

	// write file data to response
	sfs.writeResponse(w, r, fpath)
}

// serverPush pushes content to the client
func (sfs *staticFileServer) serverPush(w http.ResponseWriter, fpath string) {
	// return early when push support is not enabled
	if !sfs.pushSupport {
		return
	}

	pusher, ok := w.(http.Pusher)
	if !ok {
		return
	}

	// push content to the client
	for key, pushInfo := range sfs.pushContent {
		if key == fpath {
			for _, target := range pushInfo.pushFiles {
				pusher.Push(target, sfs.pushOptions)
			}
			break
		}
	}
}

func (sfs *staticFileServer) writeResponse(w http.ResponseWriter, r *http.Request, name string) {
	// get the static file
	sfile, ok := sfs.getStaticFile(name)
	if !ok {
		http.Error(w, "file data does not exist", http.StatusInternalServerError)
		return
	}

	// set headers
	w.Header().Set("Content-Type", sfile.ctype)
	w.Header().Set("Content-Length", strconv.FormatInt(sfile.finfo.Size(), 10))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Last-Modified", sfile.finfo.ModTime().UTC().Format(http.TimeFormat))

	_, err := w.Write(sfile.data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

// getStaticFile retrieves the static file in a concurrent safe manner
func (sfs *staticFileServer) getStaticFile(fpath string) (*staticFile, bool) {
	sfs.mu.RLock()
	sf, ok := sfs.files[fpath]
	sfs.mu.RUnlock()

	return sf, ok
}

// addNotFoundPath adds the path to list of notFoundPaths; not concurrent safe;
func (sfs *staticFileServer) addNotFoundPath(fpath string) {
	sfs.notFoundPaths[fpath] = 2
}

// addStaticFile adds the static file to the map for faster subsequent retrievals on similar path
func (sfs *staticFileServer) addStaticFile(path string) error {
	filePath := filepath.Join(sfs.rootDir, path)
	// read file content
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// get file stats
	finfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	// find mime of file
	ctype := mime.TypeByExtension(filepath.Ext(filePath))
	if ctype == "" {
		ctype = http.DetectContentType(bs)
	}

	sfs.mu.Lock()
	// add the static files map entry without any data races
	sfs.files[path] = &staticFile{
		data:  bs,
		finfo: finfo,
		ctype: ctype,
	}
	sfs.mu.Unlock()

	return nil
}
