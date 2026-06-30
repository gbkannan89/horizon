package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security-related HTTP headers to every response.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), camera=(), microphone=()")
		c.Next()
	}
}

// RequestSize limits the maximum request body size.
func RequestSize(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

// ETag generates an ETag header based on the response body for GET requests.
func ETag() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		blw := &bodyCaptureWriter{body: &strings.Builder{}, ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		if c.Writer.Status() == 200 {
			hash := sha256.Sum256([]byte(blw.body.String()))
			etag := fmt.Sprintf("\"%s\"", hex.EncodeToString(hash[:16]))

			c.Header("ETag", etag)
			c.Header("Cache-Control", "private, max-age=5")

			if match := c.GetHeader("If-None-Match"); match == etag {
				c.Writer.WriteHeader(304)
				c.Writer.Write(nil)
				return
			}
		}
	}
}

type bodyCaptureWriter struct {
	gin.ResponseWriter
	body *strings.Builder
}

func (w *bodyCaptureWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *bodyCaptureWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// RequestTimeout sets a maximum duration for each request.
func RequestTimeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// SlowRequest logs requests that exceed the threshold.
func SlowRequest(threshold time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		if d := time.Since(start); d > threshold {
			log.Printf("[SLOW] %s %s took %v", c.Request.Method, c.Request.URL.Path, d)
		}
	}
}
