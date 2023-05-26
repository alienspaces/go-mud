package salesforce

import (
	"fmt"
	"net/url"
	"time"

	"github.com/nimajalali/go-force/force"

	"gitlab.com/alienspaces/go-mud/backend/core/model"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

// Salesforce unites your marketing, sales, commerce, service, and IT teams from anywhere with
// Customer 360 — one integrated CRM platform that powers our entire suite of connected apps.
// With Customer 360, you can focus your employees on what’s important right now: stabilising
// your business, reopening, and getting back to delivering exceptional customer experiences.

// Client -
type Client struct {
	Log       logger.Logger
	Client    *force.ForceApi
	Config    ClientConfig
	ConnectTS time.Time
}

type ClientConfig struct {
	Version       string
	Username      string
	Password      string
	SecurityToken string
	ClientID      string
	ClientSecret  string
	Environment   string
}

type SalesForceClient interface {
	Connect() error

	Query(queryString string, results interface{}) error
	Get(path string, params url.Values, results interface{}) error
	Patch(path string, params url.Values, payload, out interface{}) error
	Delete(path string, params url.Values) error

	GetSObject(id string, fields []string, out force.SObject) (err error)
	GetSObjectByExternalID(id string, fields []string, out force.SObject) (err error)
	UpsertByExternalID(id string, in force.SObject) (*force.SObjectResponse, error)
	DeleteByExternalID(id string, in force.SObject) error
	InsertSObject(force.SObject) (*force.SObjectResponse, error)
	DeleteSObject(id string, in force.SObject) error

	Version() string
	Environment() string
}

const (
	DefaultAPIVersion  = "v54.0"
	DefaultEnvironment = "sandbox"
)

// Verify our implementation satisfies our interface
var _ SalesForceClient = &Client{}

// NewClient - Returns a new salesforce API client. If c.IsFaked, the caller should set Client field.
func NewClient(l logger.Logger, c ClientConfig) (*Client, error) {

	scc := &Client{
		Log:    l.WithPackageContext("servicecloud"),
		Config: c,
	}

	if c.Version == "" {
		c.Version = DefaultAPIVersion
	}

	if c.Environment == "" {
		c.Environment = DefaultEnvironment
	}

	err := scc.verifyConfig(c)
	if err != nil {
		l.Warn("NewClient failed verify config >%v<", err)
		return nil, err
	}

	return scc, nil
}

func (scc *Client) Connect() error {
	l := scc.Log.WithFunctionContext("Connect")

	if !scc.shouldReconnect() {
		return nil
	}

	l.Debug("connecting to SalesForce API with client key >%s<", model.TruncateID(scc.Config.ClientID))

	api, err := force.Create(
		scc.Config.Version,
		scc.Config.ClientID,
		scc.Config.ClientSecret,
		scc.Config.Username,
		scc.Config.Password,
		scc.Config.SecurityToken,
		scc.Config.Environment,
	)
	if err != nil {
		l.Warn("Client create failed >%v<", err)
		return err
	}
	scc.Client = api
	scc.ConnectTS = time.Now()

	return nil
}

func (scc *Client) shouldReconnect() bool {
	l := scc.Log.WithFunctionContext("shouldReconnect")

	// If there is no client then we've never connected
	if scc.Client == nil {
		return true
	}

	// If we've been connected for longer that 30 seconds we'll reconnect
	// as we have no way of knowing when the current token will expire.
	tn := time.Now()
	td := tn.Sub(scc.ConnectTS)
	d := time.Duration(30 * time.Second)

	l.Debug("Now >%s< Connected >%s< Expires >%s<", tn.String(), td.String(), d.String())

	return td > d
}

func (scc *Client) Version() string {
	return scc.Config.Version
}

func (scc *Client) Environment() string {
	return scc.Config.Environment
}

func (scc *Client) verifyConfig(cfg ClientConfig) error {
	l := scc.Log.WithFunctionContext("verifyConfig")

	if cfg.ClientID == "" {
		err := fmt.Errorf("salesforce client config missing ClientID")
		l.Warn(err.Error())
		return err
	}

	if cfg.ClientSecret == "" {
		err := fmt.Errorf("salesforce client config missing ClientSecret")
		l.Warn(err.Error())
		return err
	}

	if cfg.Username == "" {
		err := fmt.Errorf("salesforce client config missing Username")
		l.Warn(err.Error())
		return err
	}

	if cfg.Password == "" {
		err := fmt.Errorf("salesforce client config missing Password")
		l.Warn(err.Error())
		return err
	}

	return nil
}
