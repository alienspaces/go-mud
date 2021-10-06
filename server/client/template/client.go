package client

import (
	"fmt"
	"time"

	"gitlab.com/alienspaces/go-boilerplate/server/core/client"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
)

// Client -
type Client struct {
	client.Client
}

// Request -
type Request struct {
	client.Request
	Data Data `json:"data"`
}

// Response -
type Response struct {
	client.Response
	Data []Data `json:"data"`
}

// Data -
type Data struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// NewClient -
func NewClient(c configurer.Configurer, l logger.Logger) (*Client, error) {

	cl := Client{
		client.Client{
			Config: c,
			Log:    l,
		},
	}

	// Base path for all requests
	cl.Path = "/api/templates"

	return &cl, nil
}

// GetTemplate -
func (c *Client) GetTemplate(templateID string) (*Response, error) {

	// Required parameters
	params := map[string]string{
		"id": templateID,
	}

	respData := Response{}
	err := c.Get("", params, &respData)
	if err != nil {
		msg := fmt.Sprintf("Failed getting resource >%v<", err)
		c.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}
	return &respData, nil
}

// GetTemplates -
func (c *Client) GetTemplates(params map[string]string) (*Response, error) {

	respData := Response{}
	err := c.Get("", params, &respData)
	if err != nil {
		msg := fmt.Sprintf("Failed getting resources >%v<", err)
		c.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}
	return &respData, nil
}

// CreateTemplate -
func (c *Client) CreateTemplate(templateID string, reqData *Request) (*Response, error) {

	// Required parameters
	params := map[string]string{}
	if templateID != "" {
		params["id"] = templateID
	}

	respData := Response{}
	err := c.Create("", params, reqData, &respData)
	if err != nil {
		msg := fmt.Sprintf("Failed creating resource >%v<", err)
		c.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}
	return &respData, nil
}

// UpdateTemplate -
func (c *Client) UpdateTemplate(templateID string, reqData *Request) (*Response, error) {

	// Required parameters
	params := map[string]string{
		"id": templateID,
	}

	respData := Response{}
	err := c.Update("", params, reqData, &respData)
	if err != nil {
		msg := fmt.Sprintf("Failed updating resource >%v<", err)
		c.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}
	return &respData, nil
}

// DeleteTemplate -
func (c *Client) DeleteTemplate(templateID string) (*Response, error) {

	// Required parameters
	params := map[string]string{
		"id": templateID,
	}

	respData := Response{}
	err := c.Delete("", params, &respData)
	if err != nil {
		msg := fmt.Sprintf("Failed deleting resource >%v<", err)
		c.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}
	return &respData, nil
}
