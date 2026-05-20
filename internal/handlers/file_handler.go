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

func (h *FileHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *FileHandler) GetById(c *gin.Context) {
	id := c.Param("id")

	file, err := h.service.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "File Not Found",
		})
		return
	}
	c.JSON(http.StatusOK, file)
}

func (h *FileHandler) Download(c *gin.Context) {
	id := c.Param("id")

	fullPath, mime, err := h.service.GetStoragePath(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Header("Content-Type", mime)
	c.File(fullPath)
}
