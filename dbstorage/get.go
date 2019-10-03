package dbstorage

import (
	fs "github.com/gidyon/file-handlers"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"net/http"
)

func (fsDBH *fileDBHandler) getFile(w http.ResponseWriter, r *http.Request, key string) {
	var (
		data    string
		err     error
		caching = fsDBH.redisCaching && r.URL.Query().Get(urlQueryCacheKey) != ""
	)

	// get in cache first if enabled
	if caching {
		strCMD := fsDBH.redisClient.Get(key)
		data, err = strCMD.Result()
		if err != nil && err != redis.Nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// get data from database
	if data == "" || err == redis.Nil {
		file := &fs.FileData{}
		err := fsDBH.db.First(file, "id=?", key).Error
		if gorm.IsRecordNotFoundError(err) {
			http.Error(w, "FILE_NOT_FOUND", http.StatusNotFound)
			return
		}

		if caching {
			// set file data in cache only if its size is lower than size required
			if file.Size <= maxRedisFileSize {
				statusCMD := fsDBH.redisClient.Set(key, file.Data, 0)
				if err := statusCMD.Err(); err != nil {
					http.Error(w, "SET_FILE_IN_CACHE_FAILED", http.StatusNotFound)
					return
				}
			}
		}

		writeResponse(w, r, file.Data)
		return
	}

	writeResponse(w, r, []byte(data))
}
