package pay

import (
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

		var req struct {
			UserID string `form:"user"`
		}

		if err := c.ShouldBind(&req); err != nil {
			log.Printf("Error in binding the request payload: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		subID, err := userIdToSubscriptionId(database, req.UserID)
		if err != nil || subID == "" {
			log.Printf("Error in getting subID for user: %s; %s", req.UserID, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Some issue with user"})
			return
		}

		stripeSub, err := sub.Get(subID, nil)
		if err != nil {
			log.Printf("Error in getting sub for sub ID: %s; for user: %s; %s", subID, req.UserID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscription from Stripe"})
			return
		}

		endTime := time.Unix(stripeSub.CurrentPeriodEnd, 0)

		_, err = sub.Cancel(subID, nil)
		if err != nil {
			log.Printf("Error in cancelling sub for sub ID: %s; for user: %s; %s", subID, req.UserID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription on stripe side"})
			return
		}

		cancelID, err := backupCancellation(database, subID, req.UserID, endTime)
		if err != nil {
			log.Printf("Error in pushing backup post for sub ID: %s; for user: %s; %s", subID, req.UserID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cancellation backup"})
			return
		}

		scheduler := gocron.NewScheduler(time.UTC)
		scheduleCancellation(scheduler, database, req.UserID, cancelID, endTime)
		scheduler.StartAsync()

		c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled successfully"})
	}
}
