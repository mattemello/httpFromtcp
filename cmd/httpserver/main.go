package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/mattemello/httpFromtcp/internal/headers"
	"github.com/mattemello/httpFromtcp/internal/request"
	"github.com/mattemello/httpFromtcp/internal/response"
	"github.com/mattemello/httpFromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		destination := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")

		w.WriteStatusLine(response.Ok)
		header := headers.NewHeaders()
		header.Add("Connection", "close")
		header.Add("Content-Type", "text/plain")
		header.Add("Transfer-Encoding", "chunked")
		header.Add("Trailer", "X-Content-SHA256, X-Content-Length")
		w.WriteHeaders(header)

		trailer := headers.NewHeaders()

		resp, err := http.Get("https://httpbin.org" + destination)
		if err != nil {
			os.Exit(1)
		}

		var buff = make([]byte, 1024)
		var allBuff = make([]byte, 0)

		for {
			n, err := resp.Body.Read(buff)
			if err != nil {
				if err == io.EOF {
					hasMessage := sha256.Sum256(allBuff)
					trailer.Add("X-Content-SHA256", fmt.Sprintf("%x", hasMessage))
					trailer.Add("X-Content-Length", strconv.Itoa(len(allBuff)))

					w.WriteTrailers(trailer)

					return
				}
				os.Exit(1)
			}

			if n == 0 {
				hasMessage := sha256.Sum256(allBuff)
				trailer.Add("X-Content-SHA256", fmt.Sprintf("%x", hasMessage))
				trailer.Add("X-Content-Length", strconv.Itoa(len(allBuff)))

				w.WriteTrailers(trailer)

				return

			}

			_, err = w.WriteChunkBody(buff[:n])

			copyBuf := allBuff
			allBuff = make([]byte, len(allBuff)+n)
			copy(allBuff, copyBuf)
			allBuff = append(allBuff, buff[:n]...)

		}

	}

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {

		file, err := os.ReadFile("./assets/vim.mp4")
		if err != nil {
			log.Println("Can't read the file", err)
			return
		}

		w.WriteStatusLine(response.Ok)
		header := headers.NewHeaders()
		header.Add("Connection", "video/mp4")
		header.Add("Content-Type", "text/plain")
		header.Add("Content-Length", strconv.Itoa(len(file)))
		w.WriteHeaders(header)

		w.WriteBody(file)

		return
	}

	const field = "Content-Type"
	const value = "text-html"

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		errStatusLine := w.WriteStatusLine(response.BadRequest)
		if errStatusLine != nil {
			log.Println("Error in the status Line", errStatusLine)
		}

		header := headers.NewHeaders()
		header = response.GetDefaultHeaders(len(yourproblem))
		header.Add(field, value)

		errHeader := w.WriteHeaders(header)
		if errHeader != nil {
			log.Println("Error in the status Line", errHeader)
		}

		_, err := w.WriteBody([]byte(yourproblem))
		if err != nil {
			log.Println("Error in the status Line", err)
		}
		break

	case "/myproblem":
		errStatusLine := w.WriteStatusLine(response.ServerError)
		if errStatusLine != nil {
			log.Println("Error in the status Line", errStatusLine)
		}

		header := headers.NewHeaders()
		header = response.GetDefaultHeaders(len(myproblem))
		header.Add(field, value)

		errHeader := w.WriteHeaders(header)
		if errHeader != nil {
			log.Println("Error in the status Line", errHeader)
		}

		_, err := w.WriteBody([]byte(myproblem))
		if err != nil {
			log.Println("Error in the status Line", err)
		}
		break
	default:
		errStatusLine := w.WriteStatusLine(response.Ok)
		if errStatusLine != nil {
			log.Println("Error in the status Line", errStatusLine)
		}

		header := headers.NewHeaders()
		header = response.GetDefaultHeaders(len(allRight))
		header.Add(field, value)

		errHeader := w.WriteHeaders(header)
		if errHeader != nil {
			log.Println("Error in the status Line", errHeader)
		}

		_, err := w.WriteBody([]byte(allRight))
		if err != nil {
			log.Println("Error in the status Line", err)
		}
		break

	}

}

const yourproblem = `
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

const myproblem = `
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`

const allRight = `
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`
