package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type ServerStatus int

type Server struct {
	listener   net.Listener
	connection net.Conn
	closed     atomic.Bool
}

const (
	Listening ServerStatus = iota
	Closed
)

func (s *Server) Close() {
	s.closed.Store(true)
	s.listener.Close()
	s.connection.Close()
}

func (s *Server) listen() {

	if !s.closed.Load() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				log.Fatal("error accepting connection : \n", err)
			}
			s.connection = conn
			go s.handle(conn)
		}
	}

}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	resp := []byte("HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello World!")
	conn.Write(resp)
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
