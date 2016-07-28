package gojsonrpc

import (
	"encoding/json"
	"reflect"
)

type requestData struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      uint        `json:"id"`
}

// Request is a struct that holds information about a request.
type Request struct {
	requestData
}

// JSONRPCVersion returns the version of the protocol being used.
func (n *Request) JSONRPCVersion() string {
	return n.requestData.Jsonrpc
}

// Method returns the request's method.
func (r *Request) Method() string {
	return r.requestData.Method
}

// Params returns the request's params - this may be nil. Use a type assertion
// to recover the variable's type.
func (r *Request) Params() interface{} {
	return r.requestData.Params
}

// ID returns the request's ID.
func (r *Request) ID() uint {
	return r.requestData.ID
}

// This function returns a map whose keys are all the possible fields in a
// request, mapped to whether they are required fields (e.g. params is not a
// required field for a request, so it maps to false). This mapping is used by
// ParseIncoming to determine the type of the incoming message.
func RequestValidAndExpectedKeys() map[string]bool {
	return map[string]bool{"jsonrpc": true, "method": true, "params": false, "id": true}
}

// MakeRequest is used to create Request structs - do not try to use a struct
// literal. You may pass nil for the params argument. Else, params must be an
// array/slice or a map with string keys.
func MakeRequest(method string, params interface{}, id uint) (*Request, error) {
	if params != nil {
		// Params must be either an array/slice or a map with string keys.
		value := reflect.ValueOf(params)
		kind := value.Kind()
		if kind != reflect.Array &&
			kind != reflect.Slice &&
			kind != reflect.Map &&
			kind != reflect.Struct &&
			kind != reflect.Ptr {
			return nil, InvalidRequestInvalidParamsType
		} else if kind == reflect.Map && reflect.TypeOf(params).Key().Kind() != reflect.String {
			return nil, InvalidRequestInvalidParamsType
		} else if kind == reflect.Ptr {
			elemKind := reflect.TypeOf(params).Elem().Kind()
			if elemKind != reflect.Array &&
				elemKind != reflect.Struct {
				return nil, InvalidRequestInvalidParamsType
			}
		}
	}

	return &Request{
		requestData{
			Jsonrpc: Version,
			Method:  method,
			Params:  params,
			ID:      id,
		},
	}, nil
}

// Do not use this method directly. Instead, use json.Marshal with a Request
// as the argument.
func (r *Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.requestData)
}

// Do not use this method directly. Instead, use ParseIncoming and type assert
// the returned value to a Request.
func (r *Request) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &(r.requestData))
}
