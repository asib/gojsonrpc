package gojsonrpc

import (
	"encoding/json"
	"reflect"
	"testing"
)

var (
	testResultResponseId      = uint(1)
	testResultResponseResult  = "test"
	testResultResponseResult2 = map[string]interface{}{"key1": "value1", "key2": 2, "key3": true}
	testErrorResponseId       = uint(1)
	testErrorResponseError    = MakeError(testErrorCode, testErrorMessage, testErrorData)
)

func TestCreateResponseWithResult(t *testing.T) {
	r := MakeResponseWithResult(testResultResponseResult, testResultResponseId)
	if r == nil {
		t.Error("should return response")
	} else {
		if r.Result() != testResultResponseResult {
			t.Error("result not set correctly")
		}
		if r.Error() != nil {
			t.Error("error should be nil")
		}
		if r.ID() != testResultResponseId {
			t.Error("id not set correctly")
		}
		if !r.IsResult() {
			t.Error("IsResult should be true")
		}
		if r.IsError() {
			t.Error("IsError should be false")
		}
	}
}

func TestCreateResponseWithResult2(t *testing.T) {
	r := MakeResponseWithResult(testResultResponseResult2, testResultResponseId)
	if r == nil {
		t.Error("should return response")
	} else {
		if r.Result() == nil {
			t.Error("result should not be nil")
		} else {
			if reflect.ValueOf(r.Result()).Kind() != reflect.Map {
				t.Error("result should be map")
			} else {
				mapResult := r.Result().(map[string]interface{})
				for k, v := range testResultResponseResult2 {
					if !mapContains(mapResult, k, v) {
						t.Errorf("result {%v:%v} not present\n", k, v)
					}
				}
			}
		}
		if r.Error() != nil {
			t.Error("error should be nil")
		}
		if r.ID() != testResultResponseId {
			t.Error("id not set correctly")
		}
		if !r.IsResult() {
			t.Error("IsResult should be true")
		}
		if r.IsError() {
			t.Error("IsError should be false")
		}
	}
}

func TestCreateResponseWithError(t *testing.T) {
	r, err := MakeResponseWithError(testErrorResponseError, testErrorResponseId)
	if err != nil {
		t.Error(err)
	}
	if r == nil {
		t.Error("should return response")
	} else {
		if r.Error() != testErrorResponseError {
			t.Error("error not set correctly")
		}
		if r.ID() != testErrorResponseId {
			t.Error("id not set correctly")
		}
	}
}

func TestCreateResponseWithNilError(t *testing.T) {
	r, err := MakeResponseWithError(nil, testErrorResponseId)
	if err == nil {
		t.Error("should return error")
	} else if err != InvalidResponseNilError {
		t.Errorf("wrong error returned: expected invalidResponseNilError, got %s\n", err)
	}
	if r != nil {
		t.Error("response should be nil")
	}
}

func TestMarshalThenUnmarshalResponseWithError(t *testing.T) {
	re := MakeError(testErrorCode, testErrorMessage, nil)
	r, err := MakeResponseWithError(re, testErrorResponseId)
	if err != nil {
		t.Fatal(err)
	}

	jsonResp, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	unmarshalR := new(Response)
	if err = json.Unmarshal(jsonResp, unmarshalR); err != nil {
		t.Fatal(err)
	}

	if unmarshalR.Result() != nil {
		t.Error("result should be nil")
	}
	if unmarshalR.Error().Code() != r.Error().Code() {
		t.Error("error code not correct")
	}
	if unmarshalR.Error().Message() != r.Error().Message() {
		t.Error("error message not correct")
	}
	if unmarshalR.Error().Data() != nil {
		t.Error("error data should be nil")
	}
	if unmarshalR.ID() != r.ID() {
		t.Error("id not correct")
	}
	if unmarshalR.IsResult() {
		t.Error("IsResult should be false")
	}
	if !unmarshalR.IsError() {
		t.Error("IsError should be true")
	}
}

func TestMarshalThenUnmarshalResponseWithResult(t *testing.T) {
	r := MakeResponseWithResult(testResultResponseResult, testResultResponseId)

	jsonResp, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	unmarshalR := new(Response)
	if err = json.Unmarshal(jsonResp, unmarshalR); err != nil {
		t.Fatal(err)
	}

	if unmarshalR.Result() != r.Result() {
		t.Error("result not correct")
	}
	if unmarshalR.Error() != nil {
		t.Error("error should be nil")
	}
	if unmarshalR.ID() != r.ID() {
		t.Error("id not correct")
	}
	if !unmarshalR.IsResult() {
		t.Error("IsResult should be true")
	}
	if unmarshalR.IsError() {
		t.Error("IsError should be false")
	}
}

