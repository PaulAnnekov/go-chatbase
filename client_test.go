package chatbase

import (
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		c := NewClient("foo-bar-baz")
		if c.String() != "foo-bar-baz" {
			t.Errorf("Expected foo-bar-baz, got %v", c.String())
		}
	})
}

func TestMessage_Client(t *testing.T) {
	oldTimeStamp := TimeStamp
	defer func() { TimeStamp = oldTimeStamp }()
	TimeStamp = func() int { return 998877 }
	c := NewClient("foo-bar-baz")

	t.Run("default", func(t *testing.T) {
		expected := &Message{
			APIKey:    "foo-bar-baz",
			Type:      "agent",
			UserID:    "abc123",
			TimeStamp: 998877,
			Platform:  "fantasy-chat",
		}
		m := c.Message(MessageTypeAgent, "abc123", "fantasy-chat")
		if !reflect.DeepEqual(expected, m) {
			t.Errorf("Expected %v, got %v", expected, m)
		}
	})
	t.Run("agent message", func(t *testing.T) {
		expected := &Message{
			APIKey:    "foo-bar-baz",
			Type:      "agent",
			UserID:    "abc123",
			TimeStamp: 998877,
			Platform:  "fantasy-chat",
		}
		m := c.AgentMessage("abc123", "fantasy-chat")
		if !reflect.DeepEqual(expected, m) {
			t.Errorf("Expected %v, got %v", expected, m)
		}
	})
	t.Run("user message", func(t *testing.T) {
		expected := &Message{
			APIKey:    "foo-bar-baz",
			Type:      "user",
			UserID:    "abc123",
			TimeStamp: 998877,
			Platform:  "fantasy-chat",
		}
		m := c.UserMessage("abc123", "fantasy-chat")
		if !reflect.DeepEqual(expected, m) {
			t.Errorf("Expected %v, got %v", expected, m)
		}
	})
}

func TestEvent_Client(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		c := NewClient("foo-bar-baz")
		expected := &Event{
			APIKey: "foo-bar-baz",
			UserID: "abc-123",
			Intent: "test-things",
		}
		e := c.Event("abc-123", "test-things")
		if !reflect.DeepEqual(expected, e) {
			t.Errorf("Expected %v, got %v", expected, e)
		}
	})
}

func TestEvent_Update(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		c := NewClient("foo-bar-baz")
		expected := &Update{
			APIKey:    "foo-bar-baz",
			MessageID: "abc123",
		}
		u := c.Update("abc123")
		if !reflect.DeepEqual(expected, u) {
			t.Errorf("Expected %v, got %v", expected, u)
		}
	})
}
