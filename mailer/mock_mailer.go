package mailer

import (
	"context"
	"log"
	"time"

	"email-dispatcher/models"
)

type MockMailer struct{}

// Send simulates sending email
func (m *MockMailer) Send(ctx context.Context, email models.EmailRequest) error {
	log.Println("Mock sending email to:", email.Recipient)

	// Simulate delay
	time.Sleep(2 * time.Second)

	log.Println("Mock email sent to:", email.Recipient)
	return nil
}
