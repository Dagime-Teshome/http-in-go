package headers

import (
	"bytes"
	"errors"
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
		return 0, false, err
	}
	h.Set(key, value)
	return crlfIndex + 2, false, nil
}

func (h Headers) Set(key string, value string) {
	existingValue, exists := h[key]
	if exists {
		h[key] = existingValue + "," + value
		return
	}
	h[key] = value
}
func (h Headers) SetOVR(key string, value string) {
	_, exists := h[key]
	if exists {
		h[key] = value
		return
	}
}
func (h Headers) Get(key string) (string, error) {
	keyLower := strings.ToLower(key)
	existingValue, exists := h[keyLower]
	if exists {
		return existingValue, nil
	}
	return "", errors.New("key doesn't exist")
}

func getHeaderFromString(s string) (string, string, error) {
	colonIndex := strings.Index(s, ":")
	if colonIndex == -1 {
		return "", "", errors.New("Empty header")
	}
	key := s[:colonIndex]
	value := s[colonIndex+1:]

	if strings.Contains(key, " ") {
		return "", "", errors.New("Bad Request")
	}
	if len(key) < 1 {
		return "", "", errors.New("Invalid length for key")
	}
	if !Validate(key) {
		return "", "", errors.New("Invalid character used in key")
	}
	value = strings.TrimSpace(value)
	key = strings.ToLower(key)
	return key, value, nil

}

func Validate(s string) bool {
	specialChars := "!#$%&'*+-.^_`|~"
	for i := 0; i < len(s); i++ {
		c := s[i]

		// A–Z, a–z, 0–9
		if (c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') {
			continue
		}

		// Special characters via map
		if strings.Contains(specialChars, string(c)) {
			continue
		}

		return false
	}
	return true
}
