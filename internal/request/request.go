package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	read, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var request Request
	req := strings.Split(string(read), "\r\n")

	var subdivision = strings.Split(req[0], " ")
	if subdivision[0] != strings.ToUpper(subdivision[0]) {
		return nil, errors.New("Method not valid")
	}
	request.RequestLine.Method = subdivision[0]

	if subdivision[len(subdivision)-1] != "HTTP/1.1" {
		return nil, errors.New("Version of http not valid, only 1.1 can be used")
	}

	request.RequestLine.HttpVersion = strings.Split(subdivision[len(subdivision)-1], "/")[1]

	request.RequestLine.RequestTarget = subdivision[1]

	return &request, nil
}
