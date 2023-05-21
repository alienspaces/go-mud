package template

import (
	"time"

	"gitlab.com/alienspaces/go-mud/backend/schema"
)

// Response -
type Response struct {
	schema.Response
	Data *Data `json:"data"`
}

type CollectionResponse struct {
	schema.Response
	Data []*Data `json:"data"`
}

// Request -
type Request struct {
	schema.Request
	Data *Data `json:"data"`
}

// Data -
type Data struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
