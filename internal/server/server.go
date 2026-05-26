package server

import (
	"fmt"
	"io"
	"net"

	"github.com/rohit-Jung/http-protocol/internal/request"
	"github.com/rohit-Jung/http-protocol/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w response.Writer, req *request.Request)
type (
	Server struct {
		lsnr    net.Listener
		handler Handler
	}
)

func newServer(lsnr net.Listener, handler Handler) *Server {
	return &Server{
		lsnr,
		handler,
	}
}

func Serve(port int, handler Handler) (*Server, error) {
	lsnr, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := newServer(lsnr, handler)
	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.lsnr.Close()
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.lsnr.Accept()
		if err != nil {
			s.Close()
			return
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn io.ReadWriteCloser) {
	defer conn.Close()

	responseWriter := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	headers := response.GetDefaultHeaders(0)

	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequst)
		responseWriter.WriteHeaders(headers)
		return
	}

	s.handler(*responseWriter, req)
}
