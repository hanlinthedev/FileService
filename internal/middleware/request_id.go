package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequstId = "request_id"

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.NewString()
		c.Set(RequstId, requestId)
		c.Writer.Header().Set("X-Request-ID", requestId)
		c.Next()
	}
}
