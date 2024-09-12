package utils

import (
	"database/sql"
	"fmt"
	"os"
)

// SaveFileMetadata saves file metadata to the database.
func SaveFileMetadata(db *sql.DB, fileID, filename string, size int) error {
	query := `INSERT INTO files (id, filename, size) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, fileID, filename, size)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

// SplitFile splits the file into chunks of specified size.
func SplitFile(fileBytes []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(fileBytes); i += chunkSize {
		end := i + chunkSize
		if end > len(fileBytes) {
			end = len(fileBytes)
		}
		chunks = append(chunks, fileBytes[i:end])
	}
	return chunks
}

// SaveChunkToFile saves a chunk of data to a file.
func SaveChunkToFile(chunk []byte, filePath string) error {
	err := os.WriteFile(filePath, chunk, 0644)
	if err != nil {
		return fmt.Errorf("failed to save chunk to file: %w", err)
	}
	return nil
}

// SaveChunkToDatabase saves a file chunk to the database.
func SaveChunkToDatabase(db *sql.DB, fileID string, chunk []byte) error {
	query := `INSERT INTO file_chunks (file_id, chunk_data) VALUES ($1, $2)`
	_, err := db.Exec(query, fileID, chunk)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}
