package pay

import (
	"context"
	"fmt"
	"i9pay/db"
	"i9pay/platform/emails"
	"time"

	"firebase.google.com/go/auth"
	"github.com/sendgrid/sendgrid-go"
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

func setUserPaying(database *mongo.Database, userID string, expires time.Time) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"paying":   true,
			"provider": "Stripe",
		},
	}

	collection := database.Collection("user")
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	fmt.Println(expires)
	fmt.Println(primitive.NewDateTimeFromTime(expires))
	userPaymentFilter := bson.M{"userid": userID}
	userPaymentUpdate := bson.M{
		"$set": bson.M{
			"processing": false,
			"expires":    primitive.NewDateTimeFromTime(expires),
		},
	}

	userPaymentCollection := database.Collection("userpayment")
	_, err = userPaymentCollection.UpdateOne(context.TODO(), userPaymentFilter, userPaymentUpdate)
	if err != nil {
		return err
	}

	return nil
}

func setUserPayingPartial(database *mongo.Database, subscriptionID, firebaseID, length, userid string) error {

	partial := db.UserPayment{
		UserMongoID:    userid,
		Username:       firebaseID,
		Provider:       "Stripe",
		SubscriptionID: subscriptionID,
		SubLength:      length,
		EndDate:        primitive.DateTime(0),
		Processing:     true,
		Ending:         false,
	}

	collection := database.Collection("userpayment")
	_, err := collection.InsertOne(context.TODO(), partial)
	if err != nil {
		return err
	}

	return nil
}

func SetUserNotPaying(client *sendgrid.Client, auth *auth.Client, database *mongo.Database, userID string, email bool) error {
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

	userPaymentFilter := bson.M{"userid": userID}
	userPaymentCollection := database.Collection("userpaying")
	_, err = userPaymentCollection.DeleteOne(context.TODO(), userPaymentFilter)
	if err != nil {
		return err
	}

	var user db.User
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return err
	}

	userRecord, err := auth.GetUser(context.Background(), user.Username)
	if err != nil {
		return err
	}

	if email {
		if err := emails.SendOver(client, userRecord.Email, user.Name); err != nil {
			return err
		}
	}

	return nil
}

func UserIdToSubscriptionId(database *mongo.Database, userID string) (string, error) {
	filter := bson.M{"userid": userID}
	projection := bson.M{"subid": 1}

	var result struct {
		SubID string `bson:"subid"`
	}

	collection := database.Collection("userpayment")
	err := collection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.SubID, nil
}

func getUserPayment(database *mongo.Database, userID string) (*db.UserPayment, error) {
	filter := bson.M{"userid": userID}

	var result db.UserPayment

	collection := database.Collection("userpayment")
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func setUserPaymentEnding(database *mongo.Database, userID string, status bool, ending time.Time) error {
	filter := bson.M{"userid": userID}
	update := bson.M{
		"$set": bson.M{
			"ending": status,
			"end":    primitive.NewDateTimeFromTime(ending),
		},
	}

	collection := database.Collection("userpayment")
	_, err := collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	return nil
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

func getCancellationByUser(database *mongo.Database, userID string) (string, error) {
	collection := database.Collection("cancellations")

	filter := bson.M{"user_id": userID}

	var result struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.ID.Hex(), nil
}

func getUserName(database *mongo.Database, userID string) (string, error) {
	collection := database.Collection("user")

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", err
	}

	filter := bson.M{"_id": objID}
	projection := bson.M{"name": 1}

	var result struct {
		Name string `bson:"name"`
	}

	err = collection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Name, nil
}

func UpdateSubscriptionLength(db *mongo.Database, userid, frequency string) error {
	collection := db.Collection("userpayment")
	filter := bson.M{"userid": userid}
	update := bson.M{"$set": bson.M{"length": frequency}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}
