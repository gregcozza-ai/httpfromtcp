package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/gregcozza-ai/httpfromtcp/internal/response"
	
)
// Server is an HTTP 1.1 server 
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

	// Write status line
	response.WriteStatusLine(conn, response.StatusCodeSuccess)
	
	// Get default headers with correct content length
	headers := response.GetDefaultHeaders(0)

	// Write headers
	if err := response.WriteHeaders(conn, headers); err != nil {
		log.Printf("error writing headers: %v", err)
		return
	}

}