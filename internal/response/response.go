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
	writerStateTrailer    writerState = "writerStateTrailer"
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
	header.Set("Content-Length", fmt.Sprint(contentLen))
	header.Set("Connection", "close")
	header.Set("Content-Type", "text/plain")
	return header
}

func GetChunkedEncodingHeaders() headers.Headers {
	headers := GetDefaultHeaders(0)
	headers.Delete("Content-Length")
	headers.Set("Transfer-Encoding", "chunked")
	return headers
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
	return n, err
}

func (w *Writer) WriteChunkedBody(b []byte) (int, error) {
	chunkHeader := fmt.Appendf(nil, "%x\r\n", len(b))

	body := make([]byte, 0, len(chunkHeader)+len(b)+2)
	body = append(body, chunkHeader...)
	body = append(body, b...)
	body = append(body, '\r', '\n')

	return w.WriteBody(body)
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	bytesWritten, err := w.WriteBody([]byte("0\r\n"))
	w.writerState = writerStateDone
	return bytesWritten, err
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	w.writerState = writerStateHeaders
	err := w.WriteHeaders(h)
	w.writer.Write([]byte("\r\n"))
	w.writerState = writerStateDone
	return err
}
