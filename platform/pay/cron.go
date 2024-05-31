package pay

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/stripe/stripe-go/v72/sub"
	"go.mongodb.org/mongo-driver/mongo"
)

func scheduleCancellation(scheduler *gocron.Scheduler, database *mongo.Database, subID, userID, cancelID string, endTime time.Time) error {
	job, err := scheduler.At(endTime).Do(func() {

		_, err := sub.Cancel(subID, nil)
		if err != nil {
			log.Printf("Error in actual cancellation for user: %s; subID: %s; cancelID: %s; %s", userID, subID, cancelID, err.Error())
			return
		}

		if err := setUserNotPaying(database, userID); err == nil {
			err = deleteCancellation(database, cancelID)
			if err != nil {
				log.Printf("Error in uncancelling backup for user: %s; subID: %s; cancelID: %s; %s", userID, subID, cancelID, err.Error())
			}
			return
		}
		log.Printf("Error in cancelling actual for user: %s; subID: %s; cancelID: %s; %s", userID, subID, cancelID, err.Error())
	})

	job.Tag(cancelID)

	if err != nil {
		return err
	}
	return nil
}

func deleteScheduledJob(scheduler *gocron.Scheduler, cancelID string) error {
	jobs, err := scheduler.FindJobsByTag(cancelID)
	if err != nil {
		return err
	}

	if len(jobs) == 0 {
		return fmt.Errorf("no job found with cancel ID: %s", cancelID)
	}

	for _, job := range jobs {
		scheduler.RemoveByReference(job)
	}
	return nil
}
