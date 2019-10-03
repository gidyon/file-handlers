package dbstorage

import (
	fs "github.com/gidyon/file-handlers"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"mime"
	"net/http"
	"time"
)

func (fsDBH *fileDBHandler) saveFile(w http.ResponseWriter, r *http.Request, key, path string) {
	var (
		bs      = make([]byte, 0)
		err     error
		caching = fsDBH.redisCaching && r.URL.Query().Get(urlQueryCacheKey) != ""
	)

	// validate size
	r.Body = http.MaxBytesReader(w, r.Body, fsDBH.maxUploadSize)
	err = r.ParseMultipartForm(fsDBH.maxUploadSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get file content
	file, header, err := r.FormFile(urlQueryKeyFormFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// read file data
	bs, err = ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// content-type
	ctype := http.DetectContentType(bs)

	// file name
	fileName, err := func() (string, error) {
		if header.Filename != "" {
			return header.Filename, nil
		}
		fileEndings, err := mime.ExtensionsByType(ctype)
		if err != nil {
			return "", errors.New("CANT_READ_FILE_EXT_TYPE")
		}
		return uuid.New().String() + fileEndings[0], nil
	}()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get owner id
	ownerID := func() string {
		id := r.URL.Query().Get(urlQueryKeyOwnerID)
		if id != "" {
			return id
		}
		return "global"
	}()

	// Create file metadata together with its data
	fileData := fs.FileData{
		FileMeta: fs.FileMeta{
			ID:       key,
			OwnerID:  ownerID,
			OwnerTag: r.URL.Query().Get(urlQueryKeyOwnerTag),
			Mime:     ctype,
			Size:     header.Size,
			Name:     fileName,
			Path:     path,
		},
		Data: bs,
		Model: fs.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Save file in db
	err = fsDBH.db.Unscoped().Save(fileData).Error
	if err != nil {
		http.Error(w, "DB_SAVE_FILE_FAILED", http.StatusInternalServerError)
		return
	}

	msg := "SUCCESS"
	if caching {
		// set file data in cache only if its lower than size required
		if fileData.Size <= maxRedisFileSize {
			// Save file data in redis if the size does not extend limit
			statusCMD := fsDBH.redisClient.Set(key, fileData.Data, 0)
			if err := statusCMD.Err(); err != nil {
				http.Error(w, "CACHE_SAVE_FILE_FAILED", http.StatusNotFound)
				return
			}
		} else {
			msg += " :file too big to be saved in cache"
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(msg))
}
