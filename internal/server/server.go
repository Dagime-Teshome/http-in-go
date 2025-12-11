package server

import (
	"MODULE_NAME/internal/response"
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type ServerStatus int

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

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
	response.WriteStatusLine(conn, 200)
	headers := response.GetDefaultHeaders(0)
	response.WriteHeaders(conn, headers)
	// resp := []byte("HTTP/1.1 200 OK\r\n" +
	// 	"Content-Type: text/plain\r\n" +
	// 	"Content-Length: 13\r\n" +
	// 	"\r\n" +
	// 	"Hello World!")
	// conn.Write(resp)
}
func Serve(port int) (*Server, error) {
	portString := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", portString)

	if err != nil {
		return nil, fmt.Errorf("error creating server")
	}
	server := &Server{
		listener: listener,
	}
	server.closed.Store(false)
	go server.listen()
	return server, nil
}
