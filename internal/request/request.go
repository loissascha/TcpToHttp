package request

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"tcpToHttp/internal/headers"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateDone    parserState = "done"
	StateError   parserState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
	Headers     *headers.Headers
	Body        string
}

func (r *Request) hasBody() bool {
	// TODO: when doing chunked encoding, update this method
	length := getIntHeader(r.Headers, "content-length", 0)
	return length != 0
}

func getIntHeader(headers *headers.Headers, name string, defaultValue int) int {
	valueStr, exists := headers.Get(name)
	if !exists {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func newRequest() *Request {
	h := headers.NewHeaders()
	return &Request{
		state:   StateInit,
		Headers: &h,
		Body:    "",
	}
}

var ERROR_BAD_STARTLINE = fmt.Errorf("bad request line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var ERROR_INCOMPLETE_STARTLINE = fmt.Errorf("incomplete request line")
var ERROR_PARSER_STATE_ERROR = fmt.Errorf("parser state error")
var SEPERATOR = []byte("\r\n")

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break outer
		}

		switch r.state {
		case StateError:
			return 0, ERROR_PARSER_STATE_ERROR
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}

			// not enough data yet, keep on parsing
			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.state = StateHeaders

		case StateHeaders:

			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}

			// keep parsing
			if n == 0 {
				break outer
			}

			read += n

			if done {
				if r.hasBody() {
					r.state = StateBody
				} else {
					r.state = StateDone
				}
			}

		case StateBody:
			length := getIntHeader(r.Headers, "content-length", 0)
			if length == 0 {
				panic("chunked not implemented")
			}

			remaining := min(length-len(r.Body), len(currentData))
			r.Body += string(currentData[:remaining])
			read += remaining

			// slog.Info("parseState#StateBody", "remaining", remaining, "body", r.Body, "read", read)

			if len(r.Body) == length {
				r.state = StateDone
			}

		case StateDone:
			break outer

		default:
			panic("somehow we have programmed poorly")
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPERATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, nil
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ERROR_BAD_STARTLINE
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	return rl, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	// NOTE: buffer could get overrun... a header that exceeds 1k would do that...
	// or the body
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		slog.Info("RequestFromHeader", "state", request.state, "bufLen", bufLen)
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

		copy(buf, buf[readN:readN+bufLen])
		bufLen -= readN

		slog.Info("RequestFromHeader", "newBufLen", bufLen, "n", n, "readN", readN)
	}

	return request, nil
}
