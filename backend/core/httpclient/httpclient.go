package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

const (
	maxRetries int = 5
)

type AuthType string

// Client -
type Client struct {
	Log logger.Logger

	MaxRetries int
	// Path is the base path for all requests
	Path string
	// Host is the host for all requests
	Host string
	// Setting AuthKey will result in the header "Authorisation: [AuthKey]"
	AuthKey string
	// Setting AuthToken will result in the header "Authorisation: Bearer [AuthKey]"
	AuthToken string
	// Setting AuthUser/AuthPass will result in the header "Basic: [AuthKey]"
	AuthUser string
	AuthPass string

	// Verbose will log requests at info log level as opposed to the default debug log level. This
	// helps by limiting noise from other packages when attempting to determine HTTP client interactions
	Verbose bool

	// RequestHeaders provides a way to add additional request headers
	RequestHeaders map[string]string
}

// Error implements the error interface while also exposing the underlying HTTP response
type Error struct {
	Code         int
	RequestError error
	DecodeError  error
}

func (e Error) Error() string {
	return fmt.Sprintf("Code >%d< RequestError >%v< DecodeError >%v<", e.Code, e.RequestError, e.DecodeError)
}

// Request -
type Request struct {
	Pagination *RequestPagination `json:"pagination,omitempty"`
}

// RequestPagination -
type RequestPagination struct {
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
}

// Response -
type Response struct {
	Error      *ResponseError      `json:"error,omitempty"`
	Pagination *ResponsePagination `json:"pagination,omitempty"`
}

// ResponseError -
type ResponseError struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// ResponsePagination -
type ResponsePagination struct {
	Number int `json:"page_number"`
	Size   int `json:"page_size"`
	Count  int `json:"page_count"`
}

// NewClient -
func NewClient(l logger.Logger) (*Client, error) {
	cl := Client{
		Log: l,
	}
	return &cl, nil
}

// setHeaders sets various headers as per configuration
func (c *Client) setHeaders(req *http.Request) error {

	if c.RequestHeaders != nil {
		for key, value := range c.RequestHeaders {
			if c.Verbose {
				c.Log.Info("Setting header >%s< .%s<", key, value)
			} else {
				c.Log.Info("Setting header >%s< .%s<", key, value)
			}
			req.Header.Add(key, value)
		}
	}

	return nil
}

// setAuthHeaders sets authentication headers on an request object based
// on client authentication configuration
func (c *Client) setAuthHeaders(req *http.Request) error {

	// Authorization with key
	if c.AuthKey != "" {
		req.Header.Add("Authorization", c.AuthKey)
		return nil
	}

	// Authorization with bearer token
	if c.AuthToken != "" {
		var bearer = "Bearer " + c.AuthToken
		req.Header.Add("Authorization", bearer)
		return nil
	}

	// Authorization with user and pass
	if c.AuthUser != "" && c.AuthPass != "" {
		req.SetBasicAuth(c.AuthUser, c.AuthPass)
		return nil
	}

	return nil
}

// buildURL replaces placeholder parameters and adds query parameters
// The parameter "id" or ":id" has special behaviour. When provided the
// returned URL will have "/:id" appended and replaced with whatever
// the parameter value for "id" or ":id" was.
func (c *Client) buildURL(method, requestURL string, params map[string]string) (string, error) {
	l := loggerWithFunctionContext(c.Log, "buildURL")

	// Request URL
	requestURL = c.Host + c.Path + requestURL

	// Replace placeholders and add query parameters
	paramString := ""
	for param, value := range params {

		// do not allow empty param values
		if value == "" {
			return requestURL, fmt.Errorf("param >%s< has empty value", param)
		}

		found := false
		if strings.Contains(requestURL, "/:"+param) {
			requestURL = strings.Replace(requestURL, "/:"+param, "/"+value, 1)
			found = true
		}
		if strings.Contains(requestURL, "/"+param) {
			requestURL = strings.Replace(requestURL, "/"+param, "/"+value, 1)
			found = true
		}
		if strings.Contains(requestURL, "/"+strings.Replace(param, ":", "", 1)) {
			requestURL = strings.Replace(requestURL, "/"+strings.Replace(param, ":", "", 1), "/"+value, 1)
			found = true
		}
		if !found {
			param = strings.Replace(param, ":", "", 1)
			if paramString != "" {
				paramString = paramString + "&"
			}
			paramString = paramString + param + "=" + url.QueryEscape(value)
		}
	}

	if paramString != "" {
		requestURL = requestURL + "?" + paramString
	}

	// do not allow missing parameters
	if strings.Contains(requestURL, "/:") {
		return requestURL, fmt.Errorf("URL >%s< still contains placeholders", requestURL)
	}

	if c.Verbose {
		l.Info("Request URL >%s<", requestURL)
	} else {
		l.Debug("Request URL >%s<", requestURL)
	}

	return requestURL, nil
}

// encodeData is a convenience function that encodes struct data into bytes
func (c *Client) encodeData(data interface{}) ([]byte, error) {
	l := loggerWithFunctionContext(c.Log, "encodeData")

	dataBytes, err := json.Marshal(data)
	if err != nil {
		l.Warn("failed encoding data >%v<", err)
		return nil, err
	}
	return dataBytes, nil
}

// decodeData is a convenience function that decodes bytes into struct data
func (c *Client) decodeData(rc io.ReadCloser, data interface{}) error {
	l := loggerWithFunctionContext(c.Log, "decodeData")

	// close before returning
	defer rc.Close()

	err := json.NewDecoder(rc).Decode(&data)
	if err != nil && err.Error() != "EOF" {
		l.Warn("failed decoding data >%v<", err)
		return err
	}
	return nil
}

// loggerWithFunctionContext - Returns a logger with package context and provided function context
func loggerWithFunctionContext(l logger.Logger, functionName string) logger.Logger {
	return l.WithPackageContext("client").WithFunctionContext(functionName)
}
