package gojsonrpc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var (
	testRequestMethod  = "test"
	testRequestId      = uint(1)
	testRequestParams  = map[string]interface{}{"key1": "value1", "key2": float64(2), "key3": false}
	testRequestParams2 = []interface{}{"val1", float64(2), true, 34.21}
)

func TestCreateRequestWithoutParams(t *testing.T) {
	r, err := MakeRequest(testRequestMethod, nil, testRequestId)
	if err != nil {
		t.Error(err)
	}
	if r == nil {
		t.Error("request should not be nil")
	}
	if r.Method() != testRequestMethod {
		t.Error("method not set correctly")
	}
	if r.Params() != nil {
		t.Error("params should be nil")
	}
}

func TestCreateRequestWithParams(t *testing.T) {
	r, err := MakeRequest(testRequestMethod, testRequestParams, testRequestId)
	if err != nil {
		t.Error(err)
	}
	if r == nil {
		t.Error("request should not be nil")
	}
	if r.Method() != testRequestMethod {
		t.Error("method not set correctly")
	}
	if r.Params() == nil {
		t.Error("params should not be nil")
	} else {
		if reflect.ValueOf(r.Params()).Kind() != reflect.Map {
			t.Error("params should be map")
		} else {
			mapParams := r.Params().(map[string]interface{})
			for k, v := range testRequestParams {
				if !mapContains(mapParams, k, v) {
					t.Errorf("param {%v:%v} not present\n", k, v)
				}
			}
		}
	}
}
func TestCreateRequestWithParams2(t *testing.T) {
	r, err := MakeRequest(testRequestMethod, testRequestParams2, testRequestId)
	if err != nil {
		t.Error(err)
	}
	if r == nil {
		t.Error("request should not be nil")
	}
	if r.Method() != testRequestMethod {
		t.Error("method not set correctly")
	}
	if r.Params() == nil {
		t.Error("params should not be nil")
	} else {
		if reflect.ValueOf(r.Params()).Kind() != reflect.Slice {
			t.Error("params should be slice")
		} else {
			sliceParams := r.Params().([]interface{})
			for i := range testRequestParams2 {
				// check for each of the params in the request
				if !sliceContains(sliceParams, testRequestParams2[i]) {
					t.Errorf("param %v not present\n", testRequestParams2[i])
				}
			}
		}
	}
}

func TestCreateRequestWithInvalidParamsType(t *testing.T) {
	r, err := MakeRequest(testRequestMethod, "invalid params type", testRequestId)
	if err == nil {
		t.Error("should have returned an error")
	} else if err != invalidRequestInvalidParamsType {
		t.Error("wrong error returned")
	}
	if r != nil {
		t.Error("request should be nil")
	}
}

func TestCreateRequestWithInvalidParamsType2(t *testing.T) {
	r, err := MakeRequest(testRequestMethod, map[int]interface{}{1: "test", 2: true}, testRequestId)
	if err == nil {
		t.Error("should have returned an error")
	} else if err != invalidRequestInvalidParamsType {
		t.Error("wrong error returned")
	}
	if r != nil {
		t.Error("request should be nil")
	}
}

func TestMarshalThenUnmarshalRequest(t *testing.T) {
	r, err := MakeRequest(testRequestMethod, testRequestParams, testRequestId)
	if err != nil {
		t.Fatal(err)
	}
	jsonReq, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	unmarshalR := new(Request)
	if err = json.Unmarshal(jsonReq, unmarshalR); err != nil {
		t.Fatal(err)
	}

	if unmarshalR.Method() != r.Method() {
		t.Error("method not correct")
	}
	if unmarshalR.Params() != nil {
		rMapParams := r.Params().(map[string]interface{})
		unmarshalRMapParams := unmarshalR.Params().(map[string]interface{})
		for k, v := range rMapParams {
			if !mapContains(unmarshalRMapParams, k, v) {
				t.Errorf("param {%v:%v} not present\n", k, v)
			}
		}
	} else {
		t.Error("params should not be nil")
	}
	if unmarshalR.ID() != r.ID() {
		t.Error("id not correct")
	}
}

func TestMarshalRequestWithParams(t *testing.T) {
	r, err := MakeRequest(testRequestMethod, testRequestParams2, testRequestId)
	if err != nil {
		t.Fatal(err)
	}
	jsonReq, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := fmt.Sprintf(`{"%s":"%s","%s":"%s","%s":["%s",%v,%v,%v],"%s":%d}`,
		versionKey, version,
		methodKey, testRequestMethod,
		paramsKey, testRequestParams2[0], testRequestParams2[1], testRequestParams2[2], testRequestParams2[3],
		idKey, testRequestId)
	if string(jsonReq) != expectedJSON {
		t.Errorf("expected %s, got %s\n", expectedJSON, jsonReq)
	}
}

func TestMarshalRequestWithoutParams(t *testing.T) {
	r, err := MakeRequest(testRequestMethod, nil, testRequestId)
	if err != nil {
		t.Fatal(err)
	}
	jsonReq, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	expectedJSON := fmt.Sprintf(`{"%s":"%s","%s":"%s","%s":%d}`,
		versionKey, version,
		methodKey, testRequestMethod,
		idKey, testRequestId)
	if string(jsonReq) != expectedJSON {
		t.Errorf("expected %s, got %s\n", expectedJSON, jsonReq)
	}
}

func TestUnmarshalRequestWithParams(t *testing.T) {
	jsonReq := `{"jsonrpc":"2.0", "method":"test", "params":{"key1":"value1","key2":2,"key3":false}, "id":1}`
	matchingR, err := MakeRequest(testRequestMethod, testRequestParams, testRequestId)
	if err != nil {
		t.Fatal("could not make matching request")
	}
	r := new(Request)
	if err = json.Unmarshal([]byte(jsonReq), r); err != nil {
		t.Fatal(err)
	}

	if r.Method() != matchingR.Method() {
		t.Error("method not set correctly")
	}
	if r.Params() != nil {
		rMapParams := r.Params().(map[string]interface{})
		matchingRMapParams := matchingR.Params().(map[string]interface{})
		for k, v := range matchingRMapParams {
			if !mapContains(rMapParams, k, v) {
				t.Errorf("param {%v:%v} not present\n", k, v)
			}
		}
	} else {
		t.Error("params should not be nil")
	}
	if r.ID() != matchingR.ID() {
		t.Error("id not correct")
	}
}

func TestUnmarshalRequestWithoutParams(t *testing.T) {
	jsonReq := `{"jsonrpc":"2.0", "method":"test", "id":1}`
	matchingR, err := MakeRequest(testRequestMethod, nil, testRequestId)
	if err != nil {
		t.Fatal("could not make matching request")
	}
	r := new(Request)
	if err = json.Unmarshal([]byte(jsonReq), r); err != nil {
		t.Fatal(err)
	}

	if r.Method() != matchingR.Method() {
		t.Error("method not set correctly")
	}
	if r.Params() != nil {
		t.Error("params should be nil")
	}
	if r.ID() != matchingR.ID() {
		t.Error("id not correct")
	}
}
