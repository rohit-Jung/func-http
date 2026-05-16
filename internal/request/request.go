// Package request
package request

import (
	"errors"
	"io"
	"strings"
)

type State string

const (
	Init State = "init"
	Done State = "done"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HTTPVersion   string
}

type Request struct {
	RequestLine *RequestLine
	Headers     map[string]string
	Body        string
	State       State
}

const CRLF = "\r\n"

var (
	malformedRequest       = "[MALFORMED] error while reading request"
	malformedRequestLine   = "[MALFORMED] error while reading request line"
	malformedHTTPVersion   = "[MALFORMED] couldn't determine the httpversion"
	malformedHTTPMethod    = "[MALFORMED] http method is malformed"
	unsupportedHTTPVersion = "[UNSUPPORTED] http version currently not supported"
)

// func (r *Request) parse(data []byte) (int, error) {}

func parseRequestLine(buf []byte) (*RequestLine, error) {
	body := strings.Split(string(buf), CRLF)
	if len(body) < 2 {
		return nil, errors.New(malformedRequest)
	}

	startLine := body[0]

	// GET / HTTP/1.1
	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, errors.New(malformedRequestLine)
	}

	method := parts[0]
	if strings.ToUpper(method) != method {
		return nil, errors.New(malformedHTTPMethod)
	}

	httpVersion := strings.Split(parts[2], "/")
	if len(httpVersion) != 2 {
		return nil, errors.New(malformedHTTPVersion)
	}

	if httpVersion[1] != "1.1" {
		return nil, errors.New(unsupportedHTTPVersion)
	}

	reqLine := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HTTPVersion:   httpVersion[1],
	}

	return reqLine, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New(malformedRequestLine)
	}

	reqLine, err := parseRequestLine(buf)
	if err != nil {
		return nil, err
	}

	req := &Request{
		RequestLine: reqLine,
		Headers:     map[string]string{},
		Body:        "",
	}

	return req, nil
}
