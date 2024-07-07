package login

import (
	"context"
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ResetPasswordDate(authClient *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(400, gin.H{"Error": "No valid token in request"})
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			c.JSON(400, gin.H{"Error": "No valid token in request"})
			return
		}

		idToken := splitToken[1]

		fullToken, err := authClient.VerifyIDToken(context.TODO(), idToken)
		if err != nil {
			c.JSON(400, gin.H{"Error": "No valid token in request"})
			return
		}

		collection := database.Collection("user")

		filter := bson.M{"username": fullToken.UID}
		update := bson.M{
			"$set": bson.M{
				"reset": primitive.NewDateTimeFromTime(time.Now()),
			},
		}

		if _, err := collection.UpdateOne(context.TODO(), filter, update); err != nil {
			c.JSON(400, gin.H{"Error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"Message": "Success"})
	}
}
