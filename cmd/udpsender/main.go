package main

import (
	"log"
	"net"
	"os"
	"bufio"
)

func main () {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	conn , err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			continue 
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Println("Error sending data:", err)
		}
	}
}
