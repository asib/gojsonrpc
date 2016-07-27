package gojsonrpc

import "encoding/json"

type ResponseType int

const (
	ResponseTypeError ResponseType = iota
	ResponseTypeResult
)

type responseData struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Err     *Error      `json:"error,omitempty"`
	ID      uint        `json:"id"`
	_type   ResponseType
}

type Response struct {
	responseData
}

func (r *Response) Result() interface{} {
	return r.responseData.Result
}

func (r *Response) Error() *Error {
	return r.responseData.Err
}

func (r *Response) ID() uint {
	return r.responseData.ID
}

func (r *Response) IsResult() bool {
	return r.responseData._type == ResponseTypeResult
}

func (r *Response) IsError() bool {
	return !r.IsResult()
}

func ResultResponseValidAndExpectedKeys() map[string]bool {
	return map[string]bool{"jsonrpc": true, "result": true, "id": true}
}

func ErrorResponseValidAndExpectedKeys() map[string]bool {
	return map[string]bool{"jsonrpc": true, "error": true, "id": true}
}

func makeResponse(result interface{}, err *Error, id uint, _type ResponseType) *Response {
	return &Response{
		responseData{
			Jsonrpc: version,
			Result:  result,
			Err:     err,
			ID:      id,
			_type:   _type,
		},
	}
}

func MakeResponseWithResult(result interface{}, id uint) *Response {
	return makeResponse(result, nil, id, ResponseTypeResult)
}

func MakeResponseWithError(err *Error, id uint) (*Response, error) {
	if err == nil {
		return nil, invalidResponseNilError
	}

	return makeResponse(nil, err, id, ResponseTypeError), nil
}

func (r *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.responseData)
}

func (r *Response) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &(r.responseData))
	if err != nil {
		return err
	}

	if r.Error() != nil && r.Result() == nil {
		r.responseData._type = ResponseTypeError
	} else if r.Result() != nil && r.Error() == nil {
		r.responseData._type = ResponseTypeResult
	} else {
		err = invalidMessage
	}

	return err
}
