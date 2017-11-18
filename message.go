package chatbase

import (
	"encoding/json"
	"io"
)

// MessageType describes the source of a message
type MessageType string

// types of messages used to identify the origin of
// a message in Chatbase
const (
	UserType  MessageType = "user"
	AgentType MessageType = "agent"
)

var (
	messagesEndpoint = "https://chatbase.com/api/messages"
	messageEndpoint  = "https://chatbase.com/api/message"
)

// Message contains data about a generic chat message
type Message struct {
	APIKey     string      `json:"api_key"`
	Type       MessageType `json:"type"`
	UserID     string      `json:"user_id"`
	TimeStamp  int         `json:"time_stamp"`
	Platform   string      `json:"platform"`
	Message    string      `json:"message,omitempty"`
	Intent     string      `json:"intent,omitempty"`
	NotHandled bool        `json:"not_handled,omitempty"`
	Feedback   bool        `json:"feedback,omitempty"`
	Version    string      `json:"version,omitempty"`
}

// SetMessage adds an optional "message" value to a message
func (m *Message) SetMessage(msg string) *Message {
	m.Message = msg
	return m
}

// SetIntent adds an optional "intent" value to a message
func (m *Message) SetIntent(i string) *Message {
	m.Intent = i
	return m
}

// SetNotHandled adds an optional "not handled" value to a message
func (m *Message) SetNotHandled(n bool) *Message {
	m.NotHandled = n
	return m
}

// SetFeedback adds an optional "feedback" value to a message
func (m *Message) SetFeedback(f bool) *Message {
	m.Feedback = f
	return m
}

// SetVersion adds an optional "version" value to a message
func (m *Message) SetVersion(v string) *Message {
	m.Version = v
	return m
}

// SetTimeStamp overrides the message's "timestamp" value
func (m *Message) SetTimeStamp(t int) *Message {
	m.TimeStamp = t
	return m
}

// Submit tries to deliver the message to Chatbase
func (m *Message) Submit() (*MessageResponse, error) {
	return newMessageResponse(func() (io.ReadCloser, error) {
		return apiPost(messageEndpoint, m)
	})
}

// MessageResponse describes a Chatbase response to the submission of
// a single message. It is also used to represent the result of an item
// of a collection of messages that have been submitted.
type MessageResponse struct {
	MessageID MessageID `json:"message_id"`
	Status    Status    `json:"status"`
	Reason    string    `json:"reason,omitempty"`
}

// Messages is a collection of Message
type Messages []Message

// MarshalJSON ensures the list of messages is wrapped in a
// top-level object before being serialized into JSON
func (m Messages) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"messages": []Message(m),
	})
}

// Submit tries to deliver the set of messages to Chatbase
func (m *Messages) Submit() (*MessagesResponse, error) {
	return newMessagesResponse(func() (io.ReadCloser, error) {
		return apiPost(messagesEndpoint, m)
	})
}

// Append adds a message to the the collection
func (m *Messages) Append(addition ...*Message) *Messages {
	for _, a := range addition {
		*m = append(*m, *a)
	}
	return m
}

// MessagesResponse describes a Chatbase response to submitting a set of messages
type MessagesResponse struct {
	AllSucceeded bool              `json:"all_succeeded"`
	Status       Status            `json:"status"`
	Responses    []MessageResponse `json:"responses"`
	Reason       string            `json:"reason,omitempty"`
}
