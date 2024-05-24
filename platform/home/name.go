package home

import (
	"context"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Name(auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm("id")
		name := c.PostForm("name")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.HTML(200, "error.tmpl", nil)
			return
		}

		filter := bson.M{"_id": objectID}
		update := bson.M{"$set": bson.M{"name": name}}

		collection := database.Collection("user")
		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.HTML(200, "error.tmpl", nil)
			return
		}

		c.Redirect(http.StatusFound, "/")
	}
}
