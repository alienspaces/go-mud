package httpclient

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

// RetryRequest -
func (c *Client) RetryRequest(method, path string, params map[string]string, reqData interface{}, respData interface{}) error {
	l := loggerWithFunctionContext(c.Log, "RetryRequest")

	var err error

	// Replace placeholder parameters and add query parameters
	url, err := c.buildURL(method, path, params)
	if err != nil {
		l.Warn("failed building URL >%v<", err)
		return err
	}

	// data
	data, err := c.encodeData(reqData)
	if err != nil {
		l.Warn("failed marshalling request data >%v<", err)
		return err
	}

	// default maximum retries
	if c.MaxRetries == 0 {
		c.MaxRetries = maxRetries
	}
	retries := 0

	// Response
	var resp *http.Response

	// Error
	clientErr := Error{}

RETRY:
	for {
		retries++

		resp, err = c.Request(method, url, data)
		if err != nil {
			// Client error
			clientErr = Error{
				RequestError: err,
			}

			// Not guaranteed to get a response
			if resp != nil {
				clientErr.Code = resp.StatusCode
				buf := new(bytes.Buffer)
				if resp.Body != nil {
					_, err := buf.ReadFrom(resp.Body)
					if err != nil {
						return err
					}
				}
				respData := buf.String()
				l.Warn("client request failed retries >%d< >%v< response >%s<", retries, err, respData)

				// Never retry 4xx errors
				if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode <= http.StatusUnavailableForLegalReasons {
					l.Warn("client 4xx error, giving up immediately")
					return clientErr
				}
			}

			// Max retries
			if retries == c.MaxRetries {
				l.Warn("client request exceeded max retries, giving up now")
				return clientErr
			}

			time.Sleep(time.Duration(retries) * time.Second)
			continue RETRY
		}
		break
	}

	err = c.decodeData(resp.Body, &respData)
	if err != nil {
		msg := fmt.Sprintf("failed decoding response >%v<", err)
		l.Warn(msg)

		clientErr = Error{
			DecodeError: err,
		}

		return clientErr
	}

	return nil
}
