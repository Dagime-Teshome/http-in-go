package response

import (
	"MODULE_NAME/internal/headers"
	"fmt"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOK            StatusCode = 200
	StatusBadRequest    StatusCode = 400
	StatusInternalError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
	case StatusBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
	case StatusInternalError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
	default:
		_, err := w.Write([]byte(""))
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	s := strconv.Itoa(contentLen)
	newHeader := headers.NewHeaders()
	newHeader.Set("Content-Length", s)
	newHeader.Set("Connection", "close")
	newHeader.Set("Content-Type", "text/plain")
	return newHeader
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	headerString := ""
	for key, value := range headers {
		headerString += fmt.Sprintf("%s:%s\r\n", key, value)
	}
	headerString += "\r\n"
	_, err := w.Write([]byte(headerString))
	if err != nil {
		return err
	}
	return nil
}
