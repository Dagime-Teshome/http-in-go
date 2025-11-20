package request

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
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
	request, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}
	reqLine, err := parseRequestLine(request)

	if err != nil {
		return nil, err
	}
	var req Request
	req.RequestLine = *reqLine
	return &req, nil
}

func parseRequestLine(request []byte) (*RequestLine, error) {
	var RequestLine RequestLine
	fullRequest := string(request)
	requestLine := strings.Split(strings.Split(fullRequest, "\r\n")[0], " ")
	fmt.Println(len(requestLine))
	if len(requestLine) != 3 {
		return nil, errors.New("invalid number of requests")
	}

	if !isAllCapsAlpha(requestLine[0]) {
		return nil, errors.New("Http method not valid")
	}
	versionNum := strings.Split(requestLine[2], "/")[1]
	if !veryifyVersion(versionNum) {
		return nil, errors.New("invalid version number")
	}
	RequestLine.HttpVersion = versionNum
	RequestLine.RequestTarget = requestLine[1]
	RequestLine.Method = requestLine[0]
	return &RequestLine, nil

}

func isAllCapsAlpha(s string) bool {
	var onlyCaps = regexp.MustCompile(`^[A-Z]+$`)
	return onlyCaps.MatchString(s)
}

func veryifyVersion(version string) bool {

	if version == "1.1" {
		return true
	}
	return false
}
