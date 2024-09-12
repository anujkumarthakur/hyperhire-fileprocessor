package storage

import (
	"fileprocessor/pkg/models"
	"fileprocessor/pkg/psql"

	_ "github.com/lib/pq"
)

func GetFileMetadata(id string) (models.FileMetadata, error) {
	var file models.FileMetadata
	db := psql.GetDB()
	err := db.QueryRow("SELECT id, filename, size FROM files WHERE id = $1", id).
		Scan(&file.ID, &file.Filename, &file.Size)
	if err != nil {
		return models.FileMetadata{}, err
	}
	return file, nil
}

func DownloadFile(id string) ([]byte, error) {
	db := psql.GetDB()
	rows, err := db.Query("SELECT chunk_data FROM file_chunks WHERE file_id = $1 ORDER BY id", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fileData []byte
	for rows.Next() {
		var chunk []byte
		if err := rows.Scan(&chunk); err != nil {
			return nil, err
		}
		fileData = append(fileData, chunk...)
	}
	return fileData, nil
}
