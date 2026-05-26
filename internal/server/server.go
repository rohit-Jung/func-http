package server

import (
	"bytes"
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

type Handler func(w io.Writer, req *request.Request) *HandlerError

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

func WriteResponse(conn io.ReadWriteCloser, statusCode response.StatusCode, message []byte) {
	response.WriteStatusLine(conn, statusCode)
	headers := response.GetDefaultHeaders(len(message))
	response.WriteHeaders(conn, headers)
	conn.Write(message)
}

func (s *Server) handle(conn io.ReadWriteCloser) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		WriteResponse(conn, response.StatusBadRequst, []byte(""))
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	handlerErr := s.handler(buffer, req)
	if handlerErr != nil {
		WriteResponse(conn, handlerErr.StatusCode, []byte(handlerErr.Message))
		return
	}

	body := buffer.Bytes()
	WriteResponse(conn, response.StatusOk, body)
}
