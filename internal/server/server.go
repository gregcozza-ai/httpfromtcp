package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"bytes"
	"io"

	"github.com/gregcozza-ai/httpfromtcp/internal/response"
	"github.com/gregcozza-ai/httpfromtcp/internal/request"

	
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode 	response.StatusCode
	Message		string 
}

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(w, headers)
	w.Write(messageBytes)
}

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
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	hErr := s.handler(buf, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}
	b := buf.Bytes()
	response.WriteStatusLine(conn, response.StatusCodeSuccess)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)
	conn.Write(b)
	return
}