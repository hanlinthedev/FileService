package transport

import (
	stdErrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"

	customerrors "github.com/hanlinthedev/file-service/internal/errors"
	"github.com/hanlinthedev/file-service/internal/response"
)

func HandleError(c *gin.Context, err error) {

	switch {

	case stdErrors.Is(err, customerrors.ErrEmptyFile):
		response.Error(c, http.StatusBadRequest, err.Error())

	case stdErrors.Is(err, customerrors.ErrFileTooLarge):
		response.Error(c, http.StatusBadRequest, err.Error())

	case stdErrors.Is(err, customerrors.ErrUnsupportedMime):
		response.Error(c, http.StatusBadRequest, err.Error())

	case stdErrors.Is(err, customerrors.ErrUnsupportedExt):
		response.Error(c, http.StatusBadRequest, err.Error())

	case stdErrors.Is(err, customerrors.ErrFileNotFound):
		response.Error(c, http.StatusNotFound, err.Error())

	default:
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}
