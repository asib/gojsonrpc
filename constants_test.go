package gojsonrpc

import "testing"

func TestParseErrorInvalidVersionString(t *testing.T) {
	expected := "gojsonrpc: parse error: InvalidVersion"
	if InvalidVersion.Error() != expected {
		t.Errorf("expected %q got %q", expected, InvalidVersion.Error())
	}
}

func TestParseErrorInvalidMessageString(t *testing.T) {
	expected := "gojsonrpc: parse error: InvalidMessage"
	if InvalidMessage.Error() != expected {
		t.Errorf("expected %q got %q", expected, InvalidMessage.Error())
	}
}

func TestParseErrorInvalid(t *testing.T) {
	expected := "gojsonrpc: parse error: ParseError(999)"
	if ParseError(999).Error() != expected {
		t.Errorf("expected %q got %q", expected, ParseError(999).Error())
	}
}

func TestObjectErrorInvalidNotificationInvalidParamsTypeString(t *testing.T) {
	expected := "gojsonrpc: object error: InvalidNotificationInvalidParamsType"
	if InvalidNotificationInvalidParamsType.Error() != expected {
		t.Errorf("expected %q got %q", expected, InvalidNotificationInvalidParamsType.Error())
	}
}

func TestObjectErrorInvalidRequestInvalidParamsType(t *testing.T) {
	expected := "gojsonrpc: object error: InvalidRequestInvalidParamsType"
	if InvalidRequestInvalidParamsType.Error() != expected {
		t.Errorf("expected %q got %q", expected, InvalidRequestInvalidParamsType.Error())
	}
}

func TestObjectErrorInvalidResponseNilError(t *testing.T) {
	expected := "gojsonrpc: object error: InvalidResponseNilError"
	if InvalidResponseNilError.Error() != expected {
		t.Errorf("expected %q got %q", expected, InvalidResponseNilError.Error())
	}
}

func TestObjectErrorInvalid(t *testing.T) {
	expected := "gojsonrpc: object error: ObjectError(999)"
	if ObjectError(999).Error() != expected {
		t.Errorf("expected %q got %q", expected, ObjectError(999).Error())
	}
}
