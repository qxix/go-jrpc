package jrpc

import (
	"encoding/json"
	"errors"
)

type Request struct {
	Id     interface{}     `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	UserId int             `json:"-"`
}

func (r *Request) Bind(obj interface{}) error {
	err := json.Unmarshal(r.Params, obj)
	if err != nil {
		return InvalidParamsError
	}

	return nil
}

func (r *Request) Result(res interface{}) Response {
	return Response{Id: r.Id, Result: res}
}

func (r *Request) Error(err error) Response {

	var jrpcErr *Error
	if errors.As(err, &jrpcErr) {
		return Response{Id: r.Id, Error: err}
	}

	return Response{Id: r.Id, Error: &Error{
		Code:    -32603,
		Message: err.Error(),
	}}
}
