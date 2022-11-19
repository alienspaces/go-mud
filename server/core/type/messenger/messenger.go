package messenger

// Message
type Message struct {
	ID         string
	GroupID    string
	Subject    string
	Message    string
	Attributes map[string]string
}

// Messenger -
type Messenger interface {
	Publish(topic string, message Message) (messageID string, err error)
	Consume(topic string) (message *Message, err error)
}
