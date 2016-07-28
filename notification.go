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

// Notification is a struct for holding information about a notification (a
// JSON-RPC request without an ID).
type Notification struct {
	notificationData
}

// Method returns the notification's method.
func (n *Notification) Method() string {
	return n.notificationData.Method
}

// Params returns the notification's params - this may be nil. Use type assertion
// to recover the variable's type.
func (n *Notification) Params() interface{} {
	return n.notificationData.Params
}

// This function returns a map whose keys are all the possible fields in a
// notification, mapped to whether they are required fields (e.g. params is not
// a required field for a notification, so it maps to false). This mapping is
// used by ParseIncoming to determine the type of the incoming message.
func NotificationValidAndExpectedKeys() map[string]bool {
	return map[string]bool{"jsonrpc": true, "method": true, "params": false}
}

// MakeNotification is used to create Notification structs - do not try to create
// Notification's using a struct literal. You may pass nil for the params argument.
// Else, params MUST be an array/slice or a map with string keys.
func MakeNotification(method string, params interface{}) (*Notification, error) {
	if params != nil {
		// Params must be either an array/slice or a map with string keys.
		value := reflect.ValueOf(params)
		kind := value.Kind()
		if kind != reflect.Array &&
			kind != reflect.Slice &&
			kind != reflect.Map &&
			kind != reflect.Struct &&
			kind != reflect.Ptr { // could be a pointer to one of the above
			return nil, InvalidNotificationInvalidParamsType
		} else if kind == reflect.Map && reflect.TypeOf(params).Key().Kind() != reflect.String { // For maps, we need to make sure the keys are strings.
			return nil, InvalidNotificationInvalidParamsType
		} else if kind == reflect.Ptr {
			elemKind := reflect.TypeOf(params).Elem().Kind()
			if elemKind != reflect.Array &&
				elemKind != reflect.Struct {
				return nil, InvalidNotificationInvalidParamsType
			}
		}
	}

	return &Notification{
		notificationData{
			Jsonrpc: Version,
			Method:  method,
			Params:  params,
		},
	}, nil
}

// Do not used this method directly. Instead, call json.Marshal with a
// Notification as the argument.
func (n *Notification) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.notificationData)
}

// Do not use this method directly. Instead, use ParseIncoming and type assert
// the returned value to a Notification.
func (n *Notification) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &(n.notificationData))
}
