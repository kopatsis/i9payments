package pay

import (
	"context"
	"encoding/json"
	"fmt"
	"i9pay/db"
	"i9pay/platform/emails"
	"log"
	"net/http"
	"os"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/stripe/stripe-go/v72/webhook"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func WebhookConfirm(client *sendgrid.Client, auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		const MaxBodyBytes = int64(65536)
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

		payload, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
			return
		}

		endpointSecret := os.Getenv("END_SECR_CONF")

		event, err := webhook.ConstructEvent(payload, c.GetHeader("Stripe-Signature"), endpointSecret)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error verifying webhook signature"})
			return
		}

		if event.Type == "invoice.payment_succeeded" {
			var invoice stripe.Invoice
			if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing webhook JSON"})
				return
			}

			subscription := invoice.Subscription
			customerID := invoice.Customer.ID

			customerReal, err := customer.Get(customerID, nil)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error customer doesn't exist"})
				return
			}

			paymentIntent, _ := paymentintent.Get(invoice.PaymentIntent.ID, nil)
			userId := paymentIntent.Metadata["userId"]

			if userId == "" {
				userId = customerReal.Metadata["userId"]
			}

			if userId == "" {
				var userPayment db.UserPayment
				filter := bson.M{"subid": subscription.ID}
				err := database.Collection("userpayment").FindOne(context.TODO(), filter).Decode(&userPayment)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Error user id doesn't exist in either metafield or userpayment collection"})
					return
				}
				userId = userPayment.UserMongoID
			}

			fmt.Println(subscription.CurrentPeriodEnd)
			if err := setUserPaying(database, subscription.ID, userId, time.Unix(subscription.CurrentPeriodEnd, 0)); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error updating the user"})
				return
			}

			params := &stripe.SubscriptionParams{
				DefaultPaymentMethod: stripe.String(paymentIntent.PaymentMethod.ID),
			}
			sub.Update(subscription.ID, params)

			objID, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error converting user id to ObjectID"})
				return
			}
			var user db.User
			userFilter := bson.M{"_id": objID}
			err = database.Collection("user").FindOne(context.TODO(), userFilter).Decode(&user)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving user from user collection"})
				return
			}

			userRecord, err := auth.GetUser(context.Background(), user.Username)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving user from firebase"})
				return
			}

			if err := emails.SendConfirmation(client, userRecord.Email, user.Name); err != nil {
				log.Printf("Error in emailing user: %s; %s", invoice.CustomerEmail, err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to email user for cancel"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "success"})

		}
	}
}

func WebhookFail(client *sendgrid.Client, auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		const MaxBodyBytes = int64(65536)
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

		payload, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
			return
		}

		endpointSecret := os.Getenv("END_SECR_FAIL")

		event, err := webhook.ConstructEvent(payload, c.GetHeader("Stripe-Signature"), endpointSecret)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error verifying webhook signature"})
			return
		}

		if event.Type == "invoice.payment_failed" {
			var invoice stripe.Invoice
			if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing webhook JSON"})
				return
			}

			subscription := invoice.Subscription
			customerID := invoice.Customer.ID

			customerReal, err := customer.Get(customerID, nil)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error customer doesn't exist"})
				return
			}

			paymentIntent, _ := paymentintent.Get(invoice.PaymentIntent.ID, nil)
			userId := paymentIntent.Metadata["userId"]

			if userId == "" {
				userId = customerReal.Metadata["userId"]
			}

			if userId == "" {
				var userPayment db.UserPayment
				filter := bson.M{"subid": subscription.ID}
				err := database.Collection("userpayment").FindOne(context.TODO(), filter).Decode(&userPayment)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Error user id doesn't exist in either metafield or userpayment collection"})
					return
				}
				userId = userPayment.UserMongoID
			}

			objID, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error converting user id to ObjectID"})
				return
			}
			var user db.User
			userFilter := bson.M{"_id": objID}
			err = database.Collection("user").FindOne(context.TODO(), userFilter).Decode(&user)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving user from user collection"})
				return
			}

			userRecord, err := auth.GetUser(context.Background(), user.Username)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving user from firebase"})
				return
			}

			if err := emails.SendFailureNotification(client, userRecord.Email, user.Name); err != nil {
				log.Printf("Error in emailing user: %s; %s", userRecord.Email, err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to email user for payment failure"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "success"})
		}
	}
}
