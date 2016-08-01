package gojsonrpc

import "encoding/json"

type responseType int

const (
	responseTypeError responseType = iota
	responseTypeResult
)

type responseData struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Err     *Error      `json:"error,omitempty"`
	ID      uint        `json:"id"`
	_type   responseType
}

// Response is a struct that holds information about a response. It is very
// strict about being specifically either a result Response or an error Response.
type Response struct {
	responseData
}

func (r *Response) JSONRPCVersion() string {
	return r.responseData.Jsonrpc
}

// Result returns the response's result.
func (r *Response) Result() interface{} {
	return r.responseData.Result
}

// Error returns the response's error.
func (r *Response) Error() *Error {
	return r.responseData.Err
}

// ID returns the response's ID.
func (r *Response) ID() uint {
	return r.responseData.ID
}

// IsResult is a helper method for determining if a Response is an error or
// result Response. It relies on internal state, so will only work if the Response
// was created using the public interface (i.e. MakeResponseWithXXXXX or ParseIncoming).
func (r *Response) IsResult() bool {
	return r.responseData._type == responseTypeResult
}

// IsError is a helper method for determining if a Response is an error or
// result Response. It relies on internal state, so will only work if the Response
// was created using the public interface (i.e. MakeResponseWithXXXXX or ParseIncoming).
func (r *Response) IsError() bool {
	return !r.IsResult()
}

// ResultResponseValidAndExpectedKeys is a map whose keys are all the possible fields in a result
// response, mapped to whether they are required fields. This mapping is used
// by ParseIncoming to determine the type of the incoming message.
var ResultResponseValidAndExpectedKeys = map[string]bool{"jsonrpc": true, "result": true, "id": true}

// ErrorResponseValidAndExpectedKeys is a map whose keys are all the possible fields in an error
// response, mapped to whether they are required fields. This mapping is used
// by ParseIncoming to determine the type of the incoming message.
var ErrorResponseValidAndExpectedKeys = map[string]bool{"jsonrpc": true, "error": true, "id": true}

func makeResponse(result interface{}, err *Error, id uint, _type responseType) *Response {
	return &Response{
		responseData{
			Jsonrpc: Version,
			Result:  result,
			Err:     err,
			ID:      id,
			_type:   _type,
		},
	}
}

// Use this method to create result Response structs - do not try to use a struct
// literal. You may pass nil for the result argument - it will marshal to JSON's
// null.
func MakeResponseWithResult(result interface{}, id uint) *Response {
	return makeResponse(result, nil, id, responseTypeResult)
}

// Use this method to create error Response structs - do not try to use a struct
// literal. Passing nil for the error argument will cause an error to be returned.
func MakeResponseWithError(err *Error, id uint) (*Response, error) {
	if err == nil {
		return nil, InvalidResponseNilError
	}

	return makeResponse(nil, err, id, responseTypeError), nil
}

// Do not use this method directly. Instead use json.Marshal with a Response
// as the argument.
func (r *Response) MarshalJSON() ([]byte, error) {
	// Result is allowed to be nil, which marshals to JSON's null.
	// Using an unnamed struct here because I don't see that we'll ever use a
	// struct like this anywhere else.
	if r.IsResult() && r.Result() == nil {
		return json.Marshal(struct {
			Jsonrpc string      `json:"jsonrpc"`
			Result  interface{} `json:"result"`
			ID      uint        `json:"id"`
		}{
			Jsonrpc: r.JSONRPCVersion(),
			Result:  nil,
			ID:      r.ID(),
		})
	}

	return json.Marshal(r.responseData)
}

// Do not use this method directly. Instead use ParseIncoming and type assert
// the returned value to a Response.
func (r *Response) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &(r.responseData))
	if err != nil {
		return err
	}

	if r.Error() != nil && r.Result() == nil {
		r.responseData._type = responseTypeError
	} else if r.Result() != nil && r.Error() == nil {
		r.responseData._type = responseTypeResult
	} else if r.Error() == nil && r.Result() == nil {
		// It's possible for the result field to be `null` in JSON. This will unmarshal
		// to nil, leaving both r.Error() and r.Result() nil, so we must check to
		// see if the result field exists in the raw data.
		// If we were at least able to unmarshal into r.responseData, we should be
		// able to unmarshal into a map with string keys.
		var tmp map[string]interface{}
		if err = json.Unmarshal(data, &tmp); err != nil {
			return err
		}

		if _, ok := tmp[ResultKey]; ok {
			if _, ok2 := tmp[ErrorKey]; !ok2 { // only result key present - fine
				r.responseData._type = responseTypeResult
			} else { // both result and error keys present - error
				return InvalidMessage
			}
		} else { // neither result nor error key present - error
			return InvalidMessage
		}
	} else { // both result and error have non-nil values - error
		return InvalidMessage
	}

	return nil
}
