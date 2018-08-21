package types

import (
	"encoding/json"
	"fmt"

	"junolab.net/lib_api/core"
	"junolab.net/ms_core/errors"
)

var emptyJSON = json.RawMessage("{}")

type MsgType string

type Msg struct {
	Id      core.ID          `json:"id"`
	Src     string           `json:"src"`
	Meta    MsgMeta          `json:"-"`
	RawMeta *json.RawMessage `json:"meta"`
	Body    *json.RawMessage `json:"body"`
	Error   errors.MSError   `json:"error"`
}

type TypedMsg struct {
	Msg
	Type MsgType `json:"type,omitempty"`
}

func (msg *Msg) String() string {
	// Represent message as JSON, because we want to have the same message
	// representation as incoming message (we use nats.JSON_ENCODER, so,
	// incoming message always in JSON)
	blob, _ := json.Marshal(msg)
	return string(blob)
}

func (msg *Msg) Success() bool {
	return msg.Error.Code() == ""
}

func NewMsg(id core.ID, service string, encodedMeta []byte, encodedBody []byte) *Msg {
	b := json.RawMessage(encodedBody)
	m := json.RawMessage(encodedMeta)
	return &Msg{
		Id:      id,
		Src:     service,
		RawMeta: &m,
		Body:    &b,
		Error:   errors.NoError,
	}
}

func NewTypedMsg(id core.ID, service string, encodedMeta []byte, encodedBody []byte, msgType MsgType) *TypedMsg {
	m := NewMsg(id, service, encodedMeta, encodedBody)
	return &TypedMsg{
		Msg:  *m,
		Type: msgType,
	}
}

func NewErrMsg(id core.ID, service string, err errors.Error) *Msg {
	return &Msg{
		Id:      id,
		Src:     service,
		RawMeta: &emptyJSON,
		Body:    &emptyJSON,
		Error:   err.(errors.MSError),
	}
}

// Specific JSON marshalling for MsgMeta just to make AuthInfo empty if no UserId in it
// so JSON omitempty will treat it as empty field
// func (t MsgMeta) MarshalJSON() ([]byte, error) {
// 	if t.Auth().UserId == "" {
// 		t.AuthInfo = nil
// 	}
// 	type jtype MsgMeta // Type alias for the recursive call prevention
// 	return json.Marshal(jtype(t))
// }

type StreamingResponse struct {
	Subject  string         `json:"subject"`
	Payload  *HTTPPayload   `json:"payload,omitempty"`  // First message for the stream if needed
	Payloads []*HTTPPayload `json:"payloads,omitempty"` // More initial messages for the stream if needed TODO: adding as separate field for backwards compatibility (code and ms api level) - update struct, change occurs
}

type CancelStreamingRequest struct {
	Subject string `json:"subject"`
}

type HTTPPayload struct {
	Body json.RawMessage
}

func (p HTTPPayload) MarshalJSON() ([]byte, error) {
	if len(p.Body) == 0 {
		return []byte("null"), nil
	}
	return p.Body, nil
}

func (p *HTTPPayload) UnmarshalJSON(data []byte) error {
	p.Body = json.RawMessage(data)
	return nil
}

func (p HTTPPayload) String() string {
	return string(p.Body)
}

func (r StreamingResponse) Validate() error {
	if len(r.Subject) > 256 {
		return fmt.Errorf("field Subject in StreamingResponse is more than 256")
	}
	if r.Payload != nil {
		if err := r.Payload.Validate(); err != nil {
			return err
		}
	}
	for _, x := range r.Payloads {
		if err := x.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (r *HTTPPayload) Validate() error {
	if len(r.Body) < 1 {
		return fmt.Errorf("field HTTPPayload in StreamingResponse is less than 1")
	}
	return nil
}

func (r *HTTPPayload) ValidateJSON() error {
	return r.Validate()
}

func (r *CancelStreamingRequest) validate() error {
	if len(r.Subject) < 1 {
		return fmt.Errorf("field Subject is less then 1 ")
	}
	if len(r.Subject) > 256 {
		return fmt.Errorf("field Subject is more then 256 ")
	}
	return nil
}

func (r *CancelStreamingRequest) Validate() error {
	return r.validate()
}

func (r *CancelStreamingRequest) ValidateJSON() error {
	return r.validate()
}
