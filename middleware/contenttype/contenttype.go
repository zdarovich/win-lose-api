package contenttype

import (
	"github.com/gin-gonic/gin"
)
var isTypeAllowed = map[string]bool{"game":true, "server":true, "payment":true}

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

// middleware that checks for valid source-type header
func New() gin.HandlerFunc {
	return func(c *gin.Context) {
		typeVal := c.Request.Header.Get("Source-Type")

		if !isTypeAllowed[typeVal]{
			respondWithError(c, 401, "source type is not allowed")
			return
		}

		c.Next()
	}
}