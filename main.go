package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
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
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("Error while reading")
		return
	}

	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s", line)
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
