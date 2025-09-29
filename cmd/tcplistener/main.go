package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("A connection has been accepted!")

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Printf("%s\n", line)
		}

		conn.Close()
		fmt.Println("The connection has been closed!")
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {

		var currentLine string

		data := make([]byte, 8)

		defer close(ch)
		for {
			bytesRead, err := f.Read(data)

			if err != nil {
				if errors.Is(err, io.EOF) {
					if currentLine != "" {
						ch <- currentLine
					}
				} else {
					fmt.Printf("error: %s\n", err)
				}
				break
			}

			parts := strings.Split(string(data[:bytesRead]), "\n")

			for i := 0; i < len(parts)-1; i++ {
				currentLine += parts[i]
				ch <- currentLine
				currentLine = ""
			}

			currentLine += parts[len(parts)-1]
		}
	}()
	return ch
}
