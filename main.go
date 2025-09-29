package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")

	if err != nil {
		fmt.Println(err)
	}

	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {

		var currentLine string

		data := make([]byte, 8)

		defer close(ch)
		defer f.Close()
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
