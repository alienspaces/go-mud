package message

import (
	"gitlab.com/alienspaces/go-mud/backend/core/type/messenger"
)

// Consume -
func (m *Client) Consume(queueARN string) (message *messenger.Message, err error) {
	return message, nil
}
