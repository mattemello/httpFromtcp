package request

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/mattemello/httpFromtcp/internal/headers"
)

const buffSize = 8

var CRLF = "\r\n"

type status int

const (
	requestStateParsingLine status = iota + 1
	requestStateParsingHeaders
	requestStateParsingBody
	done
)

type Request struct {
	RequestLine       RequestLine
	Headers           headers.Headers
	Body              []byte
	statusRequestLine status
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	var request Request
	request.Headers = headers.NewHeaders()
	request.statusRequestLine = requestStateParsingLine

	var buf = make([]byte, buffSize, buffSize)

	var readIndex = 0
	var parseIndex = 0
	sizeBuf := buffSize

	for request.statusRequestLine != done {

		dim, err := reader.Read(buf[readIndex:])
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			request.statusRequestLine = done
			leng, _ := strconv.Atoi(request.Headers.Get("Content-Length"))
			if leng > len(request.Body) {
				return &request, errors.New("Invalid lengh, the lenght in the header is greater of the lenght of the body")
			}
			break
		}
		readIndex += dim

		parsedCount, err := request.parse(buf[parseIndex:readIndex])
		if err != nil {
			return nil, err
		}

		if request.statusRequestLine == done {
			break
		}

		parseIndex += parsedCount

		if readIndex >= buffSize {
			sizeBuf += 10
			newBuf := make([]byte, sizeBuf, sizeBuf)
			copy(newBuf, buf)
			buf = newBuf
		}

		leng, _ := strconv.Atoi(request.Headers.Get("Content-Length"))

		if leng != 0 && leng == len(request.Body) {
			request.statusRequestLine = done
			break
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

	if !strings.Contains(subdivision[1], "/") {
		return 0, errors.New("No target founded")
	}
	r.RequestTarget = subdivision[1]

	return len([]byte(req)), nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	var n int
	var err error
	switch r.statusRequestLine {
	case requestStateParsingLine:
		if !strings.Contains(string(data), CRLF) {
			return 0, nil
		}

		request := strings.Split(string(data), CRLF)

		n, err = r.RequestLine.parseRequestLine(request[0])

		if err != nil {
			return 0, err
		} else if n == 0 {
			return 0, nil
		}
		n += 2

		r.statusRequestLine = requestStateParsingHeaders
		break

	case requestStateParsingHeaders:
		var parsedAll bool

		n, parsedAll, err = r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if parsedAll {
			r.statusRequestLine = requestStateParsingBody
		}

		break

	case requestStateParsingBody:
		leng := r.Headers.Get("Content-Length")

		if leng == "" {
			r.statusRequestLine = done
			break
		}

		request := string(data)

		dim, err := strconv.Atoi(leng)
		if err != nil {
			return 0, err
		}

		if dim != len(request) {
			break
		} else if dim < len(request) {
			return 0, errors.New("dimension not valid")
		}

		r.Body = append(r.Body, []byte(request)...)

		if len(r.Body) == dim {
			r.statusRequestLine = done
		}

		break

	case done:
		return 0, errors.New("error: trying to read data in a done state")

	default:
		return 0, errors.New("error: unkown state")

	}

	return n, nil
}

func (r *Request) parse(data []byte) (int, error) {

	var totalByteParsed = 0
	for r.statusRequestLine != done {
		n, err := r.parseSingle(data[totalByteParsed:])
		if err != nil {
			return 0, err
		}

		if n == 0 {
			break
		}

		totalByteParsed += n
	}

	return totalByteParsed, nil
}
