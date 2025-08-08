package models

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	constants "github.com/ntu-onemdp/onemdp-backend/config"
)

type FileFactory struct {
}

func NewFileFactory() *FileFactory {
	return &FileFactory{}
}

// FileMetadata models information that is returned to the frontend.
type FileMetadata struct {
	DbFile

	Author string `json:"author" db:"author_name"` // Name of the author after performing join with users table
}

// DbFile models how file metadata is stored in postgres.
type DbFile struct {
	FileId      string    `json:"file_id" db:"file_id"`
	AuthorUid   string    `json:"author_uid" db:"author"` // UID of author
	Filename    string    `json:"filename" db:"filename"`
	GCSFilename string    `json:"gcs_filename" db:"gcs_filename"`
	Status      string    `json:"status" db:"status"`
	TimeCreated time.Time `json:"time_created" db:"time_created"`

	TimeDeleted *time.Time `json:"time_deleted" db:"time_deleted"`
	DeletedBy   *string    `json:"deleted_by" db:"deleted_by"`
	FileGroup   *string    `json:"file_group" db:"file_group"` // Arbitary user defined group. (e.g. Week 1, Week 2, etc.)
}

func (f *FileFactory) New(author string, filename string, filegroup *string) *DbFile {
	return &DbFile{
		FileId:      "f" + gonanoid.Must(constants.CONTENT_ID_LENGTH),
		AuthorUid:   author,
		Filename:    filename,
		GCSFilename: time.Now().Format("060102150405") + filename,
		Status:      "available",
		TimeCreated: time.Now(),
		FileGroup:   filegroup,
	}
}
