package api

import (
	"database/sql"
	"fileprocessor/pkg/psql"
	"fileprocessor/pkg/storage"
	"fileprocessor/pkg/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var db *sql.DB

func UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}
	defer file.Close()

	// Get the filename from the header
	filename := header.Filename

	// Generate a unique fileID using UUID
	fileID := uuid.New().String()

	// Read file content into a byte slice
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Get the initialized database instance
	db := psql.GetDB()

	// Save the file metadata to the database
	err = utils.SaveFileMetadata(db, fileID, filename, len(fileBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
		return
	}

	// Split the file into chunks
	chunks := utils.SplitFile(fileBytes, 1024*1024) // Split into 1 MB chunks

	// Directory to save chunks
	chunkDir := "chunks"
	// Ensure the directory exists
	if err := os.MkdirAll(chunkDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory for chunks"})
		return
	}

	// Save each chunk to the database and file system
	for i, chunk := range chunks {
		// Save chunk to the database
		err = utils.SaveChunkToDatabase(db, fileID, chunk)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save chunk %d to database: %v", i+1, err)})
			return
		}

		// Save chunk to file system
		chunkFilename := filepath.Join(chunkDir, fmt.Sprintf("%s_part_%d", fileID, i+1))
		err = utils.SaveChunkToFile(chunk, chunkFilename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save chunk %d to file: %v", i+1, err)})
			return
		}
	}

	// Return the fileID to the client
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded and split successfully", "fileID": fileID})
}

func GetFileData(c *gin.Context) {
	id := c.Param("id")
	file, err := storage.GetFileMetadata(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.JSON(http.StatusOK, file)
}

func DownloadFile(c *gin.Context) {
	id := c.Param("id")
	fileData, err := storage.DownloadFile(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", fileData)
}
