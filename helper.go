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
	incomingMap := make(map[string]json.RawMessage)
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
	} else if isErrorResponse(keys) || isResultResponse(keys) {
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

func areKeySetsMatching(existingKeys []string, expectedKeys map[string]bool) bool {
	// Make sure all required keys exist.
	for k, expected := range expectedKeys {
		// If this key is expected, look for it in existingKeys
		if expected {
			found := false
			for _, existing := range existingKeys {
				if k == existing {
					found = true
					break
				}
			}
			// If we didn't find it, then we already know keysets aren't matching,
			// so we return false.
			if !found {
				return false
			}
		}
	}

	// Make sure we don't have any extra, unexpected keys.
	for _, existing := range existingKeys {
		// If there is a mapping whose key is `existing`, then regardless of whether
		// it is required or not, it's valid, so we allow it. Only if the key is
		// invalid (i.e. not required and not even an optional field) do we return
		// false.
		if _, ok := expectedKeys[existing]; !ok { // here we are checking if the key `existing` is in the map `expectedKeys`
			return false
		}
	}

	// If we managed to reach this point, then the keysets are matching.
	return true
}

func isNotification(keys []string) bool {
	return areKeySetsMatching(keys, NotificationValidAndExpectedKeys())
}

func isRequest(keys []string) bool {
	return areKeySetsMatching(keys, RequestValidAndExpectedKeys())
}

func isErrorResponse(keys []string) bool {
	return areKeySetsMatching(keys, ErrorResponseValidAndExpectedKeys())
}

func isResultResponse(keys []string) bool {
	return areKeySetsMatching(keys, ResultResponseValidAndExpectedKeys())
}
