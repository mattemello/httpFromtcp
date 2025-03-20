package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var buffSize = 8
var CRLF = "\r\n"

type status int

const (
	intialized status = iota + 1
	done
)

type Request struct {
	RequestLine RequestLine
	status      status
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	var request Request
	request.status = intialized

	var buf = make([]byte, buffSize, buffSize)

	var readIndex = 0

	for request.status != done {

		dim, err := reader.Read(buf[readIndex:])
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			request.status = done
			fmt.Println(request.RequestLine)
			break
		}

		_, err = request.parse(buf)
		if err != nil {
			return nil, err
		}

		if request.status == done {
			break
		}

		readIndex += dim

		if readIndex >= buffSize {
			buffSize *= 2
			newBuf := make([]byte, buffSize, buffSize)
			copy(newBuf, buf)
			buf = newBuf
		}

	}
	return &request, nil
}

func (r *RequestLine) parseRequestLine(req string) (int, error) {

	var subdivision = strings.Split(req, " ")

	if subdivision[0] != strings.ToUpper(subdivision[0]) {
		return 0, errors.New("Method not valid")
	}

	r.Method = subdivision[0]

	if subdivision[len(subdivision)-1] != "HTTP/1.1" {
		return 0, errors.New("Version of http not valid, only 1.1 can be used")
	}

	r.HttpVersion = strings.Split(subdivision[len(subdivision)-1], "/")[1]

	fmt.Println(subdivision[1])
	if !strings.Contains(subdivision[1], "/") {
		return 0, errors.New("No target founded")
	}
	r.RequestTarget = subdivision[1]

	return len([]byte(req)), nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.status == intialized {
		if !strings.Contains(string(data), CRLF) {
			return 0, nil
		}

		request := strings.Split(string(data), CRLF)

		fmt.Println(request[0])
		n, err := r.RequestLine.parseRequestLine(request[0])

		if err != nil {
			return 0, err
		} else if n == 0 {
			return 0, nil
		}

		r.status = done
	} else if r.status == done {
		return 0, errors.New("error: trying to read data in a done state")
	} else {

		return 0, errors.New("error: unkown state")
	}

	return 0, nil
}
