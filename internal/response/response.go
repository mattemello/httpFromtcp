package response

import (
	"errors"
	"fmt"
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
	default:
		return errors.New("Invalid status line")
	}

	_, err := w.Write(message)

	return err
}

func GetDefaultHeaders(contentLeng int) headers.Headers {

	headers := headers.NewHeaders()

	headers["Content-Lenght"] = strconv.Itoa(contentLeng)
	headers["Connection"] = "close"
	headers["Content-Type"] = "text/plain"

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	for key, elem := range headers {
		_, err := w.Write([]byte(key + ":" + elem + "\r\n"))
		if err != nil {
			return err
		}
	}

	w.Write([]byte("\r\n"))

	return nil
}

type timeStatus int

const (
	statusLineTime timeStatus = iota + 10
	headersTime
	bodyTime
	doneTime
)

type Writer struct {
	Conn        io.Writer
	writeStatus timeStatus
}

func NewWriter() *Writer {
	return &Writer{
		writeStatus: statusLineTime,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writeStatus != statusLineTime {
		return errors.New("Can't do the status line, you can't access the writer in an incorrect order")
	}

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
	default:
		return errors.New("Invalid status line")
	}

	_, err := w.Conn.Write(message)

	w.writeStatus = headersTime

	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writeStatus != headersTime {
		return errors.New("Incorrect order of the writer, before doing `WriteHeaders` you need to do `WriteStatusLine`")
	}

	for key, elem := range headers {
		_, err := w.Conn.Write([]byte(key + ":" + elem + "\r\n"))
		if err != nil {
			return err
		}
	}

	_, err := w.Conn.Write([]byte("\r\n"))

	w.writeStatus = bodyTime

	return err
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	if w.writeStatus != bodyTime {
		return 0, errors.New("Incorrect order of the writer, before doing `WriteBody` you need to do `WriteHeaders`")
	}

	n, err := w.Conn.Write(body)

	w.writeStatus = doneTime

	return n, err
}

func (w *Writer) WriteChunkBody(p []byte) (int, error) {
	if w.writeStatus != bodyTime {
		return 0, errors.New("Incorrect order of the writer, before doing `WriteBody` you need to do `WriteHeaders`")
	}

	hexNumber := fmt.Sprintf("%x", len(p))
	n1, _ := w.Conn.Write([]byte(hexNumber + "/r/n"))
	n, err := w.Conn.Write(append(p, []byte("/r/n")...))

	return n + n1, err
}

func (w *Writer) WriteChunkBodyDone() (int, error) {
	if w.writeStatus != bodyTime {
		return 0, errors.New("Incorrect order of the writer, before doing `WriteBody` you need to do `WriteHeaders`")
	}

	hexNumber := fmt.Sprintf("%x", 0)
	n1, _ := w.Conn.Write([]byte(hexNumber + "/r/n"))
	n, err := w.Conn.Write([]byte("/r/n"))

	return n + n1, err
}
