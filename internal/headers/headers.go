package headers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/misterlister/httpfromtcp/internal/request"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	endl := bytes.Index(data, []byte(request.Crlf))
	if endl == -1 {
		return 0, false, nil
	}

	if endl == 0 {
		return len(request.Crlf), true, nil
	}

	sep := bytes.Index(data, []byte(":"))
	if sep == -1 {
		return 0, false, fmt.Errorf("error: invalid key:value pair. '%s' doesn't contain the ':' separator", string(data))
	}

	key := string(data[:sep])
	val := string(data[sep+1 : endl])

	if key[sep-1] == ' ' {
		return 0, false, fmt.Errorf("error: invalid key. '%s' ends in whitespace", key)
	}

	h[strings.TrimSpace(key)] = strings.TrimSpace(val)

	return endl + len(request.Crlf), false, nil
}
