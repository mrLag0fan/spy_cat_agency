package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a Gin middleware for logging HTTP requests and responses
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		method := c.Request.Method
		url := c.Request.URL.Path

		c.Next()

		statusCode := c.Writer.Status()

		duration := time.Since(startTime)

		log.Printf("Method: %s, URL: %s, Status: %d, Duration: %s\n", method, url, statusCode, duration)
	}
}
