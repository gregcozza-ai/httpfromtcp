package main

import (
	"fmt"
	"io"
	"net"
	
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func () {
		defer f.Close()
		buffer := make([]byte, 8)
		currentLine := ""

		for {
			n, err := f.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}

			currentLine += string(buffer[:n])
			parts := strings.Split(currentLine, "\n")

			for i := 0; i < len(parts)-1; i++ {
				ch <- parts[i]
			}

			currentLine = parts[len(parts)-1]
		}

		if currentLine != "" {
			ch <- currentLine 
		}
		close(ch)
	}()
	return ch 
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("Connection accepted")
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}
		fmt.Println("Connection closed")
	}
}
