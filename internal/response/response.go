package response

import (
	"fmt"
	"io"

	"github.com/rohit-Jung/http-protocol/internal/headers"
)

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequst           StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func getReasonPhrase(statusCode StatusCode) string {
	switch statusCode {
	case StatusOk:
		return "OK"
	case StatusBadRequst:
		return "Bad Request"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return " "
	}
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, getReasonPhrase(statusCode))
	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.NewHeaders()
	header["Content-Length"] = fmt.Sprint(contentLen)
	header["Connection"] = "close"
	header["Content-Type"] = "text/plain"
	return header
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		headerStr := fmt.Sprintf("%s: %s\r\n", key, val)
		_, err := w.Write([]byte(headerStr))
		if err != nil {
			return err
		}
	}

	w.Write([]byte("\r\n"))
	return nil
}
