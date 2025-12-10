package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/misterlister/httpfromtcp/internal/consts"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func AddHeaderValue(h Headers, key, val string) {
	_, keyExists := h[key]

	if !keyExists {
		h[key] = val
	} else {
		h[key] = h[key] + ", " + val
	}
}

var ValidHeaderCharacters = map[rune]bool{
	'!':  true,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'*':  true,
	'+':  true,
	'-':  true,
	'.':  true,
	'^':  true,
	'_':  true,
	'`':  true,
	'|':  true,
	'~':  true,
}

func (h Headers) Get(key string) (string, bool) {
	value := h[strings.ToLower(key)]
	valid := value != ""
	return value, valid
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	endl := bytes.Index(data, []byte(consts.Crlf))
	if endl == -1 {
		return 0, false, nil
	}

	if endl == 0 {
		return len(consts.Crlf), true, nil
	}

	sep := bytes.Index(data, []byte(":"))
	if sep == -1 {
		return 0, false, fmt.Errorf("error: invalid 'key: value' pair. '%s' doesn't contain the ':' separator", string(data))
	}

	keyString := string(data[:sep])
	valString := string(data[sep+1 : endl])

	if len(keyString) == 0 {
		return 0, false, fmt.Errorf("error: invalid key. Key cannot be empty")
	}

	if keyString[len(keyString)-1] == ' ' {
		return 0, false, fmt.Errorf("error: invalid key. '%s' ends in whitespace", keyString)
	}

	key := strings.ToLower(strings.TrimSpace(keyString))
	if !isValidHeaderString(key) {
		return 0, false, fmt.Errorf("error: invalid key. '%s' contains invalid characters", key)
	}

	val := strings.TrimSpace(valString)

	AddHeaderValue(h, key, val)

	return endl + len(consts.Crlf), false, nil
}

func isValidHeaderString(s string) bool {
	for _, r := range s {
		switch {
		case unicode.IsLetter(r):
			continue
		case unicode.IsDigit(r):
			continue
		case ValidHeaderCharacters[r]:
			continue
		default:
			return false
		}
	}
	return true
}
