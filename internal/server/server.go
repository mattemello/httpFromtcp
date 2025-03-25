package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/mattemello/httpFromtcp/internal/request"
	"github.com/mattemello/httpFromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Server struct {
	Connection net.Listener
	Port       int
	on         bool
}

func WriteHandlerError(w io.Writer)

func Serve(port int, Hander Handler) (*Server, error) {
	portString := strconv.Itoa(port)
	liss, err := net.Listen("tcp", ":"+portString)
	if err != nil {
		return nil, err
	}

	server := &Server{Connection: liss, Port: port, on: false}

	go server.listen(Hander)

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

func (s *Server) listen(handler Handler) {
	log.Println("> Start the listening of the server")
	for !s.on {
		conn, err := s.Connection.Accept()
		if err != nil {
			s.Close()
		}

		go s.handle(conn, handler)
	}
}

func (s *Server) handle(conn net.Conn, handler Handler) {

	req, err := request.RequestFromReader(conn)
	if err != nil {
		//todo: respond with an error
	}

	var buffer bytes.Buffer

	erro := handler(&buffer, req)
	if erro != nil {
		//todo: respond with an error
	}

	defHandler := response.GetDefaultHeaders()

	conn.Close()
}
