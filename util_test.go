package gojsonrpc

import "testing"

func TestSliceContains(t *testing.T) {
	if sliceContains([]interface{}{"yes"}, "no") {
		t.Error("should not have been true")
	}
}

func TestMapContainsIncorrectKeyIncorrectValue(t *testing.T) {
	if mapContains(map[string]interface{}{"ey": "no"}, "key", "val") {
		t.Error("should not have been true")
	}
}

func TestMapContainsCorrectKeyIncorrectValue(t *testing.T) {
	if mapContains(map[string]interface{}{"key": "no"}, "key", "val") {
		t.Error("should not have been true")
	}
}
