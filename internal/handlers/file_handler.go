package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanlinthedev/file-service/internal/services"
)

type FileHandler struct {
	service services.FileService
}

func NewFileHandler(service services.FileService) *FileHandler {
	return &FileHandler{
		service: service,
	}
}

func (h *FileHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File is required.",
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open file.",
		})
		return
	}

	defer func() {
		_ = file.Close()
	}()

	res, err := h.service.Upload(file, fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, res)
}
