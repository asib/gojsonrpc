//go:generate stringer -type=ParseError
//go:generate stringer -type=ObjectError
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

// Use this type where you only want a JSON-RPC message (notification, request
// or response).
type Message interface {
	JSONRPCVersion() string
}

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

// This method returns the string representation of a ParseError.
func (e ParseError) Error() string {
	return fmt.Sprintf("gojsonrpc: parse error: %s", e.String())
}

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// ObjectError's may be thrown by any method that must create a Notification,
// Request or Response struct.
type ObjectError int

const (
	InvalidNotificationInvalidParamsType ObjectError = iota
	InvalidRequestInvalidParamsType
	InvalidResponseNilError
)

// This method returns the string representation of an ObjectError.
func (e ObjectError) Error() string {
	return fmt.Sprintf("gojsonrpc: object error: %s", e.String())
}
