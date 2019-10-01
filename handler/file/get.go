package file

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (fsh *fsHandler) getFile(w http.ResponseWriter, r *http.Request, key, path string) {
	// if user has specified to get file from a given directory, use it
	dir := r.URL.Query().Get(urlQueryKeyDirectory)

	if dir != "" {
		dir = filepath.Clean(dir)
		if !fsh.isDirAllowed(dir) {
			http.Error(w, "NOT_ALLOWED_ACESS_TO_DIRECTORY", http.StatusBadRequest)
			return
		}
	}

	// use default uploads when user has not specified what directory to retrieve file
	if dir == "" {
		dir = fsh.defaultDir
	}

	filePath := filepath.Join(dir, key)

	// read file contents
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fsh.notFoundHandler.ServeHTTP(w, r)
			return
		}
	}

	// get file stats
	finfo, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// find mime of file
	ctype := http.DetectContentType(bs)

	w.Header().Set("Content-Type", ctype)
	w.Header().Set("Content-Length", strconv.FormatInt(finfo.Size(), 10))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Last-Modified", finfo.ModTime().UTC().Format(http.TimeFormat))

	_, err = w.Write(bs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
