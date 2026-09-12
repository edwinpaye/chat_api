package ws

import "encoding/json"

// Message is the wire envelope used both directions.
type Message struct {
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func OKReply(id, typ string, payload any) Message {
	raw, _ := json.Marshal(payload)
	return Message{ID: id, Type: typ + ".ok", Payload: raw}
}

func ErrReply(id, typ, code string, err error) Message {
	raw, _ := json.Marshal(ErrorPayload{Code: code, Message: err.Error()})
	return Message{ID: id, Type: typ + ".error", Payload: raw}
}