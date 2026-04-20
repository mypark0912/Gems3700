package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Session() gin.HandlerFunc {
	store := cookie.NewStore([]byte("nteksystem2026_idpm300"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		SameSite: 2, // Lax
	})
	return sessions.Sessions("gin_session", store)
}

// ClearOldCookies expires leftover cookies from the previous server (FastAPI "session", Express "connect.sid").
func ClearOldCookies() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, ck := range c.Request.Cookies() {
			if ck.Name == "gin_session" {
				continue
			}
			if ck.Name == "session" || strings.Contains(ck.Name, "connect") {
				http.SetCookie(c.Writer, &http.Cookie{
					Name:   ck.Name,
					Value:  "",
					Path:   "/",
					MaxAge: -1,
				})
			}
		}
		c.Next()
	}
}

// SessionRefresh extends session on API requests if user is logged in.
func SessionRefresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user != nil {
			session.Set("user", user)
			session.Set("userRole", session.Get("userRole"))
			session.Save()
		}
		c.Next()
	}
}
