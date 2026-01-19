package mw

import (
	"net/http"
	"strings"

	"happy-api/internal/config"
	"happy-api/internal/repo"
	"happy-api/internal/security"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CookieName  = "happy_token"
	CtxUserID   = "uid"
	CtxUserRole = "role"
)

func Auth(cfg config.Config, r *repo.Repo) gin.HandlerFunc {
	_ = r // repo сейчас не нужен, но сигнатуру оставляем как у тебя в роутере
	return func(c *gin.Context) {
		token := ""

		// 1) Authorization: Bearer ...
		auth := c.GetHeader("Authorization")
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			token = strings.TrimSpace(auth[7:])
		}

		// 2) Cookie
		if token == "" {
			if v, err := c.Cookie(CookieName); err == nil {
				token = v
			}
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "UNAUTHORIZED"})
			return
		}

		claims, err := security.ParseJWT(cfg.JWTSecret, token)
		if err != nil || claims == nil || claims.UserID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "UNAUTHORIZED"})
			return
		}

		uid, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "UNAUTHORIZED"})
			return
		}

		c.Set(CtxUserID, uid)
		c.Set(CtxUserRole, claims.Role)
		c.Next()
	}
}

// безопасно: никогда не паникует
func UserID(c *gin.Context) primitive.ObjectID {
	v, ok := c.Get(CtxUserID)
	if !ok || v == nil {
		return primitive.NilObjectID
	}
	uid, ok := v.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID
	}
	return uid
}

func Role(c *gin.Context) string {
	v, ok := c.Get(CtxUserRole)
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}
