package jrpc

import (
	"encoding/json"
	"errors"
	"log"
)

type Payload struct {
	Body      []byte
	AuthToken string
}

func parseRequest(body []byte) (*Request, error) {
	var t Request
	err := json.Unmarshal(body, &t)

	var jsonSyntaxErr *json.SyntaxError
	if errors.As(err, &jsonSyntaxErr) {
		return nil, ParseError
	}

	if t.Method == "" {
		return nil, InvalidRequestError
	}

	return &t, err
}

type HandlerFunc func(r *Request) Response
type AuthFunc func(r *Request, token string) bool

type method struct {
	fn        HandlerFunc
	protected bool
	verbose   bool
}

type JRPC struct {
	methods      map[string]method
	authCallback AuthFunc
}

func NewJRPC() *JRPC {
	return &JRPC{
		methods: make(map[string]method),
		authCallback: func(r *Request, token string) bool {
			return true
		},
	}
}

func (j *JRPC) Method(mName string, h HandlerFunc, opts ...func(*method)) {

	m := method{
		fn: h,
	}

	for _, opt := range opts {
		opt(&m)
	}

	j.methods[mName] = m
}

func Protected() func(*method) {
	return func(m *method) {
		m.protected = true
	}
}

func Verbose() func(*method) {
	return func(m *method) {
		m.verbose = true
	}
}

func (j *JRPC) RegisterAuthCallback(cb AuthFunc) {
	j.authCallback = cb
	print(j.authCallback)
}

func (j *JRPC) Handle(p Payload) Response {
	request, err := parseRequest(p.Body)
	if err != nil {
		return Response{
			Id:     nil,
			Result: nil,
			Error:  err,
		}
	}

	m, ok := j.methods[request.Method]
	if !ok {
		return Response{
			Id:     request.Id,
			Result: nil,
			Error:  MethodNotFoundError,
		}
	}

	if m.verbose {
		log.Default().Printf("jrpc call: method=%s params=%s\n", request.Method, request.Params)
	}

	if !j.authCallback(request, p.AuthToken) && m.protected {
		return Response{
			Id:     request.Id,
			Result: nil,
			Error:  UnauthorizedError,
		}
	}

	if request.Id != nil {
		resp := m.fn(request)
		resp.Id = request.Id
		return resp
	} else {
		go m.fn(request)
		return Response{}
	}

}
