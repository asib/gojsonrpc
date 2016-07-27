package gojsonrpc

import "testing"

func TestIsNotification(t *testing.T) {
	areNotifications := [][]string{
		[]string{"jsonrpc", "method"},
		[]string{"jsonrpc", "method", "params"},
	}
	areNotNotifications := [][]string{
		[]string{"jsonrpc", "method", "id"},
		[]string{"jsonrpc", "method", "params", "id"},
		[]string{"jsonrpc", "method", "unexpected"},
		[]string{"jsonrpc", "method", "params", "unexpected"},
	}

	for _, slice := range areNotifications {
		if !isNotification(slice) {
			t.Errorf("isNotification should be true: %v\n", slice)
		}
	}

	for _, slice := range areNotNotifications {
		if isNotification(slice) {
			t.Errorf("isNotification should be false: %v\n", slice)
		}
	}
}

func TestIsRequest(t *testing.T) {
	areNotRequests := [][]string{
		[]string{"jsonrpc", "method"},
		[]string{"jsonrpc", "method", "params"},
		[]string{"jsonrpc", "method", "id", "unexpected"},
		[]string{"jsonrpc", "method", "id", "params", "unexpected"},
	}
	areRequests := [][]string{
		[]string{"jsonrpc", "method", "id"},
		[]string{"jsonrpc", "method", "params", "id"},
	}

	for _, slice := range areRequests {
		if !isRequest(slice) {
			t.Errorf("isRequest should be true: %v\n", slice)
		}
	}

	for _, slice := range areNotRequests {
		if isRequest(slice) {
			t.Errorf("isRequest should be false: %v\n", slice)
		}
	}
}

func TestIsErrorResponse(t *testing.T) {
	areErrorResponses := [][]string{
		[]string{"jsonrpc", "error", "id"},
	}
	areNotErrorResponses := [][]string{
		[]string{"jsonrpc", "result", "id"},
		[]string{"jsonrpc", "error", "id", "unexpected"},
		[]string{"jsonrpc", "error"},
		[]string{"jsonrpc", "id"},
	}

	for _, slice := range areErrorResponses {
		if !isErrorResponse(slice) {
			t.Errorf("isErrorResponse should be true: %v\n", slice)
		}
	}

	for _, slice := range areNotErrorResponses {
		if isErrorResponse(slice) {
			t.Errorf("isErrorResponse should be false: %v\n", slice)
		}
	}
}

func TestIsResultResponse(t *testing.T) {
	areResultResponses := [][]string{
		[]string{"jsonrpc", "result", "id"},
	}
	areNotResultResponses := [][]string{
		[]string{"jsonrpc", "error", "id"},
		[]string{"jsonrpc", "result", "id", "unexpected"},
		[]string{"jsonrpc", "result"},
		[]string{"jsonrpc", "id"},
	}

	for _, slice := range areResultResponses {
		if !isResultResponse(slice) {
			t.Errorf("isResultResponse should be true: %v\n", slice)
		}
	}

	for _, slice := range areNotResultResponses {
		if isResultResponse(slice) {
			t.Errorf("isResultResponse should be false: %v\n", slice)
		}
	}
}

func TestParseIncomingWithMalformedMessage(t *testing.T) {
	rawMsg := `{"key":"value"}`
	if _, err := ParseIncoming(rawMsg); err != invalidMessage {
		t.Error("should have returned invalid message error")
	}
}

func TestParseIncomingWithInvalidVersion(t *testing.T) {
	rawMsg := `{"jsonrpc":"1.0", "method":"test"}`
	if _, err := ParseIncoming(rawMsg); err != invalidVersion {
		t.Error("should have returned invalid version error")
	}
}

func TestParseIncomingWithNotificationWithoutParams(t *testing.T) {
	rawMsg := `{"jsonrpc":"2.0", "method":"test"}`
	msg, err := ParseIncoming(rawMsg)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := msg.(*Notification); !ok {
		t.Error("should have returned a notification")
	}
}

func TestParseIncomingWithNotificationWithParams(t *testing.T) {
	rawMsg := `{"jsonrpc":"2.0", "method":"test", "params":["test1", "test2"]}`
	msg, err := ParseIncoming(rawMsg)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := msg.(*Notification); !ok {
		t.Error("should have returned a notification")
	}
}

func TestParseIncomingWithNotificationWithInvalidParams(t *testing.T) {
	rawMsg := `{"jsonrpc":"2.0", "method":"test", "params":"test1"}`
	if _, err := ParseIncoming(rawMsg); err != invalidMessage {
		t.Error("should have returned invalid message error")
	}
}

func TestParseIncomingWithRequestWithoutParams(t *testing.T) {
	rawMsg := `{"jsonrpc":"2.0", "method":"test", "id":1}`
	msg, err := ParseIncoming(rawMsg)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := msg.(*Request); !ok {
		t.Error("should have returned request")
	}
}

func TestParseIncomingWithRequestWithParams(t *testing.T) {
	rawMsg := `{"jsonrpc":"2.0", "method":"test", "params":["bla"], "id":1}`
	msg, err := ParseIncoming(rawMsg)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := msg.(*Request); !ok {
		t.Error("should have returned request")
	}
}

func TestParseIncomingWithRequestWithInvalidParams(t *testing.T) {
	rawMsg := `{"jsonrpc":"2.0", "method":"test", "params":"test1", "id":1}`
	if _, err := ParseIncoming(rawMsg); err != invalidMessage {
		t.Error("should have returned invalid message error")
	}
}

func TestParseIncomingWithErrorResponse(t *testing.T) {
	rawMsg := `{"jsonrpc":"2.0", "error":{"code":1, "message":"test"}, "id":1}`
	msg, err := ParseIncoming(rawMsg)
	if err != nil {
		t.Fatal(err)
	}

	if r, ok := msg.(*Response); !ok {
		t.Error("should have returned response")
	} else if !r.IsError() || r.IsResult() {
		t.Error("IsError should be true and IsResult should be false")
	}
}

func TestParseIncomingWithResultResponse(t *testing.T) {
	rawMsg := `{"jsonrpc":"2.0", "result":"test", "id":1}`
	msg, err := ParseIncoming(rawMsg)
	if err != nil {
		t.Fatal(err)
	}

	if r, ok := msg.(*Response); !ok {
		t.Error("should have returned response")
	} else if r.IsError() || !r.IsResult() {
		t.Error("IsError should be false and IsResult should be true")
	}
}

func TestParseIncomingWithMalformedNotification(t *testing.T) {
	rawMsg := `{"jsonrpc":"2.0", "method":"test", "unexpected":"bla"}`
	if _, err := ParseIncoming(rawMsg); err != invalidMessage {
		t.Error("should have returned invalid message error")
	}
}
