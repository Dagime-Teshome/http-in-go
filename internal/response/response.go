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

type Writer struct {
	ResWriter io.Writer
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	err := WriteStatusLine(w.ResWriter, statusCode)
	if err != nil {
		return err
	}
	return nil
}
func (w *Writer) WriteHeaders(headers headers.Headers) error {

	err := WriteHeaders(w.ResWriter, headers)
	if err != nil {
		return err
	}
	return nil
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	if len(p) <= 0 {
		return 0, fmt.Errorf("Empty body write")
	}
	_, err := w.ResWriter.Write(p)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	return len(p), nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {

	chunkSize := len(p)

	nTotal := 0
	n, err := fmt.Fprintf(w.ResWriter, "%x\r\n", chunkSize)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.ResWriter.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.ResWriter.Write([]byte("\r\n"))
	if err != nil {
		return nTotal, err
	}
	nTotal += n
	return nTotal, nil
}
func (w *Writer) WriteChunkedBodyDone() (int, error) {

	n, err := w.ResWriter.Write([]byte("0\r\n"))
	if err != nil {
		return n, err
	}
	return n, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	headerString := ""
	for key, value := range h {
		headerString += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	headerString += "\r\n"
	_, err := w.ResWriter.Write([]byte(headerString))
	if err != nil {
		return err
	}
	return nil
}
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
		headerString += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	headerString += "\r\n"
	_, err := w.Write([]byte(headerString))
	if err != nil {
		return err
	}
	return nil
}
