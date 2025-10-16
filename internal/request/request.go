package request

import (
	"fmt"
	"io"
	"strings"
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

var ERROR_BAD_STARTLINE = fmt.Errorf("bad request line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var ERROR_INCOMPLETE_STARTLINE = fmt.Errorf("incomplete request line")
var SEPERATOR = "\r\n"

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case StateInit:
		case StateDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone
}

func parseRequestLine(b string) (*RequestLine, int, error) {
	idx := strings.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, 0, ERROR_INCOMPLETE_STARTLINE
	}

	startLine := b[:idx]
	restOfMessage := b[idx+len(SEPERATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, restOfMessage, ERROR_BAD_STARTLINE
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" {
		return nil, restOfMessage, ERROR_BAD_STARTLINE
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	return rl, restOfMessage, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	// NOTE: buffer could get overrun... a header that exceeds 1k would do that...
	// or the body
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])

		// TODO: what to do here?
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen+n])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
