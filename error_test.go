package gojsonrpc

import (
	"encoding/json"
	"reflect"
	"testing"
)

var (
	testErrorCode    = 1
	testErrorMessage = "test message"
	testErrorData    = "test"
	testErrorData2   = map[string]interface{}{"key1": "value1", "key2": 2, "key3": true}
)

func TestCreateErrorWithoutData(t *testing.T) {
	e := MakeError(testErrorCode, testErrorMessage, nil)
	if e.Code() != testErrorCode {
		t.Error("code not set correctly")
	}
	if e.Message() != testErrorMessage {
		t.Error("message not set correctly")
	}
	if e.Data() != nil {
		t.Error("data should be nil")
	}
}

func TestCreateErrorWithData(t *testing.T) {
	e := MakeError(testErrorCode, testErrorMessage, testErrorData)
	if e.Code() != testErrorCode {
		t.Error("code not set correctly")
	}
	if e.Message() != testErrorMessage {
		t.Error("message not set correctly")
	}
	if e.Data() != testErrorData {
		t.Error("data not set correctly")
	}
}

func TestCreateErrorWithData2(t *testing.T) {
	e := MakeError(testErrorCode, testErrorMessage, testErrorData2)
	if e.Code() != testErrorCode {
		t.Error("code not set correctly")
	}
	if e.Message() != testErrorMessage {
		t.Error("message not set correctly")
	}
	if e.Data() == nil {
		t.Error("data should not be nil")
	} else {
		if reflect.ValueOf(e.Data()).Kind() != reflect.Map {
			t.Error("data should be map")
		} else {
			mapData := e.Data().(map[string]interface{})
			for k, v := range testErrorData2 {
				if !mapContains(mapData, k, v) {
					t.Errorf("data {%v:%v} not present\n", k, v)
				}
			}
		}
	}
}

func TestMarshalThenUnmarshalError(t *testing.T) {
	e := MakeError(testErrorCode, testErrorMessage, testErrorData)
	jsonErr, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	unmarshalE := new(Error)
	if err = json.Unmarshal(jsonErr, unmarshalE); err != nil {
		t.Fatal(err)
	}

	if unmarshalE.Code() != e.Code() {
		t.Error("code not correct")
	}
	if unmarshalE.Message() != e.Message() {
		t.Error("message not correct")
	}
	if unmarshalE.Data() != e.Data() {
		t.Error("data not correct")
	}
}

func TestMarshalWithData(t *testing.T) {
	e := MakeError(testErrorCode, testErrorMessage, testErrorData)
	jsonErr, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := `{"code":1,"message":"test message","data":"test"}`

	if expectedJSON != string(jsonErr) {
		t.Errorf("expected %s got %s\n", expectedJSON, string(jsonErr))
	}
}

func TestMarshalWithoutData(t *testing.T) {
	e := MakeError(testErrorCode, testErrorMessage, nil)
	jsonErr, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := `{"code":1,"message":"test message"}`

	if expectedJSON != string(jsonErr) {
		t.Errorf("expected %s got %s\n", expectedJSON, string(jsonErr))
	}
}

func TestUnmarshalWithData(t *testing.T) {
	jsonErr := `{"code":1,"message":"test message","data":"test"}`
	matchingE := MakeError(testErrorCode, testErrorMessage, testErrorData)

	e := new(Error)
	if err := json.Unmarshal([]byte(jsonErr), e); err != nil {
		t.Fatal(err)
	}

	if e.Code() != matchingE.Code() {
		t.Error("code not correct")
	}
	if e.Message() != matchingE.Message() {
		t.Error("message not correct")
	}
	if e.Data() != matchingE.Data() {
		t.Error("data not correct")
	}
}

func TestUnmarshalWithoutData(t *testing.T) {
	jsonErr := `{"code":1,"message":"test message"}`
	matchingE := MakeError(testErrorCode, testErrorMessage, nil)

	e := new(Error)
	if err := json.Unmarshal([]byte(jsonErr), e); err != nil {
		t.Fatal(err)
	}

	if e.Code() != matchingE.Code() {
		t.Error("code not correct")
	}
	if e.Message() != matchingE.Message() {
		t.Error("message not correct")
	}
	if e.Data() != nil {
		t.Error("data should be nil")
	}
}
