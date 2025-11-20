package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(file io.ReadCloser) <-chan string {
	c := make(chan string)

	go func() {
		currentLine := ""
		defer file.Close()
		defer close(c)
		for {
			fileContent := make([]byte, 8)
			n, fileerror := file.Read(fileContent)
			if fileerror != nil {
				if errors.Is(fileerror, io.EOF) {
					// fmt.Printf("read: %s\n", currentLine)
					c <- currentLine
					break
				}
				fmt.Printf("error: %s\n", fileerror.Error())
				break
			}
			str := string(fileContent[:n])
			parts := strings.Split(str, "\n")

			if len(parts) == 1 {
				currentLine += string(parts[0])
			} else {
				currentLine += parts[0]
				// fmt.Printf("read: %s\n", currentLine)
				c <- currentLine
				currentLine = ""
				currentLine = parts[1]
			}
			// fmt.Printf("read: %s\n", currentLine)
			// fmt.Printf("read: %s\n", str)
			// fmt.Printf("read: %s\n", string(fileContent[:n]))
		}
	}()
	return c
}

func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted Connection From", connection.RemoteAddr())
		line := getLinesChannel(connection)
		for line := range line {
			fmt.Println(line)
		}
		fmt.Println("Connection to ", connection.RemoteAddr(), "closed")

	}
}
