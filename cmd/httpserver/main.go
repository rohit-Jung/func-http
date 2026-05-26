package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rohit-Jung/http-protocol/internal/headers"
	"github.com/rohit-Jung/http-protocol/internal/request"
	"github.com/rohit-Jung/http-protocol/internal/response"
	"github.com/rohit-Jung/http-protocol/internal/server"
)

const PORT = 42069

func getHtml(statusCode response.StatusCode) string {
	switch statusCode {
	case response.StatusBadRequst:
		return `
		<html>
			<head>
				<title>400 Bad Request</title>
			</head>
			<body>
				<h1>Bad Request</h1>
				<p>Your request honestly kinda sucked.</p>
			</body>
		</html>
		`
	case response.StatusInternalServerError:
		return `
		<html>
			<head>
				<title>500 Internal Server Error</title>
			</head>
			<body>
				<h1>Internal Server Error</h1>
				<p>Okay, you know what? This one is on me.</p>
			</body>
		</html>
		`

	case response.StatusOk:
		return `
			<html>
				<head>
					<title>200 OK</title>
				</head>
				<body>
					<h1>Success!</h1>
					<p>Your request was an absolute banger.</p>
				</body>
			</html>
		`
	default:
		return ""
	}
}

// server will error out if written in wrong order
func writeResponse(w response.Writer, statusCode response.StatusCode, message string) {
	headers := response.GetDefaultHeaders(len(message))
	w.WriteStatusLine(statusCode)
	headers.Replace("Content-Type", "text/html")

	w.WriteHeaders(headers)
	w.WriteBody([]byte(message))
}

func handlePath(w response.Writer, req *request.Request) {
	statusCode := response.StatusOk

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		statusCode = response.StatusBadRequst
		html := getHtml(statusCode)
		writeResponse(w, statusCode, html)
		return
	case "/myproblem":
		statusCode = response.StatusInternalServerError
		html := getHtml(statusCode)
		writeResponse(w, statusCode, html)
		return
	}

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		target := req.RequestLine.RequestTarget[len("/httpbin"):]
		endPoint := fmt.Sprintf("http://httpbin.org%s", target)

		res, err := http.Get(endPoint)
		if err != nil {
			statusCode = response.StatusInternalServerError
			html := getHtml(statusCode)
			writeResponse(w, statusCode, html)
		}

		defer res.Body.Close()
		chunkedEncoderHeaders := response.GetChunkedEncodingHeaders()

		w.WriteStatusLine(response.StatusOk)
		chunkedEncoderHeaders.Set("Trailer", "X-Content-SHA256")
		chunkedEncoderHeaders.Set("Trailer", " X-Content-Length")
		w.WriteHeaders(chunkedEncoderHeaders)

		fullResponseBody := []byte{}
		for {
			b := make([]byte, 32)
			n, err := res.Body.Read(b)
			if n > 0 {
				w.WriteChunkedBody(b[:n])
				fullResponseBody = append(fullResponseBody, b[:n]...)
			}

			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return
			}
		}

		w.WriteChunkedBodyDone()
		bodyHash := sha256.Sum256(fullResponseBody)
		bodyLen := len(fullResponseBody)

		trailers := headers.NewHeaders()
		trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", bodyHash))
		trailers.Set("X-Content-Length", fmt.Sprint(bodyLen))
		w.WriteTrailers(trailers)
	}
}

func main() {
	server, err := server.Serve(PORT, handlePath)
	if err != nil {
		log.Fatalf("Error starting server %v", err)
	}

	defer server.Close()
	log.Println("Server started successfully on portorico: ", PORT)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
