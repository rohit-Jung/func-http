package main 

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	url, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("Error while resolving addr")
		return
	}

	conn, err := net.DialUDP("udp", nil, url)
	if err != nil {
		log.Fatal("Error while dialing udp")
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error while reading from user")
			return
		}
		conn.Write([]byte(input))
	}
}
