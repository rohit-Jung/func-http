// Package request
package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/rohit-Jung/http-protocol/internal/headers"
)

type State string

const (
	StateInit           State = "init"
	StateDone           State = "done"
	StateError          State = "error"
	StateParsingHeaders State = "parsingHeaders"
	StateParsingBody    State = "parsingBody"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HTTPVersion   string
}

type Request struct {
	RequestLine *RequestLine
	Headers     headers.Headers
	Body        []byte
	State       State
}

const CRLF = "\r\n"

var (
	errMalformedRequestLine   = fmt.Errorf("[MALFORMED] error while reading request line")
	errMalformedHTTPVersion   = fmt.Errorf("[MALFORMED] couldn't determine the httpversion")
	errMalformedHTTPMethod    = fmt.Errorf("[MALFORMED] http method is malformed")
	errRequestState           = fmt.Errorf("[STATE ERROR] error in request state")
	errUnsupportedHTTPVersion = fmt.Errorf("[UNSUPPORTED] http version currently not supported")
	errIncompleteRequest      = fmt.Errorf("[MALFORMED] the request was incomplete")
)

func parseRequestLine(buf []byte) (*RequestLine, int, error) {
	indexOfCrlf := bytes.Index(buf, []byte(CRLF))
	if indexOfCrlf == -1 {
		return nil, 0, nil
	}

	startLine := buf[:indexOfCrlf]
	readBytes := indexOfCrlf + len(CRLF)

	// GET / HTTP/1.1
	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, errMalformedRequestLine
	}

	httpMethod := parts[0]
	if strings.ToUpper(string(httpMethod)) != string(httpMethod) {
		return nil, 0, errMalformedHTTPMethod
	}

	httpVersion := bytes.Split(parts[2], []byte("/"))
	if len(httpVersion) != 2 {
		return nil, 0, errMalformedHTTPVersion
	}

	if !bytes.Equal(httpVersion[1], []byte("1.1")) {
		return nil, 0, errUnsupportedHTTPVersion
	}

	reqLine := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HTTPVersion:   string(httpVersion[1]),
	}

	return reqLine, readBytes, nil
}

func newRequest() *Request {
	return &Request{
		State:   StateInit,
		Headers: headers.NewHeaders(),
		Body:    []byte{},
	}
}

const ContentLengthKey = "content-length"

func (r *Request) parse(buf []byte) (int, error) {
	readBytes := 0
dance:
	for {
		currentData := buf[readBytes:]
		switch r.State {
		case StateError:
			return 0, errRequestState

		case StateParsingBody:
			contentLength := r.Headers.GetIntVal(ContentLengthKey, 0)
			if contentLength == 0 {
				r.State = StateDone
				break
			}

			remainingBodyBytes := contentLength - len(r.Body)
			availableCurrentBuf := len(buf) - readBytes

			// answer this why min ?
			bytesToRead := min(remainingBodyBytes, availableCurrentBuf)

			r.Body = append(r.Body, buf[readBytes:bytesToRead+readBytes]...)
			readBytes += bytesToRead
			if bytesToRead == 0 {
				break dance
			}

			if len(r.Body) == contentLength {
				r.State = StateDone
			}

			// this might not happen but still
			if len(r.Body) > contentLength {
				r.State = StateError
			}

		case StateParsingHeaders:
			headersBytesRead, doneParsingHeaders, err := r.Headers.Parse(currentData)
			if err != nil {
				return 0, err
			}

			if headersBytesRead == 0 {
				break dance
			}

			readBytes += headersBytesRead
			if doneParsingHeaders {
				r.State = StateParsingBody
			}

		case StateInit:
			requestLime, rlParsedBytes, err := parseRequestLine(currentData)
			if err != nil {
				r.State = StateError
				return 0, err
			}

			if rlParsedBytes == 0 {
				break dance
			}

			r.RequestLine = requestLime
			readBytes += rlParsedBytes
			r.State = StateParsingHeaders

		case StateDone:
			break dance
		}
	}
	return readBytes, nil
}

func (r *Request) done() bool {
	return r.State == StateDone || r.State == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0

	for !req.done() {
		bytesRead, err := reader.Read(buf[bufLen:])

		if bytesRead > 0 {
			// parsed till read ones
			bufLen += bytesRead
			bytesParsed, err := req.parse(buf[:bufLen])
			if err != nil {
				return nil, err
			}

			// copy back whats not parsed and start from where parsing was stopped
			copy(buf, buf[bytesParsed:bufLen])
			bufLen -= bytesParsed
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				if !req.done() {
					return nil, errIncompleteRequest
				}

				break
			}

			return nil, err
		}
	}

	return req, nil
}
