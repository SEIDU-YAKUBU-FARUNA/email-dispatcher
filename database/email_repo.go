//testing code

package database

import (
	"email-dispatcher/models"
	"time"
)

// Insert new email
func CreateEmailLog(log models.EmailLog) error {
	query := `
	INSERT INTO email_logs (id, recipient, template_id, status, retry_count, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := DB.Exec(query,
		log.ID,
		log.Recipient,
		log.TemplateID,
		log.Status,
		log.RetryCount,
		log.CreatedAt,
		log.UpdatedAt,
	)

	return err
}

// Update email status
func UpdateEmailStatus(id string, status string) error {
	query := `
	UPDATE email_logs
	SET status=$1, updated_at=$2
	WHERE id=$3
	`

	_, err := DB.Exec(query, status, time.Now(), id)
	return err
}

// Get email by ID
func GetEmailByID(id string) (*models.EmailLog, error) {
	query := `
	SELECT id, recipient, template_id, status, retry_count, created_at, updated_at
	FROM email_logs WHERE id=$1
	`

	row := DB.QueryRow(query, id)

	var log models.EmailLog

	err := row.Scan(
		&log.ID,
		&log.Recipient,
		&log.TemplateID,
		&log.Status,
		&log.RetryCount,
		&log.CreatedAt,
		&log.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &log, nil
}

func IncrementRetry(id string) {
	query := `
	UPDATE email_logs
	SET retry_count = retry_count + 1
	WHERE id=$1
	`
	DB.Exec(query, id)
}

//production ready code
/**

package database

import (
	"context"
	"email-dispatcher/models"
	"fmt"
	"regexp"
	"time"
)

// Valid statuses for email logs
var validStatuses = map[string]bool{
	"pending": true,
	"sent":    true,
	"failed":  true,
	"retry":   true,
}

// Validate email format (simple but effective)
func isValidEmail(email string) bool {
	// Basic email regex - in production, use a proper validation library
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(email)
}

// Validate UUID format for IDs (if you're using UUIDs)
func isValidUUID(id string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(id)
}

// Insert new email - IMPROVED VERSION
func CreateEmailLog(ctx context.Context, log models.EmailLog) error {
	// INPUT VALIDATION
	if log.ID == "" {
		return fmt.Errorf("email log ID cannot be empty")
	}

	if !isValidEmail(log.Recipient) {
		return fmt.Errorf("invalid recipient email address: %s", log.Recipient)
	}

	if log.TemplateID == "" {
		return fmt.Errorf("template ID cannot be empty")
	}

	if log.Status == "" {
		log.Status = "pending" // Set default status if not provided
	}

	if !validStatuses[log.Status] {
		return fmt.Errorf("invalid status: %s", log.Status)
	}

	// Use context with timeout (5 seconds max)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    INSERT INTO email_logs (id, recipient, template_id, status, retry_count, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := DB.ExecContext(ctx, query,
		log.ID,
		log.Recipient,
		log.TemplateID,
		log.Status,
		log.RetryCount,
		log.CreatedAt,
		log.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create email log for recipient %s: %w", log.Recipient, err)
	}

	return nil
}

// Update email status - IMPROVED VERSION
func UpdateEmailStatus(ctx context.Context, id string, status string) error {
	// INPUT VALIDATION
	if id == "" {
		return fmt.Errorf("email ID cannot be empty")
	}

	if !isValidUUID(id) {
		return fmt.Errorf("invalid email ID format: %s", id)
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s. Allowed: pending, sent, failed, retry", status)
	}

	// Use context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    UPDATE email_logs
    SET status=$1, updated_at=$2
    WHERE id=$3
    `

	result, err := DB.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update status for email %s: %w", id, err)
	}

	// Check if any row was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("email with ID %s not found", id)
	}

	return nil
}

// Get email by ID - IMPROVED VERSION
func GetEmailByID(ctx context.Context, id string) (*models.EmailLog, error) {
	// INPUT VALIDATION
	if id == "" {
		return nil, fmt.Errorf("email ID cannot be empty")
	}

	if !isValidUUID(id) {
		return nil, fmt.Errorf("invalid email ID format: %s", id)
	}

	// Use context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
    SELECT id, recipient, template_id, status, retry_count, created_at, updated_at
    FROM email_logs WHERE id=$1
    `

	row := DB.QueryRowContext(ctx, query, id)

	var log models.EmailLog
	err := row.Scan(
		&log.ID,
		&log.Recipient,
		&log.TemplateID,
		&log.Status,
		&log.RetryCount,
		&log.CreatedAt,
		&log.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get email by ID %s: %w", id, err)
	}

	return &log, nil
}


**/
