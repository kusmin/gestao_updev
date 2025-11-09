package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger emite logs estruturados para cada requisição HTTP usando o zap.Logger fornecido.
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		proto := c.Request.Proto
		ip := c.ClientIP()
		ua := c.Request.UserAgent()
		ref := c.Request.Referer()
		bytesIn := c.Request.ContentLength
		if bytesIn < 0 {
			bytesIn = 0
		}

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		bytesOut := c.Writer.Size()
		if bytesOut < 0 {
			bytesOut = 0
		}

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("proto", proto),
			zap.String("client_ip", ip),
			zap.String("user_agent", ua),
			zap.Int("bytes_in", int(bytesIn)),
			zap.Int("bytes_out", bytesOut),
			zap.Duration("latency", latency),
		}

		if query != "" {
			fields = append(fields, zap.String("query", query))
		}

		if ref != "" {
			fields = append(fields, zap.String("referer", ref))
		}

		reqID := c.GetString(requestIDHeader)
		if reqID != "" {
			fields = append(fields, zap.String("request_id", reqID))
		}

		if tenantID, ok := c.Get(ContextTenantIDKey); ok {
			if tenantStr, okCast := tenantID.(string); okCast && tenantStr != "" {
				fields = append(fields, zap.String("tenant_id", tenantStr))
			}
		}

		if userID, ok := c.Get(ContextUserIDKey); ok {
			if userStr, okCast := userID.(string); okCast && userStr != "" {
				fields = append(fields, zap.String("user_id", userStr))
			}
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", strings.Join(c.Errors.Errors(), "; ")))
		}

		msg := "http_request"
		switch {
		case status >= 500:
			logger.Error(msg, fields...)
		case status >= 400:
			logger.Warn(msg, fields...)
		default:
			logger.Info(msg, fields...)
		}
	}
}
