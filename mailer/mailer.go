package mailer

import (
	"context"
	"email-dispatcher/models"
)

// Mailer defines how emails are sent
// Any provider (SendGrid, Mailgun, etc.) must implement this
type Mailer interface {
	Send(ctx context.Context, email models.EmailRequest) error
}
