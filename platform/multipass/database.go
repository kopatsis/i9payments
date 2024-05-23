package multipass

import (
	"context"
	"fmt"
	"i9pay/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SpecialCodeSB struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Status string             `bson:"status"`
}

func getSpecialCode(database *mongo.Database) (string, error) {
	specialCode := SpecialCodeSB{
		Status: "Archived",
	}

	collection := database.Collection("specialcode")
	insertResult, err := collection.InsertOne(context.Background(), specialCode)
	if err != nil {
		return "", err
	}

	objectID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("invalid ID type")
	}

	return objectID.Hex(), nil
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
