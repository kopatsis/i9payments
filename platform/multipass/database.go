package multipass

import (
	"context"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SpecialCodeSB struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Status string             `bson:"status"`
}

func checkSpecialCode(code string, database *mongo.Database) bool {
	if code == "" {
		return false
	}

	collection := database.Collection("specialcode")
	objectID, err := primitive.ObjectIDFromHex(code)
	if err != nil {
		return false
	}

	filter := bson.M{"_id": objectID, "status": "Active"}
	var result SpecialCodeSB
	err = collection.FindOne(context.Background(), filter).Decode(&result)

	return err == nil
}

func deleteSpecialCode(code string, database *mongo.Database) error {
	collection := database.Collection("specialcode")
	objectID, err := primitive.ObjectIDFromHex(code)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = collection.DeleteOne(context.Background(), filter)
	return err
}

func emailAndIDfromToken(idt string, database *mongo.Database) (string, string, error) {

	token, _, err := new(jwt.Parser).ParseUnverified(idt, jwt.MapClaims{})
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", jwt.ErrInvalidKey
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", "", jwt.ErrInvalidKey
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", "", jwt.ErrInvalidKey
	}

	collection := database.Collection("user")

	var result struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	if err := collection.FindOne(
		context.Background(),
		bson.M{"username": sub},
		options.FindOne().SetProjection(bson.M{"_id": 1}),
	).Decode(&result); err != nil {
		return "", "", err
	}

	return result.ID.Hex(), email, nil
}
