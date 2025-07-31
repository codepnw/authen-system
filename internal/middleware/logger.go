package middleware

import (
	"time"

	"github.com/codepnw/go-authen-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		duration := time.Since(time.Now())

		logger.InfoMiddleware("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", duration),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
