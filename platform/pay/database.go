package pay

import (
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SubscriptionCancellation struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	SubID   string             `bson:"sub_id"`
	UserID  string             `bson:"user_id"`
	EndTime time.Time          `bson:"end_time"`
}

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

func scheduleCancellation(scheduler *gocron.Scheduler, database *mongo.Database, userID, cancelID string, endTime time.Time) {
	scheduler.At(endTime).Do(func() {
		err := setUserNotPaying(database, userID)
		if err == nil {
			err = deleteCancellation(database, cancelID)
			if err != nil {
				log.Printf("Error in uncancelling backup for user: %s; cancelID: %s; %s", userID, cancelID, err.Error())
			}
			return
		}
		log.Printf("Error in cancelling actual for user: %s; cancelID: %s; %s", userID, cancelID, err.Error())
	})
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

func backupCancellation(database *mongo.Database, subID, userID string, endTime time.Time) (string, error) {
	cancellation := SubscriptionCancellation{
		SubID:   subID,
		UserID:  userID,
		EndTime: endTime,
	}

	collection := database.Collection("cancellations")
	ret, err := collection.InsertOne(context.TODO(), cancellation)
	if err != nil {
		return "", err
	} else {
		return ret.InsertedID.(primitive.ObjectID).Hex(), nil
	}
}

func deleteCancellation(database *mongo.Database, cancelID string) error {
	objID, err := primitive.ObjectIDFromHex(cancelID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}

	collection := database.Collection("cancellations")
	_, err = collection.DeleteOne(context.TODO(), filter)
	return err
}
