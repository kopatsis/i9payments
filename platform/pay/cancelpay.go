package pay

import (
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72/sub"
	"go.mongodb.org/mongo-driver/mongo"
)

func CancelPayment(auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req struct {
			UserID string `json:"user"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		subID, err := userIdToSubscriptionId(database, req.UserID)
		if err != nil || subID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Some issue with user"})
			return
		}

		_, err = sub.Cancel(subID, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription on stripe side"})
			return
		}

		if err := setUserNotPaying(database, req.UserID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription on db side"})
			return
		}
	}
}
