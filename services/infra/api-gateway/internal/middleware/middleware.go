package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func CORS(origins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Idempotency-Key, X-Correlation-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID, X-Correlation-ID, X-RateLimit-*")
		c.Header("Access-Control-Max-Age", "86400")
		if c.Request.Method == "OPTIONS" { c.AbortWithStatus(204); return }
		c.Next()
	}
}

func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		cid := c.GetHeader("X-Correlation-ID")
		if cid == "" {
			b := make([]byte, 16)
			rand.Read(b)
			cid = hex.EncodeToString(b)
		}
		c.Set("correlation_id", cid)
		c.Header("X-Correlation-ID", cid)
		c.Next()
	}
}

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("%s %s %d %v [%s]", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start), c.GetString("correlation_id"))
	}
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(500, gin.H{"error": gin.H{"type": "INTERNAL_ERROR", "code": "PANIC", "message": "Internal server error"}})
	})
}
