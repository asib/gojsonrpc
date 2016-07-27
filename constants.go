package gojsonrpc

import "fmt"

const (
	version = "2.0"

	versionKey      = "jsonrpc"
	methodKey       = "method"
	paramsKey       = "params"
	idKey           = "id"
	resultKey       = "result"
	errorKey        = "error"
	errorCodeKey    = "code"
	errorMessageKey = "message"
	errorDataKey    = "data"
)

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

type ParseError int

const (
	invalidVersion ParseError = iota
	invalidMessage
)

var parseErrors = []string{"invalid version", "invalid message"}

func (e ParseError) Error() string {
	return fmt.Sprintf("gojsonrpc: parse error: %s", parseErrors[e])
}

type ObjectError int

const (
	invalidNotificationInvalidParamsType ObjectError = iota
	invalidRequestInvalidParamsType
	invalidResponseNilError
)

var objectErrors = []string{"invalid notification: invalid params type",
	"invalid request: invalid params type",
	"invalid response: nil error"}

func (e ObjectError) Error() string {
	return fmt.Sprintf("gojsonrpc: object error: %s", objectErrors[e])
}
