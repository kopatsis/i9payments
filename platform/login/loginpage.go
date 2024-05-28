package login

import (
	"context"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func LoginPage(auth *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := ExtractUIDFromSession(c, auth)
		if err != nil {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{})
			return
		}

		userRecord, err := auth.GetUser(context.Background(), uid)
		if err != nil {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{})
			return
		}

		email := userRecord.Email

		c.HTML(http.StatusOK, "login.tmpl", gin.H{"Email": email})
	}
}

func SignUpPage(auth *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := ExtractUIDFromSession(c, auth)
		if err != nil {
			c.HTML(http.StatusOK, "signup.tmpl", gin.H{})
			return
		}

		userRecord, err := auth.GetUser(context.Background(), uid)
		if err != nil {
			c.HTML(http.StatusOK, "signup.tmpl", gin.H{})
			return
		}

		email := userRecord.Email

		c.HTML(http.StatusOK, "signup.tmpl", gin.H{"Email": email})
	}
}
