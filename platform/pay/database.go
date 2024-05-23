package pay

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setUserPaying(database *mongo.Database, subscriptionID, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"paying":   true,
			"provider": subscriptionID,
		},
	}

	collection := database.Collection("user")
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func setUserNotPaying(database *mongo.Database, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"paying":   false,
			"provider": "",
		},
	}

	collection := database.Collection("user")
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func userIdToSubscriptionId(database *mongo.Database, userID string) (string, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", err
	}

	filter := bson.M{"_id": objID}
	projection := bson.M{"provider": 1}

	var result struct {
		Provider string `bson:"provider"`
	}

	collection := database.Collection("user")
	err = collection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Provider, nil
}
