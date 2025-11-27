package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIndex := bytes.Index(data, []byte(crlf))
	if crlfIndex == -1 {
		return 0, false, nil
	} else if crlfIndex == 0 {
		return 0, true, nil
	}
	headerString := string(data[:crlfIndex])
	headerParts, err := parseHeader(headerString)

	if err != nil {
		// handle error
		return 0, false, err
	}
	h[headerParts[0]] = headerParts[1]
	return crlfIndex, false, nil
}

func parseHeader(s string) ([]string, error) {
	return
}
func getHeaderFromString(s string) ([]string, error) {
	headerParts := strings.Split(s, ":")
	if len(headerParts) != 2 {
		return nil, errors.New("invalid header length")
	}
	return headerParts, nil

}
