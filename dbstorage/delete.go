package dbstorage

import (
	fs "github.com/gidyon/file-handlers"
	"net/http"
)

func (fsDBH *fileDBHandler) deleteFile(w http.ResponseWriter, r *http.Request, key string) {
	ownerID := r.URL.Query().Get(urlQueryKeyOwnerID)

	// soft delete in mysql db
	err := fsDBH.db.Delete(&fs.FileData{}, "id=? AND owner_id=?", key, ownerID).Error
	if err != nil {
		http.Error(w, "DB_DELETE_FILE_FAILED", http.StatusNotFound)
		return
	}

	// delete in redis
	if fsDBH.redisCaching {
		boolCMD := fsDBH.redisClient.Del(key)
		if err := boolCMD.Err(); err != nil {
			http.Error(w, "CACHE_DELETE_FILE_FAILED", http.StatusNotFound)
			return
		}
	}

	w.Write([]byte("SUCCESS"))
}
