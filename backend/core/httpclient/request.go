package httpclient

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
)

// Request -
func (c *Client) Request(method, url string, data []byte) (*http.Response, error) {
	l := loggerWithFunctionContext(c.Log, "Request")

	var err error

	if c.Verbose {
		l.Info("Client request URL >%s< data length >%d<", url, len(data))
	} else {
		l.Debug("Client request URL >%s< data length >%d<", url, len(data))
	}

	var resp *http.Response
	var req *http.Request

	client := &http.Client{}

	switch method {
	case http.MethodGet:

		// Get
		if c.Verbose {
			l.Info("Method %s", method)
		} else {
			l.Debug("Method %s", method)
		}

		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			l.Warn("failed client request >%v<", err)
			return nil, err
		}

		err := c.setHeaders(req)
		if err != nil {
			l.Warn("failed setting headers >%v<", err)
			return nil, err
		}

		err = c.setAuthHeaders(req)
		if err != nil {
			l.Warn("failed setting request auth headers >%v<", err)
			return nil, err
		}

		if c.Verbose {
			dmp, err := httputil.DumpRequest(req, true)
			if err != nil {
				l.Warn("failed request dump >%v<", err)
				return nil, err
			}
			l.Info("%s", string(dmp[:]))
		}

		resp, err = client.Do(req)
		if err != nil {
			l.Warn("failed client request >%v<", err)
			return resp, err
		}

		if c.Verbose {
			dmp, err := httputil.DumpResponse(resp, true)
			if err != nil {
				l.Warn("failed response dump >%v<", err)
				return nil, err
			}
			l.Info("%s", string(dmp[:]))
		}

	case http.MethodPost, http.MethodPut:

		// Post / Put
		if c.Verbose {
			l.Info("Method %s", method)
		} else {
			l.Debug("Method %s", method)
		}

		req, err = http.NewRequest(method, url, bytes.NewBuffer(data))
		if err != nil {
			l.Warn("failed client request >%v<", err)
			return nil, err
		}

		err := c.setHeaders(req)
		if err != nil {
			l.Warn("failed setting headers >%v<", err)
			return nil, err
		}

		err = c.setAuthHeaders(req)
		if err != nil {
			l.Warn("failed setting request auth headers >%v<", err)
			return nil, err
		}

		req.Header.Add("Content-Type", "application/json")

		if c.Verbose {
			dmp, err := httputil.DumpRequest(req, true)
			if err != nil {
				l.Warn("failed request dump >%v<", err)
				return nil, err
			}
			l.Info("%s", string(dmp[:]))
		}

		resp, err = client.Do(req)
		if err != nil {
			l.Warn("failed client request >%#v< >%v<", resp, err)
			return resp, err
		}

		if c.Verbose {
			dmp, err := httputil.DumpResponse(resp, true)
			if err != nil {
				l.Warn("failed response dump >%v<", err)
				return nil, err
			}
			l.Info("%s", string(dmp[:]))
		}

	case http.MethodDelete:

		// Post / Put
		if c.Verbose {
			l.Info("Method %s", method)
		} else {
			l.Debug("Method %s", method)
		}

		req, err = http.NewRequest(method, url, bytes.NewBuffer(data))
		if err != nil {
			l.Warn("failed client request >%v<", err)
			return nil, err
		}

		err := c.setHeaders(req)
		if err != nil {
			l.Warn("failed setting headers >%v<", err)
			return nil, err
		}

		err = c.setAuthHeaders(req)
		if err != nil {
			l.Warn("failed setting request auth headers >%v<", err)
			return nil, err
		}

		req.Header.Add("Content-Type", "application/json")

		if c.Verbose {
			dmp, err := httputil.DumpRequest(req, true)
			if err != nil {
				l.Warn("failed request dump >%v<", err)
				return nil, err
			}
			l.Info("%s", string(dmp[:]))
		}

		resp, err = client.Do(req)
		if err != nil {
			l.Warn("failed client request >%#v< >%v<", resp, err)
			return resp, err
		}

		if c.Verbose {
			dmp, err := httputil.DumpResponse(resp, true)
			if err != nil {
				l.Warn("failed response dump >%v<", err)
				return nil, err
			}
			l.Info("%s", string(dmp[:]))
		}

	default:
		// boom
		msg := fmt.Sprintf("method >%s< currently unsupported!", method)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	if c.Verbose {
		l.Info("Client response status >%s<", resp.Status)
	} else {
		l.Debug("Client response status >%s<", resp.Status)
	}

	// Check response code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		err = fmt.Errorf("Response status >%d<", resp.StatusCode)
	}

	return resp, err
}
