package schema

import "time"

// TemplateResponse -
type TemplateResponse struct {
	Response
	Data []TemplateData `json:"data"`
}

// TemplateRequest -
type TemplateRequest struct {
	Request
	Data TemplateData `json:"data"`
}

// TemplateData -
type TemplateData struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
