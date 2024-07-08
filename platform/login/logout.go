package login

import (
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func Logout(auth *auth.Client) gin.HandlerFunc {
	return CookieLogout
}

func CookieLogout(c *gin.Context) {
	CookieReset(c)

	c.Redirect(http.StatusFound, "/login")
}

func CookieReset(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	// Set the refreshToken cookie with a past expiration date
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "properLogin",
		Value:    "FALSE",
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})
}
