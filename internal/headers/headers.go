package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIndex := bytes.Index(data, []byte(crlf))
	if crlfIndex == -1 {
		return 0, false, nil
	} else if crlfIndex == 0 {
		return 2, true, nil
	}
	headerString := string(data[:crlfIndex])
	key, value, err := getHeaderFromString(headerString)

	if err != nil {
		// handle error
		return 0, false, err
	}
	h[key] = value
	fmt.Println(h, "-----------")
	return crlfIndex + 2, false, nil
}

func getHeaderFromString(s string) (string, string, error) {
	colonIndex := strings.Index(s, ":")
	key := s[:colonIndex]
	value := s[colonIndex+1:]

	if strings.Contains(key, " ") {
		return "", "", errors.New("Bad Request")
	}
	value = strings.TrimSpace(value)
	return key, value, nil

}
