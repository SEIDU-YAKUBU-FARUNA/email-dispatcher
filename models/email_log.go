package models

import "time"

type EmailLog struct {
	ID         string
	Recipient  string
	TemplateID string
	Status     string
	RetryCount int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
