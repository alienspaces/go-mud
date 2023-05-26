package message

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/google/uuid"

	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/messenger"
)

type FakedClient struct {
	Client
}

// NewFakeClient -
func NewFakeClient(l logger.Logger) (*FakedClient, error) {

	client := &FakedClient{
		Client: Client{
			log: l,
		},
	}

	return client, nil
}

func (m *FakedClient) Publish(topicARN string, message messenger.Message) (messageID string, err error) {
	l := m.logger("Publish")

	awsMessageAttributes := map[string]*sns.MessageAttributeValue{}

	for k := range message.Attributes {
		v := message.Attributes[k]
		if v == "" {
			l.Debug("excluding message attribute k >%s< v >%s<, value is empty", k, v)
			continue
		}
		l.Debug("adding message attribute k >%s< v >%s<", k, v)
		awsMessageAttributes[k] = &sns.MessageAttributeValue{StringValue: &v, DataType: aws.String("String")}
	}

	l.Info("publishing message topic >%s< message length >%d<", topicARN, len(message.Message))

	uuidByte, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uuidByte.String(), nil
}

func (m *FakedClient) Consume(queueARN string) (message *messenger.Message, err error) {
	return message, nil
}

func (m *FakedClient) logger(functionName string) logger.Logger {
	if m.log == nil {
		return nil
	}
	return m.log.WithPackageContext(fmt.Sprintf("(fake) %s", packageName)).WithFunctionContext(functionName)
}
