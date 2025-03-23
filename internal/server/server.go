package server

import (
	"log"
	"net"
	"strconv"

	"github.com/mattemello/httpFromtcp/internal/response"
)

type Server struct {
	Connection net.Listener
	Port       int
	on         bool
}

func Serve(port int) (*Server, error) {
	portString := strconv.Itoa(port)
	liss, err := net.Listen("tcp", ":"+portString)
	if err != nil {
		return nil, err
	}

	server := &Server{Connection: liss, Port: port, on: false}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	err := s.Connection.Close()
	if err != nil {
		return err
	}

	s.on = false

	return nil
}

func (s *Server) listen() {
	log.Println("> Start the listening of the server")
	for !s.on {
		conn, err := s.Connection.Accept()
		if err != nil {
			s.Close()
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {

	headers := response.GetDefaultHeaders(0)
	response.WriteStatusLine(conn, 200)
	response.WriteHeaders(conn, headers)

	conn.Close()
}
