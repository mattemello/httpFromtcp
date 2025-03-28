package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattemello/httpFromtcp/internal/request"
	"github.com/mattemello/httpFromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w io.Writer, req *request.Request) *server.HandlerError {

	switch req.RequestLine.RequestTarget {

	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: 400,
			Message:    "Your problem is not my problem\n",
		}

	case "/myproblem":
		return &server.HandlerError{
			StatusCode: 500,
			Message:    "Woopsie, my bad\n",
		}
	}

	w.Write([]byte("All good, frfr\n"))

	return nil

}
