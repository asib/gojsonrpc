package gojsonrpc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var (
	testNotificationMethod  = "test"
	testNotificationParams  = map[string]interface{}{"key1": "value1", "key2": float64(2), "key3": false}
	testNotificationParams2 = []interface{}{"val1", 2, true, 34.21}
)

func TestCreateNotificationWithoutParams(t *testing.T) {
	n, err := MakeNotification(testNotificationMethod, nil)
	if err != nil {
		t.Fatal(err)
	}
	if n == nil {
		t.Error("notification should not be nil")
	}
	if n.Method() != testNotificationMethod {
		t.Error("method not set correctly")
	}
	if n.Params() != nil {
		t.Error("params should be nil")
	}
}

func TestCreateNotificationWithParams(t *testing.T) {
	n, err := MakeNotification(testNotificationMethod, testNotificationParams)
	if err != nil {
		t.Fatal(err)
	}
	if n == nil {
		t.Error("notification should not be nil")
	}
	if n.Method() != testNotificationMethod {
		t.Error("method not set correctly")
	}
	if n.Params() == nil {
		t.Error("params should not be nil")
	} else {
		if reflect.ValueOf(n.Params()).Kind() != reflect.Map {
			t.Error("params should be map")
		} else {
			mapParams := n.Params().(map[string]interface{})
			for k, v := range testNotificationParams {
				if !mapContains(mapParams, k, v) {
					t.Errorf("param {%v:%v} not present\n", k, v)
				}
			}
		}
	}
}
func TestCreateNotificationWithParams2(t *testing.T) {
	n, err := MakeNotification(testNotificationMethod, testNotificationParams2)
	if err != nil {
		t.Fatal(err)
	}
	if n == nil {
		t.Error("notification should not be nil")
	}
	if n.Method() != testNotificationMethod {
		t.Error("method not set correctly")
	}
	if n.Params() == nil {
		t.Error("params should not be nil")
	} else {
		if reflect.ValueOf(n.Params()).Kind() != reflect.Slice {
			t.Error("params should be slice")
		} else {
			sliceParams := n.Params().([]interface{})
			for i := range testNotificationParams2 {
				// check for each of the params in the notification
				if !sliceContains(sliceParams, testNotificationParams2[i]) {
					t.Errorf("param %v not present\n", testNotificationParams2[i])
				}
			}
		}
	}
}

func TestCreateNotificationWithInvalidParamsType(t *testing.T) {
	// Params must be of type array/slice/map.
	n, err := MakeNotification(testNotificationMethod, "invalid params type")
	if err == nil {
		t.Error("should have returned an error")
	} else if err != InvalidNotificationInvalidParamsType {
		t.Error("wrong error returned")
	}
	if n != nil {
		t.Error("notification should be nil")
	}
}

func TestCreateNotificationWithInvalidParamsType2(t *testing.T) {
	// If params is a map, the keys must be strings.
	n, err := MakeNotification(testNotificationMethod, map[int]interface{}{1: "test", 2: true})
	if err == nil {
		t.Error("should have returned an error")
	} else if err != InvalidNotificationInvalidParamsType {
		t.Error("wrong error returned")
	}
	if n != nil {
		t.Error("notification should be nil")
	}
}

func TestMarshalThenUnmarshalNotification(t *testing.T) {
	n, err := MakeNotification(testNotificationMethod, testNotificationParams)
	if err != nil {
		t.Fatal(err)
	}
	jsonNotif, err := json.Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	unmarshalN := new(Notification)
	if err = json.Unmarshal(jsonNotif, unmarshalN); err != nil {
		t.Fatal(err)
	}

	if unmarshalN.Method() != n.Method() {
		t.Error("method not correct")
	}
	if unmarshalN.Params() != nil {
		nMapParams := n.Params().(map[string]interface{})
		unmarshalNMapParams := unmarshalN.Params().(map[string]interface{})
		for k, v := range nMapParams {
			if !mapContains(unmarshalNMapParams, k, v) {
				t.Errorf("param {%v:%v} not present\n", k, v)
			}
		}
	} else {
		t.Error("params should not be nil")
	}
}

func TestMarshalNotificationWithParams(t *testing.T) {
	n, err := MakeNotification(testNotificationMethod, testNotificationParams2)
	if err != nil {
		t.Fatal(err)
	}
	jsonNotif, err := json.Marshal(n)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := fmt.Sprintf(`{"%s":"%s","%s":"%s","%s":["%s",%v,%v,%v]}`,
		VersionKey, Version,
		MethodKey, testNotificationMethod,
		ParamsKey, testNotificationParams2[0], testNotificationParams2[1],
		testNotificationParams2[2], testNotificationParams2[3])
	if string(jsonNotif) != expectedJSON {
		t.Errorf("expected %s, got %s\n", expectedJSON, jsonNotif)
	}
}

func TestMarshalNotificationWithoutParams(t *testing.T) {
	n, err := MakeNotification(testNotificationMethod, nil)
	if err != nil {
		t.Fatal(err)
	}
	jsonNotif, err := json.Marshal(n)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := fmt.Sprintf(`{"%s":"%s","%s":"%s"}`,
		VersionKey, Version,
		MethodKey, testNotificationMethod)
	if string(jsonNotif) != expectedJSON {
		t.Errorf("expected %s, got %s\n", expectedJSON, jsonNotif)
	}
}

func TestUnmarshalNotificationWithParams(t *testing.T) {
	jsonNotif := `{"jsonrpc":"2.0", "method":"test", "params":{"key1":"value1","key2":2,"key3":false}}`
	matchingN, err := MakeNotification(testNotificationMethod, testNotificationParams)
	if err != nil {
		t.Fatal("could not make matching notification")
	}
	n := new(Notification)
	if err = json.Unmarshal([]byte(jsonNotif), n); err != nil {
		t.Fatal(err)
	}

	if n.Method() != matchingN.Method() {
		t.Error("method not set correctly")
	}
	if n.Params() != nil {
		nMapParams := n.Params().(map[string]interface{})
		matchingNMapParams := matchingN.Params().(map[string]interface{})
		for k, v := range matchingNMapParams {
			if !mapContains(nMapParams, k, v) {
				t.Errorf("param {%v:%v} not present\n", k, v)
			}
		}
	} else {
		t.Error("params should not be nil")
	}
}

func TestUnmarshalNotificationWithoutParams(t *testing.T) {
	jsonNotif := `{"jsonrpc":"2.0", "method":"test"}`
	matchingN, err := MakeNotification(testNotificationMethod, nil)
	if err != nil {
		t.Fatal("could not make matching notification")
	}
	n := new(Notification)
	if err = json.Unmarshal([]byte(jsonNotif), n); err != nil {
		t.Fatal(err)
	}

	if n.Method() != matchingN.Method() {
		t.Error("method not set correctly")
	}
	if n.Params() != nil {
		t.Error("params should be nil")
	}
}
