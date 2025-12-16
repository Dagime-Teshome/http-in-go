package main

import (
	"MODULE_NAME/internal/headers"
	"MODULE_NAME/internal/request"
	"MODULE_NAME/internal/response"
	"MODULE_NAME/internal/server"
	"crypto/sha256"
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
		proxyHandler(w, req)
		return
		// response, err := sendHttpRequest(req.RequestLine.RequestTarget)
		// if err != nil {
		// 	handler500(w, req)
		// }
		// chunkedHandler(w, response.Body)
		// return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {
		videoHandler(w, req)
		return
	}
}

func videoHandler(w *response.Writer, req *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	vidBytes, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		log.Fatal(err)
		handler500(w, req)
	}
	headers := response.GetDefaultHeaders(len(vidBytes))
	headers.SetOVR("Content-Type", "video/mp4")
	w.WriteHeaders(headers)
	w.WriteBody(vidBytes)

}

func chunkedHandler(w *response.Writer, p io.ReadCloser) {
	defer p.Close()
	buffer := make([]byte, 32)
	headers := response.GetDefaultHeaders(0)
	w.WriteStatusLine(200)
	delete(headers, "Content-Length")
	headers.Set("Transfer-Encoding", "chunked")
	w.WriteHeaders(headers)
	for {
		n, err := p.Read(buffer)
		if err != nil {
			if err == io.EOF {
				w.WriteChunkedBodyDone()
				break
			}
			handler500(w, &request.Request{})
			return
		}
		if n > 0 {
			w.WriteChunkedBody(buffer[:n])
		}
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
	w.WriteStatusLine(response.StatusInternalError)
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

func proxyHandler(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + target
	fmt.Println("Proxying to", url)
	resp, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(0)
	h.SetOVR("Transfer-Encoding", "chunked")
	h.SetOVR("Trailer", " X-Content-Sha256, X-Content-Length")
	h.Delete("Content-Length")
	w.WriteHeaders(h)

	const maxChunkSize = 1024
	buffer := make([]byte, maxChunkSize)
	var hashBuffer []byte
	for {
		n, err := resp.Body.Read(buffer)
		fmt.Println("Read", n, "bytes")
		if n > 0 {
			hashBuffer = append(hashBuffer, buffer[:n]...)
			_, err = w.WriteChunkedBody(buffer[:n])
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading response body:", err)
			break
		}
	}
	_, err = w.WriteChunkedBodyDone()
	trailer := addTrailer(hashBuffer)
	w.WriteTrailers(trailer)
	if err != nil {
		fmt.Println("Error writing chunked body done:", err)
	}
}

func addTrailer(buf []byte) headers.Headers {
	sum := sha256.Sum256(buf)
	headers := headers.NewHeaders()
	headers.Set("X-Content-Sha256", fmt.Sprintf("%x", sum))
	headers.Set("X-Content-Length", fmt.Sprintf("%d", len(buf)))
	return headers
}
