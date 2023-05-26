package smtp

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"gitlab.com/alienspaces/go-mud/backend/core/type/emailer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

const (
	packageName = "smtp"
	// MimeHTML - Content type
	MimeHTML = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

// SMTP -
type SMTP struct {
	log        logger.Logger
	config     Configuration
	connection *smtp.Client
	host       string
}

type Configuration struct {
	Host string
}

var _ emailer.Emailer = &SMTP{}

// NewSMTP -
func NewSMTP(l logger.Logger, c Configuration) (*SMTP, error) {

	e := &SMTP{
		config:     c,
		log:        l,
		connection: &smtp.Client{},
	}

	err := e.Init()
	if err != nil {
		l.Warn("failed email init >%v<", err)
		return nil, err
	}

	return e, nil
}

// Init -
func (e *SMTP) Init() error {
	l := e.logger("Connect")
	l.Info("Initialising")

	host := e.config.Host
	if host == "" {
		err := fmt.Errorf("missing host, cannot init initialise emailer")
		l.Warn(err.Error())
		return err
	}

	// smtp config
	e.host = host

	l.Info("SMPT host >%v<", e.host)

	err := e.Connect()
	if err != nil {
		l.Warn("failed connecting to host >%s< >%v<", e.host, err)
		return err
	}

	return nil
}

// Connect -
func (e *SMTP) Connect() error {
	l := e.logger("Connect")
	l.Info("Connecting")

	// Connect to the remote SMTP server.
	c, err := smtp.Dial(e.host)
	if err != nil {
		l.Warn("failed dialing host >%v<", err)
		return err
	}

	e.connection = c
	return nil
}

// SendEmail -
func (e *SMTP) SendEmail(from, to, subject, content string) error {
	l := e.logger("SendEmail")
	l.Info("Sending text email host >%s< from >%s< to >%s< subject >%s< content length >%d<", e.host, from, to, subject, len(content))

	err := e.Connect()
	if err != nil {
		l.Warn("failed connecting to host >%s< >%v<", e.host, err)
		return err
	}

	c := e.connection

	// Set the sender and recipient first
	if err := c.Mail(from); err != nil {
		l.Warn("failed setting from >%s< >%v<", from, err)
		return err
	}

	tos := strings.Split(to, ",")
	for _, t := range tos {
		if err := c.Rcpt(t); err != nil {
			l.Warn("failed setting recipient >%s< >%v<", t, err)
			return err
		}
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		l.Warn("failed creating client data writer >%v<", err)
		return err
	}

	// Write the headers
	header := []byte("To: " + to + "\r\n" + "Subject: " + subject + "\r\n" + "\r\n")

	_, err = wc.Write(header)
	if err != nil {
		l.Warn("failed writing header >%v<", err)
		return err
	}

	_, err = fmt.Fprint(wc, content)
	if err != nil {
		l.Warn("failed writing body >%v<", err)
		return err
	}

	err = wc.Close()
	if err != nil {
		l.Warn("failed closing client data writer >%v<", err)
		return err
	}

	return nil
}

// SendHTMLEmail -
func (e *SMTP) SendHTMLEmail(from, to, subject, content string) error {
	l := e.logger("SendHTMLEmail")
	l.Info("Sending HTML email host >%s< from >%s< to >%s< subject >%s< content length >%d<", e.host, from, to, subject, len(content))

	err := e.Connect()
	if err != nil {
		l.Warn("failed connecting to host >%s< >%v<", e.host, err)
		return err
	}

	c := e.connection

	// Set the sender and recipient first
	if err := c.Mail(from); err != nil {
		l.Warn("failed setting from >%s< >%v<", from, err)
		return err
	}

	tos := strings.Split(to, ",")
	for _, t := range tos {
		if err := c.Rcpt(t); err != nil {
			l.Warn("failed setting recipient >%s< >%v<", t, err)
			return err
		}
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		l.Warn("failed creating client data writer >%v<", err)
		return err
	}

	// Write the headers
	header := []byte("To: " + to + "\r\n" + "Subject: " + subject + "\r\n" + MimeHTML + "\r\n")

	_, err = wc.Write(header)
	if err != nil {
		l.Warn("failed writing header >%v<", err)
		return err
	}

	_, err = fmt.Fprint(wc, content)
	if err != nil {
		l.Warn("failed writing body >%v<", err)
		return err
	}

	err = wc.Close()
	if err != nil {
		l.Warn("failed closing client data writer >%v<", err)
		return err
	}

	return nil
}

// SendRawEmail - Support attachment
func (e *SMTP) SendRawEmail(from string, to string, content bytes.Buffer) error {
	l := e.logger("SendRawEmail")
	l.Info("Sending raw host >%s< from >%s< to >%s< content length >%d<", e.host, from, to, content.Len())

	err := e.Connect()
	if err != nil {
		l.Warn("failed connecting to host >%s< >%v<", e.host, err)
		return err
	}

	c := e.connection

	// Set the sender and recipient first
	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	_, err = wc.Write(content.Bytes())
	if err != nil {
		return err
	}

	err = wc.Close()

	if err != nil {
		return err
	}

	return err
}

// Quit the connection
func (e *SMTP) Quit() error {
	l := e.logger("Quit")
	l.Info("Quitting")

	// if not connected return nil
	if e.connection.Text == nil {
		return nil
	}
	err := e.connection.Quit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (e *SMTP) logger(functionName string) logger.Logger {
	if e.log == nil {
		return nil
	}
	return e.log.WithPackageContext(packageName).WithFunctionContext(functionName)
}
