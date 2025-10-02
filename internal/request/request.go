package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Status int

const (
	Initialized Status = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	State       Status
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"
const BufferSize = 8

func (r *Request) parse(data []byte) (int, error) {
	if r.State == Initialized {
		reqLine, bytesRead, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if bytesRead == 0 {
			return 0, nil
		}

		r.RequestLine = *reqLine
		r.State = Done
		return bytesRead, nil
	} else if r.State == Done {
		return 0, fmt.Errorf("error: trying to read data in a done state")
	}
	return 0, fmt.Errorf("error: unknown state")
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, BufferSize)

	readToIndex := 0

	req := Request{
		State: Initialized,
	}

	for req.State != Done {
		bufLen := len(buf)
		if readToIndex >= bufLen {
			bufLen *= 2
			newBuf := make([]byte, bufLen)
			copy(newBuf, buf)
			buf = newBuf
		}

		bytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, err
			}
			req.State = Done
			break
		}
		readToIndex += bytesRead
		bytesParsed, err := req.parse(buf)
		if err != nil {
			return nil, err
		}

		unparsedBuf := make([]byte, bufLen-bytesParsed)
		copy(unparsedBuf, buf[bytesParsed:])
		buf = unparsedBuf
		readToIndex -= bytesParsed
	}

	return &req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	index := bytes.Index(data, []byte(crlf))
	if index == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:index])

	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, index + len(crlf), nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	httpPrefix := "HTTP/"
	validVersion := "1.1"

	parts := strings.Split(str, " ")

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid number of sections (%d) in request line: %s", len(parts), str)
	}

	for _, char := range parts[0] {
		if !unicode.IsUpper(char) {
			return nil, fmt.Errorf("invalid http method: %s", parts[0])
		}
	}

	method := parts[0]
	requestTarget := parts[1]

	if !strings.HasPrefix(parts[2], httpPrefix) {
		return nil, fmt.Errorf("invalid http version declaration: %s", parts[2])
	}

	version := strings.TrimPrefix(parts[2], httpPrefix)

	if version != validVersion {
		return nil, fmt.Errorf("invalid http version: %s", version)
	}

	httpVersion := version

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   httpVersion,
	}, nil
}
