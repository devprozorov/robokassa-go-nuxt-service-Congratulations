package mw

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleAny, _ := c.Get(CtxUserRole)
		role, _ := roleAny.(string)
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"ok": false, "error": "FORBIDDEN"})
			return
		}
		c.Next()
	}
}
