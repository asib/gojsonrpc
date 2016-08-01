// Package gojsonrpc provides an interface for dealing with JSON-RPC APIs.
package gojsonrpc

import (
	"encoding/json"
	"reflect"
)

// ParseIncoming attempts to parse the supplied message into one of the three
// relevant types: Notification, Request or Response. Callers must run a type
// assertion to identify which type was returned.
func ParseIncoming(message string) (Message, error) {
	var incomingMap map[string]json.RawMessage
	err := json.Unmarshal([]byte(message), &incomingMap)
	if err != nil {
		return nil, err
	}

	// Look for jsonrpc field, return error if not present.
	if _, ok := incomingMap[VersionKey]; !ok {
		return nil, InvalidMessage
	}

	// Check version is correct.
	var incomingVersion string
	err = json.Unmarshal(incomingMap[VersionKey], &incomingVersion)
	if err != nil {
		return nil, err
	} else if incomingVersion != Version {
		return nil, InvalidVersion
	}

	// Now we need to try to match this object's keys against those of a
	// notification, request or response.
	var keys []string
	for k := range incomingMap {
		keys = append(keys, k)
	}

	if isNotification(keys) {
		return parseIncomingNotification([]byte(message))
	} else if isRequest(keys) {
		return parseIncomingRequest([]byte(message))
	} else if isErrorResponse(keys) {
		// Check that the error is valid
		var errorMap map[string]interface{}
		if err := json.Unmarshal(incomingMap[ErrorKey], &errorMap); err != nil {
			return nil, err
		}

		var errKeys []string
		for k := range errorMap {
			errKeys = append(errKeys, k)
		}
		if isValidResponseError(errKeys) {
			return parseIncomingResponse([]byte(message))
		}
	} else if isResultResponse(keys) {
		return parseIncomingResponse([]byte(message))
	}

	// If not caught by one of the above, must be an malformed message.
	return nil, InvalidMessage
}

func parseIncomingNotification(jsonNotif []byte) (*Notification, error) {
	notif := new(Notification)
	if err := json.Unmarshal(jsonNotif, notif); err != nil {
		return nil, err
	}

	if notif.Params() != nil {
		// Params must be a JSON array or object.
		value := reflect.ValueOf(notif.Params())
		kind := value.Kind()
		if kind != reflect.Slice && kind != reflect.Map {
			return nil, InvalidMessage
		}
	}

	// No need to check that the map is of type map[string]interface{} because
	// the encoding/json package does this for us.

	return notif, nil
}

func parseIncomingRequest(jsonReq []byte) (*Request, error) {
	req := new(Request)
	if err := json.Unmarshal(jsonReq, req); err != nil {
		return nil, err
	}

	if req.Params() != nil {
		// Params must be a JSON array or object.
		value := reflect.ValueOf(req.Params())
		kind := value.Kind()
		if kind != reflect.Slice && kind != reflect.Map {
			return nil, InvalidMessage
		}
	}

	// No need to check that the map is of type map[string]interface{} because
	// the encoding/json package does this for us.

	return req, nil
}

func parseIncomingResponse(jsonResp []byte) (*Response, error) {
	resp := new(Response)
	if err := json.Unmarshal(jsonResp, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func isNotification(keys []string) bool {
	return AreKeySetsMatching(keys, NotificationValidAndExpectedKeys)
}

func isRequest(keys []string) bool {
	return AreKeySetsMatching(keys, RequestValidAndExpectedKeys)
}

func isErrorResponse(keys []string) bool {
	return AreKeySetsMatching(keys, ErrorResponseValidAndExpectedKeys)
}

func isValidResponseError(keys []string) bool {
	return AreKeySetsMatching(keys, ErrorValidAndExpectedKeys)
}

func isResultResponse(keys []string) bool {
	return AreKeySetsMatching(keys, ResultResponseValidAndExpectedKeys)
}
