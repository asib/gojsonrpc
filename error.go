package gojsonrpc

import "encoding/json"

type errorData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Error struct {
	errorData
}

func (e *Error) Code() int {
	return e.errorData.Code
}

func (e *Error) Message() string {
	return e.errorData.Message
}

func (e *Error) Data() interface{} {
	return e.errorData.Data
}

func MakeError(code int, message string, data interface{}) *Error {
	return &Error{
		errorData{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.errorData)
}

func (e *Error) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &(e.errorData))
}
