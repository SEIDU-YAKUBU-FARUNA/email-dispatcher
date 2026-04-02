package worker

import (
	"context"
	"log"

	"email-dispatcher/mailer"
	"email-dispatcher/models"

	"time"

	"email-dispatcher/database"
)

// Constants for retry logic

const (
	MaxRetries = 3
	BaseDelay  = 2 // seconds

)

// getBackoffDuration calculates the delay before the next retry using exponential backoff

func getBackoffDuration(retry int) time.Duration {
	return time.Duration(BaseDelay*(1<<retry)) * time.Second
}

// Worker struct holds the tools (dependencies) the worker needs [cite: 232]
type Worker struct {
	Mailer mailer.Mailer            // The interface for sending mail [cite: 234]
	Queue  chan models.EmailRequest // The internal buffer/queue [cite: 207]
}

// NewWorker is a helper to create a worker with its dependencies
func NewWorker(m mailer.Mailer, q chan models.EmailRequest) *Worker {
	return &Worker{
		Mailer: m,
		Queue:  q,
	}
}

// Start begins the worker's loop to process incoming email jobs

func (w *Worker) Start() {
	go func() {
		for job := range w.Queue {

			log.Println("Processing:", job.Recipient)

			var err error

			for retry := 0; retry <= MaxRetries; retry++ {

				err = w.Mailer.Send(context.Background(), job)

				if err == nil {
					log.Println(" Email sent:", job.Recipient)
					database.UpdateEmailStatus(job.ID, "sent")
					break
				}

				log.Println("Attempt failed:", retry, err)

				database.IncrementRetry(job.ID)

				// If max retries reached → fail permanently
				if retry == MaxRetries {
					database.UpdateEmailStatus(job.ID, "failed")
					log.Println("Email permanently failed:", job.Recipient)
					break
				}

				// Wait before retrying
				delay := getBackoffDuration(retry)
				log.Println("Retrying in:", delay)

				time.Sleep(delay)
			}
		}
	}()
}
