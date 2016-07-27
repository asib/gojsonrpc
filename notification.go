package gojsonrpc

import (
	"encoding/json"
	"reflect"
)

type notificationData struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type Notification struct {
	notificationData
}

func (n *Notification) Method() string {
	return n.notificationData.Method
}

func (n *Notification) Params() interface{} {
	return n.notificationData.Params
}

func NotificationValidAndExpectedKeys() map[string]bool {
	return map[string]bool{"jsonrpc": true, "method": true, "params": false}
}

func MakeNotification(method string, params interface{}) (*Notification, error) {
	if params != nil {
		// Params must be either an array/slice or a map with string keys.
		value := reflect.ValueOf(params)
		kind := value.Kind()
		if kind != reflect.Array && kind != reflect.Slice && kind != reflect.Map {
			return nil, invalidNotificationInvalidParamsType
		}

		// For maps, we need to make sure the keys are strings.
		if kind == reflect.Map {
			if reflect.TypeOf(params).Key().Kind() != reflect.String {
				return nil, invalidNotificationInvalidParamsType
			}
		}
	}

	return &Notification{
		notificationData{
			Jsonrpc: version,
			Method:  method,
			Params:  params,
		},
	}, nil
}

func (n *Notification) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.notificationData)
}

func (n *Notification) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &(n.notificationData))
}
