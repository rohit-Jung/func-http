package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/rohit-Jung/http-protocol/internal/request"
)

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
		}

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("Error while accepting")
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HTTPVersion)

		fmt.Printf("Headers: \n")
		for key, val := range req.Headers {
			fmt.Printf("- %s: %s\n", strings.ToUpper(key), strings.ToUpper(val))
		}

		fmt.Printf("Body:\n")
		fmt.Printf("%s", string(req.Body))
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

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	out := make(chan string, 1)
// 	reader := bufio.NewReader(f)
//
// 	go func() {
// 		defer f.Close()
// 		defer close(out)
//
// 		for {
// 			line, err := reader.ReadString('\n')
// 			if err != nil {
// 				if len(line) > 0 {
// 					out <- line
// 				}
//
// 				if errors.Is(err, io.EOF) {
// 					return
// 				}
//
// 				return
// 			}
//
// 			out <- line
// 		}
// 	}()
//
// 	return out
// }
