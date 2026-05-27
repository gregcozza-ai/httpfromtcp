package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	
	"github.com/gregcozza-ai/httpfromtcp/internal/response"
	"github.com/gregcozza-ai/httpfromtcp/internal/request"
)

type Handler func(w *response.Writer, req *request.Request) 

// Server is an HTTP 1.1 server 
type Server struct {
	handler 	Handler 
	listener 	net.Listener
	closed		atomic.Bool 
	
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err 
	}
	s := &Server{
		handler: handler,
		listener: listener, 
	}
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
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusCodeBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.handler(w, req)
}
