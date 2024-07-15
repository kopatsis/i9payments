package login

import (
	"context"
	"i9pay/db"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoginPage(auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		returnTo := c.Query("returnTo")

		uid, iat, err := ExtractUIDFromSession(c, auth)
		if err != nil {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{"Return": returnTo})
			return
		}

		userRecord, err := auth.GetUser(context.Background(), uid)
		if err != nil {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{"Return": returnTo})
			return
		}

		user, err := UserFromUID(uid, database)
		if err != nil {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{"Return": returnTo})
			return
		}

		issuedTime := time.Unix(iat, 0)
		resetTime := user.ResetDate.Time()

		if issuedTime.Before(resetTime) {
			CookieReset(c)
			c.HTML(http.StatusOK, "login.tmpl", gin.H{"Return": returnTo})
			return
		}

		email := userRecord.Email

		c.HTML(http.StatusOK, "login.tmpl", gin.H{"Email": email, "Return": returnTo})
	}
}

func SignUpPage(auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		returnTo := c.Query("returnTo")

		uid, iat, err := ExtractUIDFromSession(c, auth)
		if err != nil {
			c.HTML(http.StatusOK, "signup.tmpl", gin.H{"Return": returnTo})
			return
		}

		userRecord, err := auth.GetUser(context.Background(), uid)
		if err != nil {
			c.HTML(http.StatusOK, "signup.tmpl", gin.H{"Return": returnTo})
			return
		}

		user, err := UserFromUID(uid, database)
		if err != nil {
			c.HTML(http.StatusOK, "login.tmpl", gin.H{"Return": returnTo})
			return
		}

		issuedTime := time.Unix(iat, 0)
		resetTime := user.ResetDate.Time()

		if issuedTime.Before(resetTime) {
			CookieReset(c)
			c.HTML(http.StatusOK, "login.tmpl", gin.H{"Return": returnTo})
			return
		}

		email := userRecord.Email

		c.HTML(http.StatusOK, "signup.tmpl", gin.H{"Email": email, "Return": returnTo})
	}
}

func UserFromUID(sub string, database *mongo.Database) (*db.User, error) {
	collection := database.Collection("user")

	var user db.User

	if err := collection.FindOne(
		context.Background(),
		bson.M{"username": sub},
	).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
