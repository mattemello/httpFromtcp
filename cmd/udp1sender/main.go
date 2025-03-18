package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpp, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	udpConn, err := net.DialUDP("udp", nil, udpp)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer udpConn.Close()

	fmt.Printf("lissening to %s type or ctrl-c to exit\n", udpp.AddrPort())

	var buff = bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		str, err := buff.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}

		fmt.Println(str)

		_, err = udpConn.Write([]byte(str))
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
	}
}
