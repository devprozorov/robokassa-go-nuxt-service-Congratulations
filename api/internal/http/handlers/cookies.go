package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) cookieDomain() string {
	// если хочешь чтобы работало и на поддоменах: ".happy4u.online"
	// если не нужно — оставь пусто ""
	if h.cfg.CookieDomain != "" {
		return h.cfg.CookieDomain
	}
	return ""
}

func (h *Handler) cookieSecure() bool {
	// в проде под https — true
	return h.cfg.AppEnv == "prod"
}

func (h *Handler) setAuthCookie(c *gin.Context, token string, ttl time.Duration) {
	secure := h.cookieSecure()
	domain := h.cookieDomain()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "happy_token",
		Value:    token,
		Path:     "/", // ВАЖНО: чтобы уходило на /api/payments/*
		Domain:   domain,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(ttl),
	})
}

func (h *Handler) clearAuthCookie(c *gin.Context) {
	secure := h.cookieSecure()
	domain := h.cookieDomain()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "happy_token",
		Value:    "",
		Path:     "/", // ВАЖНО: такой же Path как при установке
		Domain:   domain,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}
