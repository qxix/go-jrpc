package jrpc

type Response struct {
	Id     interface{} `json:"id,omitempty"`
	Result interface{} `json:"result,omitempty"`
	Error  error       `json:"error,omitempty"`
}
