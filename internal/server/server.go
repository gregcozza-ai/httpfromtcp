package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener 	net.Listener
	closed		atomic.Bool 
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err 
	}
	s := &Server{listener: listener}
	go s.listen()
	return s, nil 
}

func (s *Server) Close() error {
	if s.closed.Swap(true) {
		return nil 
	}
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		if s.closed.Load() {
			break
		}
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.closed.Load() {
				log.Printf("error accepting connections: %v", err)
			}
			return
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello World!\n"
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Printf("error writing response: %v", err)
	}
}