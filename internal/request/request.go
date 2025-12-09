package request

import (
	"MODULE_NAME/internal/headers"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Status int

const (
	initialized Status = iota
	done
	ParsingHeaders
	ParsingBody
)

type Request struct {
	RequestLine    RequestLine
	Headers        headers.Headers
	Body           []byte
	Status         Status
	bodyLengthRead int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// func (r *Request) parse(data []byte) (int, error) {
// 	switch r.Status {
// 	case initialized:
//			requestLine, startLineParsed, err := parseRequestLine(data)
//			if err != nil {
//				return 0, errors.New("error parsing stream")
//			}
//			if startLineParsed == 0 {
//				return 0, nil
//			}
//			r.RequestLine = *requestLine
//			r.Status = ParsingHeaders
//			return startLineParsed, nil
//		case ParsingHeaders:
//			headerParsed, finished, err := r.Headers.Parse(data)
//			if err != nil {
//				return 0, err
//			}
//			if finished {
//				r.Status = done
//			}
//			return headerParsed, nil
//		case done:
//			return 0, errors.New("trying to read from done parser")
//		default:
//			return 0, errors.New("unknown state")
//		}
//	}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.Status != done {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.Status {
	case initialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			// something actually went wrong
			return 0, err
		}
		if n == 0 {
			// just need more data
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.Status = ParsingHeaders
		return n, nil
	case ParsingHeaders:
		n, finished, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if finished {
			r.Status = ParsingBody
		}
		return n, nil
	case ParsingBody:
		contLen, err := r.Headers.Get("Content-Length")
		if err != nil {
			r.Status = done
			return len(data), nil
		}
		conLenInt, err := strconv.Atoi(contLen)
		if err != nil {
			return 0, err
		}
		// /

		r.Body = append(r.Body, data...)
		r.bodyLengthRead += len(data)
		if r.bodyLengthRead > conLenInt {
			return 0, fmt.Errorf("Content-Length too large")
		}
		if r.bodyLengthRead == conLenInt {
			r.Status = done
		}
		return len(data), nil

	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

const bufferSize = 8
const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, bufferSize)
	readToIndex := 0

	request := &Request{
		RequestLine: RequestLine{},
		Headers:     headers.NewHeaders(),
		Body:        []byte{},
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
			if errors.Is(err, io.EOF) {
				if request.Status != done {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", request.Status, numBytesRead)
				}
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead
		numBytesParsed, err := request.parse(buf[:readToIndex])

		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return request, nil
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
