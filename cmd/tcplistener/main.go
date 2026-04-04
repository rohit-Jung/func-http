package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)
	reader := bufio.NewReader(f)

	go func() {
		defer f.Close()
		defer close(out)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if len(line) > 0 {
					out <- line
				}

				if errors.Is(err, io.EOF) {
					return
				}

				return
			}

			out <- line
		}
	}()

	return out
}

func main() {
	net, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("Error while Listening")
		return
	}

	defer net.Close()

	for {
		conn, err := net.Accept()
		if err != nil {
			log.Fatal("Error while accepting")
			return
		}

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("%s", line)
		}
	}
}

// data := make([]byte, 8)
// n, err := file.Read(data)
// data = data[:n]
// if i := bytes.IndexByte(data, '\n'); i != -1 {
// 	str += string(data[:i])
// 	data = data[i+1:]
// 	fmt.Printf("read: %s\n", str)
// 	str = ""
// }
//
// str += string(data)
