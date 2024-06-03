package login

import (
	"context"
	"fmt"
	"i9pay/db"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func VerifyToken(authClient *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			IDToken      string `json:"idToken"`
			RefreshToken string `json:"refreshToken"`
			Name         string `json:"name"`
		}

		fmt.Println(request.RefreshToken)

		if err := c.ShouldBindJSON(&request); err != nil {
			fmt.Println("wrong req")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		token, err := authClient.VerifyIDToken(context.Background(), request.IDToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get user for cookie"})
			return
		}
		username := token.UID

		if err := postToDBUser(username, request.Name, database); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to post user to db for cookie"})
			return
		}

		fmt.Println(request)

		if err := Cookie(request.IDToken, request.RefreshToken, authClient, c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to create a session cookie"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Session cookie set successfully"})
	}
}

func postToDBUser(username, name string, database *mongo.Database) error {

	filter := bson.M{"username": username}

	collection := database.Collection("user")

	err := collection.FindOne(context.TODO(), filter).Err()
	if err == nil {
		return nil
	}
	if err != mongo.ErrNoDocuments {
		return err
	}

	user := db.User{
		Username:          username,
		Name:              name,
		PlyoTolerance:     3,
		PushupSetting:     "Knee",
		BannedExercises:   []string{},
		BannedStretches:   []string{},
		BannedParts:       []int{},
		ExerFavoriteRates: map[string]float32{},
		ExerModifications: map[string]float32{},
		TypeModifications: map[string]float32{},
		RoundEndurance:    map[int]float32{},
		TimeEndurance:     map[int]float32{},
	}

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	return nil
}
