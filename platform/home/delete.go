package home

import (
	"context"
	"i9pay/db"
	"i9pay/platform/multipass"
	"i9pay/platform/pay"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72/sub"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Delete(auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, id, err := multipass.BothIDsFromCookie(c, auth, database)
		if err != nil {
			c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
			return
		}

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
			return
		}

		var user db.User
		collection := database.Collection("user")
		err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
		if err != nil {
			c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
			return
		}

		err = auth.DeleteUser(context.TODO(), user.Username)
		if err != nil {
			c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
			return
		}

		if user.Paying && user.Provider != "" && user.Provider != "Apple" && user.Provider != "Android" {

			subID, err := pay.UserIdToSubscriptionId(database, id)
			if err != nil {
				c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
				return
			}

			if _, err = sub.Cancel(subID, nil); err != nil {
				c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
				return
			}

			if err := pay.SetUserNotPaying(nil, auth, database, id, false); err != nil {
				c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
				return
			}
		}

		if _, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objectID}); err != nil {
			c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
			return
		}

		if _, err := database.Collection("usertoken").DeleteOne(context.TODO(), bson.M{"userid": id}); err != nil {
			c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
			return
		}

		for _, cookie := range c.Request.Cookies() {
			c.SetCookie(cookie.Name, "", -1, "/", "", false, true)
		}

		c.Redirect(http.StatusFound, "/login")

	}
}
