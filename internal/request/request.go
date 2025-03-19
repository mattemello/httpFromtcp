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

	request.RequestLine, err = requestLine(req[0])
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func requestLine(req string) (RequestLine, error) {
	var subdivision = strings.Split(req, " ")

	if subdivision[0] != strings.ToUpper(subdivision[0]) {
		return RequestLine{}, errors.New("Method not valid")
	}

	var resp RequestLine
	resp.Method = subdivision[0]

	if subdivision[len(subdivision)-1] != "HTTP/1.1" {
		return RequestLine{}, errors.New("Version of http not valid, only 1.1 can be used")
	}

	resp.HttpVersion = strings.Split(subdivision[len(subdivision)-1], "/")[1]

	resp.RequestTarget = subdivision[1]

	return resp, nil
}
