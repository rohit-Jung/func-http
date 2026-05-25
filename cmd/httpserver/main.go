package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rohit-Jung/http-protocol/internal/server"
)

const PORT = 42069

func main() {
	server, err := server.Serve(PORT)
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
