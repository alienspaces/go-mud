package smtp

import (
	"bytes"
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/type/emailer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

// FakeSMTP -
type FakeSMTP struct {
	SMTP
}

var _ emailer.Emailer = &FakeSMTP{}

// NewFakeSMTP -
func NewFakeSMTP(l logger.Logger, c Configuration) (*FakeSMTP, error) {

	e := &FakeSMTP{
		SMTP: SMTP{
			config: c,
			log:    l,
		},
	}

	err := e.Init()
	if err != nil {
		l.Warn("failed fake smtp init >%v<", err)
		return nil, err
	}

	return e, nil
}

// Init -
func (e *FakeSMTP) Init() error {
	l := e.logger(" Connect")
	l.Info("Initialising")

	host := e.config.Host
	if host == "" {
		err := fmt.Errorf("missing host, cannot init initialise emailer")
		l.Warn(err.Error())
		return err
	}

	// SMTP config
	e.host = host

	return nil
}

// Connect -
func (e *FakeSMTP) Connect() error {
	return nil
}

// SendEmail -
func (e *FakeSMTP) SendEmail(from, to, subject, content string) error {
	l := e.logger("SendEmail")
	l.Info("Sending text email host >%s< from >%s< to >%s< subject >%s< content length >%d<", e.host, from, to, subject, len(content))
	return nil
}

// SendHTMLEmail -
func (e *FakeSMTP) SendHTMLEmail(from, to, subject, content string) error {
	l := e.logger("SendHTMLEmail")
	l.Info("Sending html email host >%s< from >%s< to >%s< subject >%s< content length >%d<", e.host, from, to, subject, len(content))
	return nil
}

// SendRawEmail - Support attachment
func (e *FakeSMTP) SendRawEmail(from string, to string, content bytes.Buffer) error {
	l := e.logger("SendRawEmail")
	l.Info("Sending raw email host >%s< from >%s< to >%s< content length >%d<", e.host, from, to, content.Len())
	return nil
}

// Quit the connection
func (e *FakeSMTP) Quit() error {
	return nil
}

func (e *FakeSMTP) logger(functionName string) logger.Logger {
	if e.log == nil {
		return nil
	}
	return e.log.WithPackageContext(fmt.Sprintf("(fake) %s", packageName)).WithFunctionContext(functionName)
}
