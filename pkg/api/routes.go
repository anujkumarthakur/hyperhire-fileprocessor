package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/upload", UploadFile)
	router.GET("/files/:id", GetFileData)
	router.GET("/download/:id", DownloadFile)
}
