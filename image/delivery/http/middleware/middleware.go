package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware represent the data-struct for middleware
type Middleware struct {
	// another stuff , may be needed by middleware
}

func (m *Middleware) TransferEncodingCheck(c *gin.Context) {
	if c.Request.Header.Get("Content-Type") != "" && c.Request.Header.Get("Content-Type") != "application/json" {
		c.String(http.StatusUnsupportedMediaType, "Unsupported payload format")
		c.Abort()
		return
	}
}

// InitMiddleware initialize the middleware
func InitMiddleware() *Middleware {
	return &Middleware{}
}
