package pay

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/mongo"
)

func scheduleCancellation(scheduler *gocron.Scheduler, database *mongo.Database, subID, userID, cancelID string, endTime time.Time) error {
	duration := time.Until(endTime)
	if duration <= 0 {
		return errors.New("endTime must be in the future")
	}

	fmt.Println(duration)

	job, err := scheduler.Every(duration).SingletonMode().Do(func() {
		fmt.Println("run????")
		if err := setUserNotPaying(database, userID); err == nil {
			err = deleteCancellation(database, cancelID)
			if err != nil {
				log.Printf("Error in uncancelling backup for user: %s; subID: %s; cancelID: %s; %s", userID, subID, cancelID, err.Error())
			}
			return
		}
	})

	if err != nil {
		return err
	}

	job.Tag(cancelID)

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
