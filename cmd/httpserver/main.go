package main

import (
	"MODULE_NAME/internal/request"
	"MODULE_NAME/internal/response"
	"MODULE_NAME/internal/server"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		response, err := sendHttpRequest(req.RequestLine.RequestTarget)
		if err != nil {
			handler500(w, req)
		}
		chunkedHandler(w, response.Body)
		return
	}
}

func chunkedHandler(w *response.Writer, p io.ReadCloser) {
	buffer := make([]byte, 32)
	headers := response.GetDefaultHeaders(0)
	w.WriteStatusLine(200)
	delete(headers, "Content-Length")
	headers.Set("Transfer-Encoding", "chunked")
	w.WriteHeaders(headers)
	for {
		n, err := p.Read(buffer)
		if n == 0 {
			w.WriteChunkedBodyDone()
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			handler500(w, &request.Request{})
		}
		w.WriteChunkedBody(buffer[:n])
	}
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	body := []byte(`<html>
		<head>
		<title>400 Bad Request</title>
		</head>
		<body>
		<h1>Bad Request</h1>
		<p>Your request honestly kinda sucked.</p>
		</body>
		</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.SetOVR("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	body := []byte(`<html>
		<head>
		<title>500 Internal Server Error</title>
		</head>
		<body>
		<h1>Internal Server Error</h1>
		<p>Okay, you know what? This one is on me.</p>
		</body>
		</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.SetOVR("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	body := []byte(`<html>
		<head>
		<title>200 OK</title>
		</head>
		<body>
		<h1>Success!</h1>
		<p>Your request was an absolute banger.</p>
		</body>
		</html>
		`)
	h := response.GetDefaultHeaders(len(body))
	h.SetOVR("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}
func sendHttpRequest(target string) (*http.Response, error) {
	trimTarget := strings.TrimPrefix(target, "/httpbin")
	urlString := fmt.Sprintf("https://httpbin.org%s", trimTarget)
	response, err := http.Get(urlString)
	if err != nil {
		return nil, err
	}
	return response, nil
}
