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
	Handler    Handler
}

func (errHand HandlerError) WriteHandlerError(w io.Writer) {
	defHeader := response.GetDefaultHeaders(len(errHand.Message))
	response.WriteStatusLine(w, errHand.StatusCode)
	response.WriteHeaders(w, defHeader)
	w.Write([]byte(errHand.Message))
}

func Serve(port int, Hander Handler) (*Server, error) {
	portString := strconv.Itoa(port)
	liss, err := net.Listen("tcp", ":"+portString)
	if err != nil {
		return nil, err
	}

	server := &Server{Connection: liss, Port: port, on: false, Handler: Hander}

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
			if !s.on {
				return
			}
			log.Printf("Can't accept the connection: %v", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {

	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		erroH := &HandlerError{
			StatusCode: response.BadRequest,
			Message:    err.Error(),
		}

		erroH.WriteHandlerError(conn)
		return
	}

	var buf = bytes.NewBuffer([]byte{})

	erroHand := s.Handler(buf, req)
	if erroHand != nil {
		erroHand.WriteHandlerError(conn)
		return
	}

	response.WriteStatusLine(conn, response.Ok)
	defHeader := response.GetDefaultHeaders(buf.Len())
	response.WriteHeaders(conn, defHeader)
	var read = buf.Bytes()
	conn.Write(read)
}
