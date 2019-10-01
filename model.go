package fs

import (
	"time"
)

// Model contains time information
type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// FileMeta contains metadata information for a file
type FileMeta struct {
	ID       string `gorm:"primary_key"`
	OwnerID  string `gorm:"type:varchar(50)"`
	OwnerTag string `gorm:"type:varchar(20)"`
	Mime     string `gorm:"type:varchar(40)"`
	Name     string `gorm:"type:varchar(100)"`
	Path     string `gorm:"type:text"`
	Size     int64  `gorm:"type:int"`
}

// FileInfo model stores a file metadata
type FileInfo struct {
	FileMeta
	Model
}

// FileData model stores a file metadata and its content
type FileData struct {
	FileMeta
	Data []byte `gorm:"type:blob(8192000);not null"`
	Model
}
