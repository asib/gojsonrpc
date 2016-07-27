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

type Request struct {
	requestData
}

func (r *Request) Method() string {
	return r.requestData.Method
}

func (r *Request) Params() interface{} {
	return r.requestData.Params
}

func (r *Request) ID() uint {
	return r.requestData.ID
}

func RequestValidAndExpectedKeys() map[string]bool {
	return map[string]bool{"jsonrpc": true, "method": true, "params": false, "id": true}
}

// Unexported, so we force users of the package to use MakeRequestWithArray and
// MakeRequestWithMap, which ensures params field is of the expected type.
func MakeRequest(method string, params interface{}, id uint) (*Request, error) {
	if params != nil {
		// Params must be either an array/slice or a map with string keys.
		value := reflect.ValueOf(params)
		kind := value.Kind()
		if kind != reflect.Array && kind != reflect.Slice && kind != reflect.Map {
			return nil, invalidRequestInvalidParamsType
		}

		if kind == reflect.Map {
			if reflect.TypeOf(params).Key().Kind() != reflect.String {
				return nil, invalidRequestInvalidParamsType
			}
		}
	}

	return &Request{
		requestData{
			Jsonrpc: version,
			Method:  method,
			Params:  params,
			ID:      id,
		},
	}, nil
}

func (r *Request) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.requestData)
}

func (r *Request) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &(r.requestData))
}
