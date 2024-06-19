package pay

import (
	"context"
	"i9pay/platform/emails"
	"i9pay/platform/multipass"
	"log"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
	"go.mongodb.org/mongo-driver/mongo"
)

func CancelPayment(client *sendgrid.Client, auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid, userid, err := multipass.BothIDsFromCookie(c, auth, database)
		if err != nil {
			log.Printf("Error in getting the user: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		subID, err := UserIdToSubscriptionId(database, userid)
		if err != nil || subID == "" {
			log.Printf("Error in getting subID for user: %s; %s", userid, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Some issue with user"})
			return
		}

		stripeSub, err := sub.Get(subID, nil)
		if err != nil {
			log.Printf("Error in getting sub for sub ID: %s; for user: %s; %s", subID, userid, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscription from Stripe"})
			return
		}

		params := &stripe.SubscriptionParams{
			CancelAtPeriodEnd: stripe.Bool(true),
		}

		if _, err := sub.Update(subID, params); err != nil {
			log.Printf("Error in setting sub for sub ID: %s; for user: %s; %s", subID, userid, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscription from Stripe"})
			return
		}

		endTime := time.Unix(stripeSub.CurrentPeriodEnd, 0)
		if err := setUserPaymentEnding(database, userid, true, endTime); err != nil {
			log.Printf("Error in setting user payment: %s; for user: %s; %s", subID, userid, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set user payment"})
			return
		}

		_, err = backupCancellation(database, subID, userid, endTime)
		if err != nil {
			log.Printf("Error in pushing backup post for sub ID: %s; for user: %s; %s", subID, userid, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cancellation backup"})
			return
		}

		userRecord, err := auth.GetUser(context.Background(), uid)
		if err != nil {
			log.Printf("Error in getting user for cancel: %s; for user: %s; %s", subID, userid, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user for cancel"})
			return
		}

		if err := emails.SendCancelled(client, userRecord.Email, userRecord.DisplayName); err != nil {
			log.Printf("Error in emailing user for cancel: %s; for user: %s; %s", subID, userid, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to email user for cancel"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled successfully"})
	}
}
