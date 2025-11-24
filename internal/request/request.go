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

const bufferSize = 8

const (
	initialized Status = iota
	done
)

type Request struct {
	RequestLine RequestLine
	Status      Status
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	var noBytes int
	if r.Status == initialized {
		requestLine, noBytes, err := parseRequestLine(data)
		if err != nil {
			return 0, errors.New("Error parsing stream")
		}
		if noBytes == 0 {
			return 0, nil
		} else if noBytes > 0 {
			r.RequestLine = *requestLine
			r.Status = done

		}
	} else if r.Status == done {
		return 0, errors.New("Trying to read from done parser")
	} else {
		return 0, errors.New("unknown state")
	}
	return noBytes, nil
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	// request, err := io.ReadAll(reader);
	// reqLine, err := parseRequestLine(request)

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
		numbBytesParsed, err := request.parse(buf)
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numbBytesParsed:readToIndex])
		readToIndex -= numbBytesParsed
	}

	return &request, nil
}

func parseRequestLine(request []byte) (*RequestLine, int, error) {
	crlfIndex := bytes.Index(request, []byte(crlf))
	if crlfIndex == -1 {
		return &RequestLine{}, 0, nil
	}
	requestLineText := string(request[:crlfIndex])
	requestLine, err := requestLineFromString(requestLineText)
	fmt.Println(requestLine)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, crlfIndex + 2, nil

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
