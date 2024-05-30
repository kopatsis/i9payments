package pay

import (
	"i9pay/platform/multipass"
	"log"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/stripe/stripe-go/v72/sub"
	"go.mongodb.org/mongo-driver/mongo"
)

func CancelPayment(auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		_, userid, err := multipass.BothIDsFromCookie(c, auth, database)
		if err != nil {
			log.Printf("Error in getting the user: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		subID, err := userIdToSubscriptionId(database, userid)
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

		endTime := time.Unix(stripeSub.CurrentPeriodEnd, 0)

		_, err = sub.Cancel(subID, nil)
		if err != nil {
			log.Printf("Error in cancelling sub for sub ID: %s; for user: %s; %s", subID, userid, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription on stripe side"})
			return
		}

		cancelID, err := backupCancellation(database, subID, userid, endTime)
		if err != nil {
			log.Printf("Error in pushing backup post for sub ID: %s; for user: %s; %s", subID, userid, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cancellation backup"})
			return
		}

		scheduler := gocron.NewScheduler(time.UTC)
		scheduleCancellation(scheduler, database, userid, cancelID, endTime)
		scheduler.StartAsync()

		c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled successfully"})
	}
}
