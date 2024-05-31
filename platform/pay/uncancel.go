package pay

import (
	"i9pay/platform/multipass"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostUncancel(auth *auth.Client, database *mongo.Database, scheduler *gocron.Scheduler) gin.HandlerFunc {
	return func(c *gin.Context) {

		_, userID, err := multipass.BothIDsFromCookie(c, auth, database)
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

		if err = deleteCancellation(database, cancelID); err != nil {
			log.Printf("Failed to cancel cancellation db entry for user: %s; %s", userID, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel cancellation db entry"})
			return
		}

		c.Status(204)

	}
}
