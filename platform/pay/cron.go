package pay

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DoneCancels(database *mongo.Database) {
	collection := database.Collection("cancellations")

	filter := bson.M{"end_time": bson.M{"$gt": time.Now()}}
	options := options.Find().SetSort(bson.D{{Key: "end_time", Value: 1}})

	cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		log.Printf("Can't get jobs")
		return
	}
	defer cursor.Close(context.Background())

	var cancellations []SubscriptionCancellation
	if err = cursor.All(context.Background(), &cancellations); err != nil {
		log.Printf("Can't get jobs")
		return
	}

	for _, cancellation := range cancellations {
		err := setUserNotPaying(database, cancellation.UserID)
		if err != nil {
			log.Printf("Error setting user not paying for userID: %s, error: %v", cancellation.UserID, err)
			continue
		}

		err = deleteCancellation(database, cancellation.ID.Hex())
		if err != nil {
			log.Printf("Error deleting cancellation for cancelID: %s, error: %v", cancellation.ID.Hex(), err)
		}
	}
}
