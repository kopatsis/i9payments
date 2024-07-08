package login

import (
	"context"
	"fmt"
	"i9pay/db"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBToken struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID string             `bson:"user"`
	Token  string             `bson:"token"`
}

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

		if err := postToDBUser(username, request.Name, request.RefreshToken, database); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to post user to db for cookie"})
			return
		}

		fmt.Println(request)

		if err := Cookie(request.IDToken, request.RefreshToken, authClient, c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to create a session cookie"})
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "properLogin",
			Value:    "TRUE",
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
		})

		c.JSON(http.StatusOK, gin.H{"message": "Session cookie set successfully"})
	}
}

func postToDBUser(username, name, refresh string, database *mongo.Database) error {

	filter := bson.M{"username": username}

	collection := database.Collection("user")
	var exuser db.User

	err := collection.FindOne(context.TODO(), filter).Decode(&exuser)
	if err == nil {
		refreshTokenDB(exuser.ID.Hex(), refresh, database)
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

	res, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	refreshTokenDB(res.InsertedID.(primitive.ObjectID).Hex(), refresh, database)
	return nil
}

func refreshTokenDB(userid, refreshToken string, database *mongo.Database) error {
	collection := database.Collection("usertoken")

	filter := bson.M{"user": userid}
	update := bson.M{
		"$set": bson.M{
			"token": refreshToken,
		},
		"$setOnInsert": bson.M{
			"_id":  primitive.NewObjectID(),
			"user": userid,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)

	if err != nil {
		return fmt.Errorf("failed to update or insert token: %v", err)
	}

	return nil
}
