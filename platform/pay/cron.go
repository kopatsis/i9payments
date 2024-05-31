package pay

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/mongo"
)

func scheduleCancellation(scheduler *gocron.Scheduler, database *mongo.Database, userID, cancelID string, endTime time.Time) {
	scheduler.At(endTime).Do(func() {
		err := setUserNotPaying(database, userID)
		if err == nil {
			err = deleteCancellation(database, cancelID)
			if err != nil {
				log.Printf("Error in uncancelling backup for user: %s; cancelID: %s; %s", userID, cancelID, err.Error())
			}
			return
		}
		log.Printf("Error in cancelling actual for user: %s; cancelID: %s; %s", userID, cancelID, err.Error())
	})
}
