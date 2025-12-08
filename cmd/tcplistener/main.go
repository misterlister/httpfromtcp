package main

import (
	"fmt"
	"log"
	"net"

	"github.com/misterlister/httpfromtcp/internal/request"
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
			log.Fatalf("error: %s\n", err.Error())
		}

		fmt.Println("A connection has been accepted!")

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error parsing request: %s\n", err.Error())
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)

		if len(request.Headers) > 0 {
			fmt.Println("Headers:")

			for key, val := range request.Headers {
				fmt.Printf("- %s: %s\n", key, val)
			}
		}

	}
}
