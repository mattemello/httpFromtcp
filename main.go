package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var port = ":42069"

func main() {

	liss, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		conn, err := liss.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Connection Accepted!!", conn.RemoteAddr())

		ch := getLinesChannel(conn)

		for ln := range ch {
			fmt.Printf("%s", ln)
		}
		conn.Close()
		fmt.Println("Connection Closed!!")
	}

}

func getLinesChannel(liss io.ReadCloser) <-chan string {
	prov := make(chan string)

	go func() {
		var byt []byte
		byt = make([]byte, 8)
		var current string
		for {
			n, err := liss.Read(byt)
			if err != nil || n == 0 {
				if current != "" {
					prov <- current
					current = ""
				}
				break
			}

			str := string(byt[:n])
			part := strings.Split(str, "\n")

			for i := 0; i < len(part)-1; i++ {
				prov <- current + part[i]
				current = ""
			}

			current += part[len(part)-1]
		}
	}()

	return prov
}
