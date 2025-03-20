package main

import (
	"fmt"
	"github.com/mattemello/httpFromtcp/internal/request"
	"log"
	"net"
)

var port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		requestLine, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}

		fmt.Printf("\nRequest line:\n- Method: %s\n- Target: %s\n- Version: %s\n\n",
			requestLine.RequestLine.Method,
			requestLine.RequestLine.RequestTarget,
			requestLine.RequestLine.HttpVersion)

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	lines := make(chan string)
// 	go func() {
// 		defer f.Close()
// 		defer close(lines)
// 		currentLineContents := ""
// 		for {
// 			b := make([]byte, 8, 8)
// 			n, err := f.Read(b)
// 			if err != nil {
// 				if currentLineContents != "" {
// 					lines <- currentLineContents
// 				}
// 				if errors.Is(err, io.EOF) {
// 					break
// 				}
// 				fmt.Printf("error: %s\n", err.Error())
// 				return
// 			}
// 			str := string(b[:n])
// 			parts := strings.Split(str, "\n")
// 			for i := 0; i < len(parts)-1; i++ {
// 				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
// 				currentLineContents = ""
// 			}
// 			currentLineContents += parts[len(parts)-1]
// 		}
// 	}()
// 	return lines
// }
