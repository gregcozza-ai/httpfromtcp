package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

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
	file, err := os.Open(inputFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("======================================")
	
	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
