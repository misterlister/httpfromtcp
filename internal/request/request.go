package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\r\n")

	if len(lines) == 0 {
		return nil, errors.New("no lines to read in request")
	}

	reqLine, err := parseRequestLine(lines[0])

	if err != nil {
		return nil, err
	}

	req := Request{
		RequestLine: reqLine,
	}

	return &req, nil
}

func parseRequestLine(reqLine string) (RequestLine, error) {
	parsedReq := RequestLine{}
	httpPrefix := "HTTP/"
	validVersion := "1.1"

	parts := strings.Split(reqLine, " ")

	if len(parts) != 3 {
		return parsedReq, errors.New("invalid number of sections in request line")
	}

	for _, char := range parts[0] {
		if !unicode.IsUpper(char) {
			return parsedReq, errors.New("invalid http method")
		}
	}

	parsedReq.Method = parts[0]
	parsedReq.RequestTarget = parts[1]

	if !strings.HasPrefix(parts[2], httpPrefix) {
		return parsedReq, errors.New("invalid http version declaration")
	}

	version := strings.TrimPrefix(parts[2], httpPrefix)

	if version != validVersion {
		return parsedReq, errors.New("invalid http version")
	}

	parsedReq.HttpVersion = version

	return parsedReq, nil
}
