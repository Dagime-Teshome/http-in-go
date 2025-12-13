package server

import (
	"MODULE_NAME/internal/request"
	"MODULE_NAME/internal/response"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type ServerStatus int

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

// type HandlerError struct {
// 	StatusCode response.StatusCode
// 	Message    string
// }

// func (he HandlerError) write(w io.Writer) {

// 	response.WriteStatusLine(w, he.StatusCode)
// 	headers := response.GetDefaultHeaders(len(he.Message))
// 	response.WriteHeaders(w, headers)
// 	w.Write([]byte(he.Message))
// }

type Handler func(w *response.Writer, req *request.Request)

const (
	Listening ServerStatus = iota
	Closed
)

func (s *Server) Close() {
	s.closed.Store(true)
	if s.listener != nil {
		s.listener.Close()
	}
	return
}

func (s *Server) listen() {

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}

}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	request, err := request.RequestFromReader(conn)
	if err != nil {
		responseWrite := &response.Writer{
			ResWriter: conn,
		}
		responseWrite.WriteStatusLine(400)
		header := response.GetDefaultHeaders(0)
		responseWrite.WriteHeaders(header)
		return
	}
	responseStr := &response.Writer{
		ResWriter: conn,
	}

	s.handler(responseStr, request)

}
func Serve(port int, handler Handler) (*Server, error) {
	portString := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", portString)

	if err != nil {
		return nil, fmt.Errorf("error creating server")
	}
	server := &Server{
		listener: listener,
		handler:  handler,
	}
	server.closed.Store(false)
	go server.listen()
	return server, nil
}

func successWriter(w io.Writer, buf []byte) {
	response.WriteStatusLine(w, 200)
	headers := response.GetDefaultHeaders(len(buf))
	response.WriteHeaders(w, headers)
	w.Write(buf)
}
