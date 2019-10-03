package file

import (
	"github.com/Sirupsen/logrus"
	fs "github.com/gidyon/file-handlers"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (fsh *fsHandler) saveFile(w http.ResponseWriter, r *http.Request, key, path string) {
	// validate size
	r.Body = http.MaxBytesReader(w, r.Body, fsh.maxUploadSize)
	err := r.ParseMultipartForm(fsh.maxUploadSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get file from request
	file, header, err := r.FormFile(urlQueryKeyFormFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// read content
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// detect content-type
	ctype := http.DetectContentType(bs)

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

	var tx *gorm.DB
	// Create a transaction
	if fsh.useDB {
		tx = fsh.db.Begin()
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "RUNTIME_ERROR", http.StatusInternalServerError)
			}
		}()
		if tx.Error != nil {
			http.Error(w, "TRANSACTION_BEGIN_FAILED", http.StatusInternalServerError)
			return
		}

		// save file info metadata to database
		fileInfo := fs.FileInfo{
			FileMeta: fs.FileMeta{
				ID:       key,
				OwnerID:  r.URL.Query().Get(urlQueryKeyOwnerID),
				OwnerTag: r.URL.Query().Get(urlQueryKeyOwnerTag),
				Mime:     ctype,
				Size:     header.Size,
				Name:     fileName,
				Path:     path,
			},
			Model: fs.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		// save file info
		err = tx.Unscoped().Save(fileInfo).Error
		if err != nil {
			tx.Rollback()
			logrus.Errorln(err)
			http.Error(w, "TX_SAVE_FILE_FAILED", http.StatusInternalServerError)
			return
		}
	}

	// rollbacks the transaction
	rollBack := func() {
		if fsh.useDB {
			tx.Rollback()
		}
	}

	// if user can specify which directory to save file, use it
	dir := r.URL.Query().Get(urlQueryKeyDirectory)

	if dir != "" {
		dir = filepath.Clean(dir)
		if !fsh.isDirAllowed(dir) {
			rollBack()
			http.Error(w, "NOT_ALLOWED_ACESS_TO_DIRECTORY", http.StatusBadRequest)
			return
		}
	}

	// use default uploads when user has not specified which directory to save the file
	if dir == "" {
		dir = fsh.defaultDir
	}

	// path to file
	filePath := filepath.Join(dir, key)

	// create the file
	newFile, err := os.Create(filePath)
	if err != nil {
		rollBack()
		http.Error(w, "CANT_CREATE_FILE", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	_, err = newFile.Write(bs)
	if err != nil {
		rollBack()
		http.Error(w, "CANT_WRITE_TO_FILE", http.StatusInternalServerError)
		return
	}

	// commit file to database
	if fsh.useDB {
		err = tx.Commit().Error
		if err != nil {
			rollBack()
			http.Error(w, "FAILED_TO_COMMIT_FILE", http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusCreated)
	}

	w.Write([]byte("SUCCESS"))
}
