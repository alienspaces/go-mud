package client

import (
	"fmt"
	"net/http"
)

// Get is a convenience method wrapping RetryRequest
func (c *Client) Get(path string, params map[string]string, respData interface{}) error {

	c.Log.Context("function", "Get")
	defer func() {
		c.Log.Context("function", "")
	}()

	c.Log.Debug("Request path >%s< params >%#v< respData >%#v<", path, params, respData)

	err := c.RetryRequest(
		http.MethodGet,
		path,
		params,
		nil,
		respData,
	)
	if err != nil {
		c.Log.Warn(fmt.Sprintf("Failed request >%v<", err))
		return err
	}

	return nil
}

// Create is a convenience method wrapping RetryRequest
func (c *Client) Create(path string, params map[string]string, reqData interface{}, respData interface{}) error {

	c.Log.Context("function", "Create")
	defer func() {
		c.Log.Context("function", "")
	}()

	c.Log.Debug("Request path >%s< params >%#v< reqData >%#v< respData >%#v<", path, params, reqData, respData)

	if reqData == nil {
		msg := fmt.Sprintf("Request data is nil >%v<, cannot create resource", reqData)
		c.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	err := c.RetryRequest(
		http.MethodPost,
		path,
		params,
		reqData,
		respData,
	)
	if err != nil {
		c.Log.Warn(fmt.Sprintf("Failed request >%v<", err))
		return err
	}

	return nil
}

// Update is a convenience method wrapping RetryRequest
func (c *Client) Update(path string, params map[string]string, reqData interface{}, respData interface{}) error {

	c.Log.Context("function", "UpdateTemplate")
	defer func() {
		c.Log.Context("function", "")
	}()

	c.Log.Debug("Request path >%s< params >%#v< reqData >%#v< respData >%#v<", path, params, reqData, respData)

	if reqData == nil {
		msg := fmt.Sprintf("Request data is nil >%v<, cannot update resource", reqData)
		c.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	err := c.RetryRequest(
		http.MethodPut,
		path,
		params,
		reqData,
		respData,
	)
	if err != nil {
		c.Log.Warn(fmt.Sprintf("Failed request >%v<", err))
		return err
	}

	return nil
}

// Delete is a convenience method wrapping RetryRequest
func (c *Client) Delete(path string, params map[string]string, respData interface{}) error {

	c.Log.Context("function", "Delete")
	defer func() {
		c.Log.Context("function", "")
	}()

	err := c.RetryRequest(
		http.MethodDelete,
		path,
		params,
		nil,
		respData,
	)
	if err != nil {
		c.Log.Warn(fmt.Sprintf("Failed request >%v<", err))
		return err
	}

	return nil
}
