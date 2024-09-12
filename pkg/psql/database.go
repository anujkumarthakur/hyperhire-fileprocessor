package psql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

// InitializeDB initializes the database connection using environment variables.
func InitializeDB() error {
	var err error
	// Load environment variables from .env file
	err = godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load .env file: %v", err)
	}
	// Print environment variables for debugging
	// fmt.Printf("DB_USER: %s\n", os.Getenv("DB_USER"))
	// fmt.Printf("DB_PASSWORD: %s\n", os.Getenv("DB_PASSWORD"))
	// fmt.Printf("DB_NAME: %s\n", os.Getenv("DB_NAME"))

	// Read database connection details from environment variables.
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	// Open a new database connection.
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %v", err)
	}

	// Ping the database to ensure the connection is established.
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	log.Println("Database connection established successfully")
	return nil
}

// GetDB returns the initialized database instance.
func GetDB() *sql.DB {
	return db
}
