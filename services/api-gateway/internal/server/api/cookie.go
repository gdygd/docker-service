package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, name string, value string, expTm time.Time) {
	maxAge := int(time.Until(expTm).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expTm,
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// func (server *Server) SetCookie(c *gin.Context, name string, value string, expTm time.Time) {
// 	server.SetCookieW(c.Writer, name, value, expTm)
// }

// SetCookieW gin.Context 없이 http.ResponseWriter로 직접 쿠키 설정
// ModifyResponse 등 gin.Context에 접근할 수 없는 환경에서 사용
func (server *Server) SetCookie(c *gin.Context, name string, value string, expTm time.Time) {
	maxAge := int(time.Until(expTm).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	isProd := os.Getenv("ENV") == "prod"

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expTm,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
	}

	domain := os.Getenv("COOKIE_DOMAIN")
	if domain != "" {
		cookie.Domain = domain
	}

	http.SetCookie(c.Writer, cookie)
}
