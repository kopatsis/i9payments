package pay

import (
	"i9pay/platform/multipass"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateFrequency(auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		newlength := c.PostForm("frequency")

		_, userid, err := multipass.BothIDsFromCookie(c, auth, database)
		if err != nil {
			log.Printf("Error in getting the user: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error in getting the user"})
			return
		}

		payment, err := getUserPayment(database, userid)
		if err != nil {
			log.Printf("Error in getting the user payment: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error in getting the user payment"})
			return
		}

		if payment.SubLength == newlength {
			c.Status(204)
			return
		}

		priceID := "price_1PJfbQIstWH7VBmuNNsoLTN2"
		if newlength == "yearly" {
			priceID = "price_1PJfbpIstWH7VBmu1nToVdC9"
		}

		if err := UpdateSubscriptionPlan(payment.SubscriptionID, priceID); err != nil {
			log.Printf("Error in actually updating the user frequency: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error in actually updating the user frequency"})
			return
		}

		if err := UpdateSubscriptionLength(database, userid, newlength); err != nil {
			log.Printf("Error in actually updating the db user frequency: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error in actually updating the db user frequency"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"New length": newlength})
	}
}
