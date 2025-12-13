package main

import (
	"MODULE_NAME/internal/request"
	"MODULE_NAME/internal/response"
	"MODULE_NAME/internal/server"
	"bytes"
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

func handlerFunction(w *response.Writer, req *request.Request) {
	buf := bytes.NewBuffer([]byte{})
	switch req.RequestLine.Method {
	case "GET":
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			buf.Write([]byte(`<html>
							<head>
								<title>400 Bad Request</title>
							</head>
							<body>
								<h1>Bad Request</h1>
								<p>Your request honestly kinda sucked.</p>
							</body>
							</html>`))
			w.WriteStatusLine(response.StatusBadRequest)
			header := response.GetDefaultHeaders(buf.Len())
			header.SetOVR("Content-Type", "text/html")
			w.WriteHeaders(header)
			w.WriteBody(buf.Bytes())

		case "/myproblem":
			buf.Write([]byte(`<html>
								<head>
									<title>500 Internal Server Error</title>
								</head>
								<body>
									<h1>Internal Server Error</h1>
									<p>Okay, you know what? This one is on me.</p>
								</body>
								</html>`))
			w.WriteStatusLine(response.StatusInternalError)
			header := response.GetDefaultHeaders(buf.Len())
			header.SetOVR("Content-Type", "text/html")
			w.WriteHeaders(header)
			w.WriteBody(buf.Bytes())
		case "/":
			buf.Write([]byte(`<html>
								<head>
									<title>200 OK</title>
								</head>
								<body>
									<h1>Success!</h1>
									<p>Your request was an absolute banger.</p>
								</body>
								</html>

								`))
			w.WriteStatusLine(response.StatusOK)
			header := response.GetDefaultHeaders(buf.Len())
			header.SetOVR("Content-Type", "text/html")
			w.WriteHeaders(header)
			w.WriteBody(buf.Bytes())
		}

	}

}
