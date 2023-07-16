package message

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/messenger"
)

const (
	packageName = "message"
)

// Client -
type Client struct {
	log       logger.Logger
	config    Configuration
	snsClient *sns.SNS
}

type Configuration struct {
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSAccessToken     string
}

var _ messenger.Messenger = &Client{}

// NewClient -
func NewClient(l logger.Logger, c Configuration) (*Client, error) {

	m := &Client{
		log:    l,
		config: c,
	}

	err := m.Init()
	if err != nil {
		m.log.Warn("Failed init >%v<", err)
		return nil, err
	}

	return m, nil
}

// Init -
func (m *Client) Init() error {
	l := m.logger("getAWSSession")

	// sns client
	snsClient, err := m.getSNSClient()
	if err != nil {
		l.Warn("Failed to get sns client >%v<", err)
		return err
	}
	m.snsClient = snsClient

	// TODO: Verify configuration

	return nil
}

// getAWSSession
func (m *Client) getAWSSession() (*session.Session, error) {
	l := m.logger("getAWSSession")

	awsRegion := m.config.AWSRegion
	awsAccessKeyID := m.config.AWSAccessKeyID
	awsSecretAccessKey := m.config.AWSSecretAccessKey
	awsAccessToken := m.config.AWSAccessToken

	l.Info("=== AWS access key ID >%s<", awsAccessKeyID)
	l.Info("=== AWS session access key >%s<", awsSecretAccessKey)
	l.Info("=== AWS access token >%s<", awsAccessToken)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, awsAccessToken),
	})
	if err != nil {
		l.Warn("Failed new AWS session >%v<", err)
		return nil, err
	}

	return sess, nil
}

func (m *Client) logger(functionName string) logger.Logger {
	if m.log == nil {
		return nil
	}
	return m.log.WithPackageContext(packageName).WithFunctionContext(functionName)
}
