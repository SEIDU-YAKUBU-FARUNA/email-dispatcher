package api

import (
	"fmt"
	"net/http"
	"time"

	"email-dispatcher/database"
	"email-dispatcher/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SendEmailHandler(c *gin.Context) {
	var req models.EmailRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	// Rate limit per recipient
	allowed, err := database.AllowRequest(req.Recipient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Rate limiter error",
		})
		return
	}

	if !allowed {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Too many requests. Try later.",
		})
		return
	}

	// Generate job ID
	jobID := uuid.New().String()

	req.ID = jobID

	// Save to DB FIRST (important)
	emailLog := models.EmailLog{
		ID:         jobID,
		Recipient:  req.Recipient,
		TemplateID: req.TemplateID,
		Status:     "pending",
		RetryCount: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	// Save email log to DB
	err = database.CreateEmailLog(emailLog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save email",
		})
		return
	}

	// Push to queue
	GetQueue() <- req

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Email queued successfully",
		"job_id":  jobID,
	})

	// Create unique key (important)
	key := fmt.Sprintf("email:%s:%s", req.Recipient, req.Subject)

	// Check duplicate
	var duplicate bool
	duplicate, err = database.IsDuplicate(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Server error",
		})
		return
	}

	if duplicate {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Duplicate request detected",
		})
		return
	}

}

func GetStatusHandler(c *gin.Context) {
	id := c.Param("id")

	// Fetch email log from DB

	email, err := database.GetEmailByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Email not found",
		})
		return
	}
	// Return email status and recipient for better visibility
	c.JSON(http.StatusOK, gin.H{
		"id":        email.ID,
		"status":    email.Status,
		"recipient": email.Recipient,
	})
}
