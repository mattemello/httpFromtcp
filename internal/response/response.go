package response

import (
	"io"
	"strconv"

	"github.com/mattemello/httpFromtcp/internal/headers"
)

type StatusCode int

const (
	Ok          StatusCode = 200
	BadRequest             = 400
	ServerError            = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var message []byte
	switch statusCode {
	case Ok:
		message = []byte("HTTP/1.1 200 OK\r\n")
		break
	case BadRequest:
		message = []byte("HTTP/1.1 400 Bad Request\r\n")
		break
	case ServerError:
		message = []byte("HTTP/1.1 500 Internal Server Error\r\n")
		break
	}

	_, err := w.Write(message)

	return err
}

func GetDefaultHeaders(contentLeng int) headers.Headers {

	headers := headers.NewHeaders()

	headers["content-lenght"] = strconv.Itoa(contentLeng)
	headers["connection"] = "close"
	headers["content-type"] = "text/plain"

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	for key, elem := range headers {
		_, err := w.Write([]byte(key + ":" + elem + "\r\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
