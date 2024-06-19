package pay

import (
	"context"
	"i9pay/db"
	"log"
	"time"

	"firebase.google.com/go/auth"
	"github.com/sendgrid/sendgrid-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DoneCancels(client *sendgrid.Client, database *mongo.Database, auth *auth.Client) {
	collection := database.Collection("cancellations")

	filter := bson.M{"end_time": bson.M{"$gt": time.Now()}}
	opts := options.Find().SetSort(bson.D{{Key: "end_time", Value: 1}})

	cursor, err := collection.Find(context.Background(), filter, opts)
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
		err := SetUserNotPaying(client, auth, database, cancellation.UserID, true)
		if err != nil {
			log.Printf("Error setting user not paying for userID: %s, error: %v", cancellation.UserID, err)
			continue
		}

		err = deleteCancellation(database, cancellation.ID.Hex())
		if err != nil {
			log.Printf("Error deleting cancellation for cancelID: %s, error: %v", cancellation.ID.Hex(), err)
		}
	}

	userPaymentCollection := database.Collection("userpayment")

	filter = bson.M{"expires": bson.M{"$lt": primitive.NewDateTimeFromTime(time.Now().Add(-96 * time.Hour))}}
	opts = options.Find().SetSort(bson.D{{Key: "expires", Value: 1}})

	userPaymentCursor, err := userPaymentCollection.Find(context.Background(), filter, opts)
	if err != nil {
		log.Printf("Can't get expired user payments: %v", err)
		return
	}
	defer userPaymentCursor.Close(context.Background())

	var expiredPayments []db.UserPayment
	if err = userPaymentCursor.All(context.Background(), &expiredPayments); err != nil {
		log.Printf("Can't decode expired user payments: %v", err)
		return
	}

	for _, payment := range expiredPayments {
		err := SetUserNotPaying(client, auth, database, payment.UserMongoID, true)
		if err != nil {
			log.Printf("Error setting user not paying for userID: %s, error: %v", payment.UserMongoID, err)
			continue
		}
	}

}
