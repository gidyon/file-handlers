package file

import (
	fs "github.com/gidyon/file-server"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"path/filepath"
)

func (fsh *fsHandler) deleteFile(w http.ResponseWriter, r *http.Request, key, path string) {
	var (
		ownerID = r.URL.Query().Get(urlQueryKeyOwnerID)
		tx      *gorm.DB
		err     error
	)

	// rollbacks a transaction
	rollBack := func() {
		if fsh.useDB {
			tx.Rollback()
		}
	}

	// Create a transaction
	if fsh.useDB {
		tx = fsh.db.Begin()
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "RUNTIME_ERROR", http.StatusInternalServerError)
				return
			}
		}()
		if tx.Error != nil {
			http.Error(w, "TRANSACTION_BEGIN_FAILED: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// delete file in database
		err := tx.Delete(&fs.FileInfo{}, "id=? AND owner_id=?", key, ownerID).Error
		if err != nil {
			rollBack()
			http.Error(w, "TX_DELETE_FILE_FAILED: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// if user can specify which directory to delete file from, use it
	dir := r.URL.Query().Get(urlQueryKeyDirectory)

	if dir != "" {
		dir = filepath.Clean(dir)
		if !fsh.isDirAllowed(dir) {
			rollBack()
			http.Error(w, "NOT_ALLOWED_DIRECTORY", http.StatusBadRequest)
			return
		}
	}

	// use default uploads directory when user has not specified which directory file resides
	if dir == "" {
		dir = fsh.defaultDir
	}

	// path to file
	filePath := filepath.Join(dir, key)

	// delete the file in directory
	err = os.Remove(filePath)
	if err != nil {
		rollBack()
		http.Error(w, "DELETE_FILE_FAILED: "+err.Error(), http.StatusBadRequest)
		return
	}

	// commit db transaction
	if fsh.useDB {
		err = tx.Commit().Error
		if err != nil {
			rollBack()
			http.Error(w, "FAILED_TO_COMMIT_FILE: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("SUCCESS"))
}
