package response

import (
	"fmt"
	"io"
	"log"

	"github.com/rohit-Jung/http-protocol/internal/headers"
)

type writerState string

const (
	writerStateDone       writerState = "writerStateDone"
	writerStateBody       writerState = "writerStateBody"
	writerStateHeaders    writerState = "writerStateHeaders"
	writerStateStatusLine writerState = "writerStateStatusLine"
)

type (
	StatusCode int
	Writer     struct {
		writer      io.Writer
		writerState writerState
	}
)

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer, writerState: writerStateStatusLine}
}

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

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStateStatusLine {
		log.Fatal("[RESPONSE ORDER ERROR]: Status Line Must be written First")
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, getReasonPhrase(statusCode))
	_, err := w.writer.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	w.writerState = writerStateHeaders
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.NewHeaders()
	header["Content-Length"] = fmt.Sprint(contentLen)
	header["Connection"] = "close"
	header["Content-Type"] = "text/plain"
	return header
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writerStateHeaders {
		log.Fatal("[RESPONSE ORDER ERROR]: Headers must be written after status line")
	}

	for key, val := range headers {
		headerStr := fmt.Sprintf("%s: %s\r\n", key, val)
		_, err := w.writer.Write([]byte(headerStr))
		if err != nil {
			return err
		}
	}

	w.writer.Write([]byte("\r\n"))
	w.writerState = writerStateBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		log.Fatal("[RESPONSE ORDER ERROR]: Body must be written after headers")
	}

	n, err := w.writer.Write(p)
	w.writerState = writerStateDone
	return n, err
}
