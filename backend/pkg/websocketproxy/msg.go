package websocketproxy

import "github.com/mailru/easyjson"

//go:generate easyjson
//easyjson:json
type Message struct {
	easyjson.UnknownFieldsProxy
	ID       int                    `json:"id"`
	GUID     string                 `json:"guid"`
	Method   string                 `json:"method,omitempty"`
	Params   map[string]interface{} `json:"params,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Result   interface{}            `json:"result,omitempty"`
	Error    *Error                 `json:"error,omitempty"`
}

type Error struct {
	Error ErrorPayload `json:"error,omitempty"`
}

type ErrorPayload struct {
	Name    string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
	Stack   string `json:"stack,omitempty"`
}
