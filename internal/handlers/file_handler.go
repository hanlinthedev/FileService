package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanlinthedev/file-service/internal/response"
	"github.com/hanlinthedev/file-service/internal/services"
	"github.com/hanlinthedev/file-service/internal/transport"
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
		transport.HandleError(c, err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		transport.HandleError(c, err)
		return
	}

	defer func() {
		_ = file.Close()
	}()
	ctx := c.Request.Context()
	res, err := h.service.Upload(ctx, file, fileHeader)
	if err != nil {
		transport.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, res)
}

func (h *FileHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	if err := h.service.Delete(ctx, id); err != nil {
		transport.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusNoContent, nil)
}

func (h *FileHandler) GetById(c *gin.Context) {
	id := c.Param("id")

	file, err := h.service.GetById(id)
	if err != nil {
		transport.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, file)
}

func (h *FileHandler) Download(c *gin.Context) {
	id := c.Param("id")

	fullPath, mime, err := h.service.GetStoragePath(id)
	if err != nil {
		transport.HandleError(c, err)
		return
	}
	c.Header("Content-Type", mime)
	c.File(fullPath)
}
