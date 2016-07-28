package gojsonrpc

import "fmt"

const (
	Version = "2.0"

	VersionKey      = "jsonrpc"
	MethodKey       = "method"
	ParamsKey       = "params"
	IDKey           = "id"
	ResultKey       = "result"
	ErrorKey        = "error"
	ErrorCodeKey    = "code"
	ErrorMessageKey = "message"
	ErrorDataKey    = "data"
)

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// ParseError is a type definition allowing callers of methods that may
// return an error of type ParseError to type assert the returned error.
// ParseError's may be thrown by any function that parses (chiefly ParseIncoming).
type ParseError int

const (
	InvalidVersion ParseError = iota
	InvalidMessage
)

var parseErrors = []string{"invalid version", "invalid message"}

// This method returns the string representation of a ParseError.
func (e ParseError) Error() string {
	return fmt.Sprintf("gojsonrpc: parse error: %s", parseErrors[e])
}

// ObjectError's may be thrown by any method that must create a Notification,
// Request or Response struct.
type ObjectError int

const (
	InvalidNotificationInvalidParamsType ObjectError = iota
	InvalidRequestInvalidParamsType
	InvalidResponseNilError
)

var objectErrors = []string{"invalid notification: invalid params type",
	"invalid request: invalid params type",
	"invalid response: nil result",
	"invalid response: nil error"}

// This method returns the string representation of an ObjectError.
func (e ObjectError) Error() string {
	return fmt.Sprintf("gojsonrpc: object error: %s", objectErrors[e])
}
