package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ch := getLinesChannel(file)

	for ln := range ch {
		fmt.Printf("read: %s", ln)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	prov := make(chan string)

	go func() {
		var byt []byte
		byt = make([]byte, 8)
		var current string
		for {
			n, err := f.Read(byt)
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

		close(prov)
	}()

	return prov
}
