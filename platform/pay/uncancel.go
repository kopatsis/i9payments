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
	"github.com/go-co-op/gocron"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostUncancel(auth *auth.Client, database *mongo.Database, scheduler *gocron.Scheduler) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid, userID, err := multipass.BothIDsFromCookie(c, auth, database)
		if err != nil {
			log.Printf("Failed to cancel cancellation in finding user: %s; %s", userID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel cancellation in finding user"})
			return
		}

		cancelID, err := getCancellationByUser(database, userID)
		if err != nil {
			log.Printf("Failed to get cancellation db entry for user: %s; %s", userID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cancellation db entry"})
			return
		}

		if err = deleteScheduledJob(scheduler, cancelID); err != nil {
			log.Printf("Failed to cancel cancellation schduler for user: %s; %s", userID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel cancellation schduler"})
			return
		}

		params := &stripe.SubscriptionParams{
			CancelAtPeriodEnd: stripe.Bool(false),
		}

		subID, err := userIdToSubscriptionId(database, userID)
		if err != nil {
			log.Printf("Failed to retrieve subscription from Stripe user: %s; %s", userID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscription from Stripe"})
			return
		}

		if _, err := sub.Update(subID, params); err != nil {
			log.Printf("Error in setting sub for sub ID: %s; for user: %s; %s", subID, userID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscription from Stripe"})
			return
		}

		if err = deleteCancellation(database, cancelID); err != nil {
			log.Printf("Failed to cancel cancellation db entry for user: %s; %s", userID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel cancellation db entry"})
			return
		}

		if err := setUserPaymentEnding(database, userID, false, time.Time{}); err != nil {
			log.Printf("Failed to edit db payment entry for user: %s; %s", userID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit db payment entry for user"})
			return
		}

		userRecord, err := auth.GetUser(context.Background(), uid)
		if err != nil {
			log.Printf("Error in getting user for cancel: %s; for user: %s; %s", subID, userID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user for cancel"})
			return
		}

		if err := emails.SendUnCancelled(userRecord.Email, userRecord.DisplayName); err != nil {
			log.Printf("Error in emailing user for cancel: %s; for user: %s; %s", subID, userID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to email user for cancel"})
			return
		}

		c.Status(204)

	}
}
