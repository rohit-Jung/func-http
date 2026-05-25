package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	lsnr net.Listener
}

func newServer(lsnr net.Listener) *Server {
	return &Server{
		lsnr: lsnr,
	}
}

func Serve(port int) (*Server, error) {
	lsnr, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := newServer(lsnr)
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
	msg := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!")
	conn.Write(msg)
	conn.Close()
}
