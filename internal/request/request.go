package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Status int

const (
	initialized Status = iota
	done
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Status      Status
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	var bytesParsed int
	switch r.Status {
	case initialized:

		requestLine, bytesParsed, err := parseRequestLine(data)
		if err != nil {
			return 0, errors.New("Error parsing stream")
		}
		if bytesParsed == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine

	case done:
		return 0, errors.New("Trying to read from done parser")
	default:
		return 0, errors.New("unknown state")
	}
	return bytesParsed, nil
}

const bufferSize = 8
const crlf = "\r\n"
const headerEnd = "\r\n\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	request := Request{
		RequestLine: RequestLine{},
		Status:      initialized,
	}

	for request.Status != done {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				request.Status = done
				break
			}

			return nil, err
		}
		readToIndex += numBytesRead
		numBytesParsed, err := request.parse(buf)
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:readToIndex])
		readToIndex -= numBytesParsed
	}

	return &request, nil
}

func parseRequestLine(request []byte) (*RequestLine, int, error) {
	crlfIndex := bytes.Index(request, []byte(crlf))
	if crlfIndex == -1 {
		return nil, 0, nil
	}
	requestLineText := string(request[:crlfIndex])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	// the +2 is to take into account the \r\n characters
	bytesParsed := crlfIndex + 2
	return requestLine, bytesParsed, nil

}

func isAllCapsAlpha(s string) bool {
	var onlyCaps = regexp.MustCompile(`^[A-Z]+$`)
	return onlyCaps.MatchString(s)
}

func checkHttpVersion(version string) bool {

	if version == "1.1" {
		return true
	}
	return false
}

func requestLineFromString(reqLine string) (*RequestLine, error) {
	reqParts := strings.Split(reqLine, " ")

	if len(reqParts) != 3 {
		return nil, errors.New("invalid request line")
	}

	if !isAllCapsAlpha(reqParts[0]) {
		return nil, errors.New("Http method not valid")
	}

	httpVersionParts := strings.Split(reqParts[2], "/")

	if len(httpVersionParts) != 2 {
		return nil, errors.New("Malformed http version")
	}

	if !checkHttpVersion(httpVersionParts[1]) {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpVersionParts[1])
	}

	if httpVersionParts[0] != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpVersionParts[0])
	}
	return &RequestLine{
		HttpVersion:   httpVersionParts[1],
		RequestTarget: reqParts[1],
		Method:        reqParts[0],
	}, nil
}
