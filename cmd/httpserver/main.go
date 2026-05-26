package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rohit-Jung/http-protocol/internal/request"
	"github.com/rohit-Jung/http-protocol/internal/response"
	"github.com/rohit-Jung/http-protocol/internal/server"
)

const PORT = 42069

func handlePath(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: response.StatusBadRequst,
			Message:    "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    "Woopsie, my bad\n",
		}
	default:
		w.Write([]byte("All good, frfr\n"))
		return nil
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