func TestMarshalResponseWithError(t *testing.T) {
	e := MakeError(testErrorCode, testErrorMessage, testErrorData)
	r, err := MakeResponseWithError(e, testErrorResponseId)
	if err != nil {
		t.Fatal(err)
	}

	jsonResp, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := `{"jsonrpc":"2.0","error":{"code":1,"message":"test message","data":"test"},"id":1}`

	if string(jsonResp) != expectedJSON {
		t.Errorf("got %s expected %s\n", string(jsonResp), expectedJSON)
	}
}

func TestMarshalResponseWithResult(t *testing.T) {
	r := MakeResponseWithResult(testResultResponseResult, testResultResponseId)
	jsonResp, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := `{"jsonrpc":"2.0","result":"test","id":1}`

	if string(jsonResp) != expectedJSON {
		t.Errorf("got %s expected %s\n", string(jsonResp), expectedJSON)
	}
}

func TestMarshalResultResponseWithNilResult(t *testing.T) {
	r := MakeResponseWithResult(nil, testResultResponseId)
	jsonResp, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := `{"jsonrpc":"2.0","result":null,"id":1}`

	if string(jsonResp) != expectedJSON {
		t.Errorf("got %s expected %s\n", string(jsonResp), expectedJSON)
	}
}

func TestUnmarshalResponseWithError(t *testing.T) {
	jsonResp := `{"jsonrpc":"2.0", "error":{"code":1, "message":"test message", "data":"test"}, "id":1}`
	matchingE := MakeError(testErrorCode, testErrorMessage, testErrorData)
	matchingR, err := MakeResponseWithError(matchingE, testErrorResponseId)
	if err != nil {
		t.Fatal(err)
	}

	r := new(Response)
	if err = json.Unmarshal([]byte(jsonResp), r); err != nil {
		t.Fatal(err)
	}

	if r.Error().Code() != matchingR.Error().Code() {
		t.Error("code not correct")
	}
	if r.Error().Message() != matchingR.Error().Message() {
		t.Error("message not correct")
	}
	if r.Error().Data() != matchingR.Error().Data() {
		t.Error("data not correct")
	}
	if r.ID() != matchingR.ID() {
		t.Error("id not correct")
	}
}

func TestUnmarshalResponseWithResult(t *testing.T) {
	jsonResp := `{"jsonrpc":"2.0", "result":"test", "id":1}`
	matchingR := MakeResponseWithResult(testResultResponseResult, testErrorResponseId)

	r := new(Response)
	if err := json.Unmarshal([]byte(jsonResp), r); err != nil {
		t.Fatal(err)
	}

	if r.Result() != matchingR.Result() {
		t.Error("result not correct")
	}
	if r.ID() != matchingR.ID() {
		t.Error("id not correct")
	}
}

func TestUnmarshalResponseWithNonNullErrorAndNonNullResult(t *testing.T) {
	jsonResp := `{"jsonrpc":"2.0", "error":{"code":1, "message":"test message", "data":"test"}, "result":"test", "id":1}`
	r := new(Response)
	if err := json.Unmarshal([]byte(jsonResp), r); err == nil {
		t.Error("should have returned an error")
	} else if err != InvalidMessage {
		t.Error("wrong error returned:", err)
	}
}

func TestUnmarshalResponseWithNullErrorAndNullResult(t *testing.T) {
	jsonResp := `{"jsonrpc":"2.0", "error":null, "result":null, "id":1}`
	r := new(Response)
	if err := json.Unmarshal([]byte(jsonResp), r); err == nil {
		t.Error("should have returned an error")
	} else if err != InvalidMessage {
		t.Error("wrong error returned:", err)
	}
}

func TestUnmarshalResponseWithNullResult(t *testing.T) {
	jsonResp := `{"jsonrpc":"2.0", "result":null, "id":1}`
	r := new(Response)
	if err := json.Unmarshal([]byte(jsonResp), r); err != nil {
		t.Fatal(err)
	}

	if !r.IsResult() || r.IsError() {
		t.Error("IsResult should be true and IsError should be false")
	}
}

func TestUnmarshalResponseWithoutResultAndWithoutError(t *testing.T) {
	jsonResp := `{"jsonrpc":"2.0", "id":1}`
	r := new(Response)
	if err := json.Unmarshal([]byte(jsonResp), r); err == nil {
		t.Error("should have returned an error")
	} else if err != InvalidMessage {
		t.Error("wrong error returned:", err)
	}
}
