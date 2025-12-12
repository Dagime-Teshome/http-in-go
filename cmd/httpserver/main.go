package main

import (
	"MODULE_NAME/internal/request"
	"MODULE_NAME/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handlerFunction)
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

func handlerFunction(w io.Writer, req *request.Request) *server.HandlerError {
	var errorHandler = &server.HandlerError{}
	switch req.RequestLine.Method {
	case "GET":
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			errorHandler = &server.HandlerError{
				StatusCode: 400,
				Message:    "Your problem is not my problem",
			}
		case "/myproblem":
			errorHandler = &server.HandlerError{
				StatusCode: 500,
				Message:    "Woopsie, my bad",
			}
		case "/use-nvim":
			w.Write([]byte("All good, frfr"))
			errorHandler = nil
		}

	}

	return errorHandler
}
