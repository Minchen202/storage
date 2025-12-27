package api

import (
	"net/http"
	"p2p-storage/backend/internal/db"
	"p2p-storage/backend/pkg/models"

	"github.com/gin-gonic/gin"
)

type CreateFileRequest struct {
	Filename                string `json:"filename" binding:"required"`
	FileSize                int64  `json:"file_size" binding:"required"`
	ChunkCount              int    `json:"chunk_count" binding:"required"`
	EncryptionKeyEncrypted  string `json:"encryption_key_encrypted" binding:"required"`
}

func (env *Env) CreateFileHandler(c *gin.Context) {
	var req CreateFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	file := &models.File{
		UserID:                  userID.(string),
		Filename:                req.Filename,
		FileSize:                req.FileSize,
		ChunkCount:              req.ChunkCount,
		EncryptionKeyEncrypted:  req.EncryptionKeyEncrypted,
	}

	if err := db.CreateFile(env.DB, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file record"})
		return
	}

	c.JSON(http.StatusCreated, file)
}

func (env *Env) GetFilesHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	files, err := db.GetFilesByUser(env.DB, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve files"})
		return
	}

	c.JSON(http.StatusOK, files)
}
