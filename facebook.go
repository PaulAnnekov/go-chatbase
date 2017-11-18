package chatbase

import (
	"encoding/json"
	"errors"
	"io"
	"net/url"
)

var (
	facebookMessageEndpoint  = "https://chatbase.com/api/facebook/message_received"
	facebookMessagesEndpoint = "https://chatbase.com/api/facebook/message_received_batch"
	facebookRequestEndpoint  = "https://chatbase.com/api/facebook/send_message"
	facebookRequestsEndpoint = "https://chatbase.com/api/facebook/send_message_batch"
)

// FacebookFields contains metadata about a native Facebook message
type FacebookFields struct {
	Intent     string `json:"intent"`
	NotHandled bool   `json:"not_handled"`
	Feedback   bool   `json:"feedback,omitempty"`
	Version    string `json:"version,omitempty"`
}

// FacebookMessage is a single native Facebook message
type FacebookMessage struct {
	Fields  *FacebookFields
	Payload interface{}
	APIKey  string
}

// MarshalJSON ensures the message is merged with the metadata in the way that
// Chatbase expects it to be
func (f FacebookMessage) MarshalJSON() ([]byte, error) {
	intermediate, intermediateErr := json.Marshal(f.Payload)
	if intermediateErr != nil {
		return nil, intermediateErr
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(intermediate, &m); err != nil {
		return nil, err
	}
	if f.Fields != nil {
		m["chatbase_fields"] = f.Fields
	}
	return json.Marshal(m)
}

// SetIntent adds an optional "intent" value to the message
func (f *FacebookMessage) SetIntent(i string) *FacebookMessage {
	if f.Fields == nil {
		f.Fields = &FacebookFields{}
	}
	f.Fields.Intent = i
	return f
}

// SetNotHandled adds an optional "not handled" value to the message
func (f *FacebookMessage) SetNotHandled(n bool) *FacebookMessage {
	if f.Fields == nil {
		f.Fields = &FacebookFields{}
	}
	f.Fields.NotHandled = n
	return f
}

// SetFeedback adds an optional "feedback" value to the message
func (f *FacebookMessage) SetFeedback(n bool) *FacebookMessage {
	if f.Fields == nil {
		f.Fields = &FacebookFields{}
	}
	f.Fields.Feedback = n
	return f
}

// SetVersion adds an optional "version" value to the message
func (f *FacebookMessage) SetVersion(v string) *FacebookMessage {
	if f.Fields == nil {
		f.Fields = &FacebookFields{}
	}
	f.Fields.Version = v
	return f
}

// Submit tries to deliver a single Facebook message to chatbase
func (f *FacebookMessage) Submit() (*MessageResponse, error) {
	return postSingleFacebookItem(f, f.APIKey, facebookMessageEndpoint)
}

// FacebookMessages is a collection of FacecbookMessage
type FacebookMessages []FacebookMessage

// MarshalJSON ensures the messages are wrapped in a top-level object before
// being serialized into the payload
func (f FacebookMessages) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"messages": []FacebookMessage(f),
	})
}

// Append adds an additional message to the collection. The collection
// cannot contain messages using different API keys
func (f *FacebookMessages) Append(addition *FacebookMessage) *FacebookMessages {
	*f = append(*f, *addition)
	return f
}

// Submit tries to deliver the set of messages to Chatbase. The collection
// cannot contain messages using different API keys
func (f *FacebookMessages) Submit() (*MessagesResponse, error) {
	if len(*f) == 0 {
		return nil, errors.New("cannot submit empty collection")
	}
	apiKey := (*f)[0].APIKey
	return postMultipleFacebookItems(f, apiKey, facebookMessagesEndpoint)
}

func postFacebook(endpoint, apiKey string, v interface{}) (io.ReadCloser, error) {
	u, uErr := url.Parse(endpoint)
	if uErr != nil {
		return nil, uErr
	}
	u.RawQuery = url.Values{"api_key": []string{apiKey}}.Encode()

	body, err := apiPost(u.String(), v)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func postSingleFacebookItem(v interface{}, apiKey, endpoint string) (*MessageResponse, error) {
	return newMessageResponse(func() (io.ReadCloser, error) {
		return postFacebook(endpoint, apiKey, v)
	})
}

func postMultipleFacebookItems(v interface{}, apiKey, endpoint string) (*MessagesResponse, error) {
	return newMessagesResponse(func() (io.ReadCloser, error) {
		return postFacebook(endpoint, apiKey, v)
	})
}

// FacebookRequestResponse is a payload that contains both
// request and response data
type FacebookRequestResponse struct {
	APIKey   string          `json:"-"`
	Request  interface{}     `json:"request_body"`
	Response interface{}     `json:"response_body"`
	Fields   *FacebookFields `json:"chatbase_fields"`
}

// SetIntent adds an optional "intent" value to the pair
func (f *FacebookRequestResponse) SetIntent(i string) *FacebookRequestResponse {
	if f.Fields == nil {
		f.Fields = &FacebookFields{}
	}
	f.Fields.Intent = i
	return f
}

// SetNotHandled adds an optional "not handled" value to the pair
func (f *FacebookRequestResponse) SetNotHandled(n bool) *FacebookRequestResponse {
	if f.Fields == nil {
		f.Fields = &FacebookFields{}
	}
	f.Fields.NotHandled = n
	return f
}

// SetFeedback adds an optional "feedback" value to the pair
func (f *FacebookRequestResponse) SetFeedback(n bool) *FacebookRequestResponse {
	if f.Fields == nil {
		f.Fields = &FacebookFields{}
	}
	f.Fields.Feedback = n
	return f
}

// SetVersion adds an optional "version" value to the pair
func (f *FacebookRequestResponse) SetVersion(v string) *FacebookRequestResponse {
	if f.Fields == nil {
		f.Fields = &FacebookFields{}
	}
	f.Fields.Version = v
	return f
}

// Submit tries to deliver the pair to Chatbase
func (f *FacebookRequestResponse) Submit() (*MessageResponse, error) {
	return postSingleFacebookItem(f, f.APIKey, facebookRequestEndpoint)
}

// FacebookRequestResponses is a collection of FacebookRequestResponse
type FacebookRequestResponses []FacebookRequestResponse

// MarshalJSON ensures the messages are wrapped in a top-level object before
// being serialized into the payload
func (f FacebookRequestResponses) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"messages": []FacebookRequestResponse(f),
	})
}

// Submit tries to send the collection of request/response pairs to Chatbase.
// The collection cannot contain messages using different API keys
func (f *FacebookRequestResponses) Submit() (*MessagesResponse, error) {
	if len(*f) == 0 {
		return nil, errors.New("cannot submit empty collection")
	}
	apiKey := (*f)[0].APIKey
	return postMultipleFacebookItems(f, apiKey, facebookRequestsEndpoint)
}

// Append adds an additional message to the collection. The collection cannot
// contain messages using different API keys
func (f *FacebookRequestResponses) Append(addition *FacebookRequestResponse) *FacebookRequestResponses {
	*f = append(*f, *addition)
	return f
}
