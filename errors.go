package jrpc

import "fmt"

var (
	ParseError          = &Error{Code: -32700, Message: "JSON parse error"}
	InvalidRequestError = &Error{Code: -32600, Message: "invalid request"}
	MethodNotFoundError = &Error{Code: -32601, Message: "method not found"}
	UnauthorizedError   = &Error{Code: -32099, Message: "unauthorized"}
	InvalidParamsError  = &Error{Code: -32602, Message: "invalid params"}
	InternalError       = &Error{Code: -32603, Message: "internal error"}
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}
