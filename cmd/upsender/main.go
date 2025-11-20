package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	upAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		// handle error
	}

	updDial, err := net.DialUDP("udp", nil, upAddr)
	if err != nil {
		// handle error
	}
	defer updDial.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		print("> ")
		stdString, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Something went wrong", err)
		}
		updDial.Write([]byte(stdString))
	}
}
