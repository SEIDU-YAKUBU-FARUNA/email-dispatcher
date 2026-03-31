package api

import "email-dispatcher/models"

var emailQueue = make(chan models.EmailRequest, 100)

// Getter function (clean architecture)
func GetQueue() chan models.EmailRequest {
	return emailQueue
}
