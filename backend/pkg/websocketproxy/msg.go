package websocketproxy

import (
	"encoding/json/jsontext"

	jsonv2 "encoding/json/v2"
)

type Message struct {
	Extra    jsontext.Value         `json:",inline"`
	ID       int                    `json:"id"`
	GUID     string                 `json:"guid"`
	Method   string                 `json:"method,omitzero"`
	Params   map[string]interface{} `json:"params,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Result   interface{}            `json:"result,omitzero"`
	Error    *Error                 `json:"error,omitzero"`
}

func (v Message) MarshalJSON() ([]byte, error) {
	type alias Message
	return jsonv2.Marshal(alias(v))
}

func (v *Message) UnmarshalJSON(data []byte) error {
	type alias Message
	return jsonv2.Unmarshal(data, (*alias)(v))
}

type Error struct {
	Error ErrorPayload `json:"error,omitzero"`
}

type ErrorPayload struct {
	Name    string `json:"name,omitzero"`
	Message string `json:"message,omitzero"`
	Stack   string `json:"stack,omitzero"`
}
