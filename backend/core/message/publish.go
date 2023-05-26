package message

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"

	"gitlab.com/alienspaces/go-mud/backend/core/type/messenger"
)

// Publish -
func (m *Client) Publish(topicARN string, message messenger.Message) (messageID string, err error) {
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

	l.Info("publishing message topic >%s< message >%s< attributes >%#v<", topicARN, message.Message, awsMessageAttributes)

	input := &sns.PublishInput{
		Message:           aws.String(message.Message),
		Subject:           aws.String(message.Subject),
		TopicArn:          aws.String(topicARN),
		MessageAttributes: awsMessageAttributes,
	}

	if strings.HasSuffix(topicARN, ".fifo") {
		input.MessageDeduplicationId = aws.String(message.ID)
		input.MessageGroupId = aws.String(message.GroupID)
	}

	result, err := m.snsClient.Publish(input)
	if err != nil {
		l.Warn("failed publishing message >%v<", err)
		return "", err
	}

	// No message ID, failed send
	if *result.MessageId == "" {
		msg := "response does not contain message ID"
		l.Warn(msg)
		return "", fmt.Errorf(msg)
	}

	l.Info("publish result >%+v<", result)

	return *result.MessageId, nil
}

// getSNSClient
func (m *Client) getSNSClient() (*sns.SNS, error) {

	// session
	awsSession, err := m.getAWSSession()
	if err != nil {
		m.log.Warn("Failed to get AWS session >%v<", err)
		return nil, err
	}

	// client
	snsClient := sns.New(awsSession)

	return snsClient, nil
}
