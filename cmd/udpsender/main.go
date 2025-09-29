package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")

	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		content, err := reader.ReadString('\n')

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println(err)
			continue
		}

		_, err = conn.Write([]byte(content))

		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
