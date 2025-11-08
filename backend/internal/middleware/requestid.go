package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

// RequestID injeta um identificador único em cada requisição.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(requestIDHeader)
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Set(requestIDHeader, reqID)
		c.Writer.Header().Set(requestIDHeader, reqID)
		c.Next()
	}
}
