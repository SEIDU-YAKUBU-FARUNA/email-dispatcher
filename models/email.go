package models

type EmailRequest struct {
	ID         string `json:"id"`
	Recipient  string `json:"recipient" binding:"required,email"`
	Subject    string `json:"subject" binding:"required,min=3"`
	Body       string `json:"body" binding:"required"`
	TemplateID string `json:"template_id"`
}
