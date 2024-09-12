package main

import (
	"fileprocessor/pkg/api"
	"fileprocessor/pkg/psql"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database connection.
	if err := psql.InitializeDB(); err != nil {
		log.Fatalf("Could not initialize the database: %v", err)
	}
	r := gin.Default()

	api.SetupRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
