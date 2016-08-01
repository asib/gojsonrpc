package gojsonrpc

import "encoding/json"

type errorData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error is a struct for holding error information that is part of a response.
type Error struct {
	errorData
}

// Code returns the error's code.
func (e *Error) Code() int {
	return e.errorData.Code
}

// Message returns the error's message.
func (e *Error) Message() string {
	return e.errorData.Message
}

// Data returns the error's data field. Callers should use a type assertion to
// recover the variable's type.
func (e *Error) Data() interface{} {
	return e.errorData.Data
}

// Use this function to create Error's - do not try to use a struct literal.
// You may pass nil for the data argument.
func MakeError(code int, message string, data interface{}) *Error {
	return &Error{
		errorData{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

// ErrorValidAndExpectedKeys is a map whose keys are all the possible fields in a
// response error, mapped to whether they are required fields. This mapping is
// used by ParseIncoming to determine whether a response error is valid.
var ErrorValidAndExpectedKeys = map[string]bool{"code": true, "message": true, "data": false}

// This method is used by the encoding/json package when json.Marshal is
// called on an Error struct (or a Response struct). Do not directly call this
// method. Instead, attach the Error to a Response using MakeResponseWithError,
// then run json.Marshal on the Response.
func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.errorData)
}

// This method is used by the encoding/json package when json.Unmarshal is
// called on an Error struct (or a Response struct). You're unlikely to need to
// call this method directly - if you wish to unmarshal a JSON-RPC response,
// use ParseIncoming, then type assert to a Response.
func (e *Error) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &(e.errorData))
}
