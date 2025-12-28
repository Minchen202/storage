package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (env *Env) UploadChunkHandler(c *gin.Context) {
	chunkHash := c.Param("hash")
	if chunkHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chunk hash is required"})
		return
	}

	err := env.R2Client.UploadChunk(c.Request.Context(), chunkHash, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload chunk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chunk uploaded successfully"})
}

func (env *Env) DownloadChunkHandler(c *gin.Context) {
	chunkHash := c.Param("hash")
	if chunkHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chunk hash is required"})
		return
	}

	chunkData, err := env.R2Client.DownloadChunk(c.Request.Context(), chunkHash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chunk not found"})
		return
	}
	defer chunkData.Close()

	c.DataFromReader(http.StatusOK, c.Request.ContentLength, "application/octet-stream", chunkData, nil)
}
