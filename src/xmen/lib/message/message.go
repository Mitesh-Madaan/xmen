package message

import (
	"fmt"
)

type Message struct {
	// Message
	Reason string
	// ReasonCategory
	ReasonCategory string
}

func NewMessage(reason string, reason_category string) *Message {
	// Create a new message
	return &Message{Reason: reason, ReasonCategory: reason_category}
}

func (m *Message) toStatus() string {
	// Convert the message to a status
	return fmt.Sprintf("%s: %s", m.ReasonCategory, m.Reason)
}

func DumpMessages(messages []*Message) string {
	// Dump the messages
	data := ""
	if messages == nil {
		return data
	}
	for _, message := range messages {
		data += fmt.Sprintf("%s; ", message.toStatus())
	}
	return data
}
