// SPDX-License-Identifier: MIT

package jsonrpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/caixw/apidoc/v6/internal/locale"
)

// content-type json-rpc 采用的字符集
const charset = "utf-8"

var (
	contentType   = http.CanonicalHeaderKey("content-Type")
	contentLength = http.CanonicalHeaderKey("content-length")
)

type stream struct {
	in     io.Reader
	out    io.Writer
	outMux sync.Mutex
}

func newStream(in io.Reader, out io.Writer) *stream {
	return &stream{
		in:  in,
		out: out,
	}
}

func (s *stream) read(req *Request) error {
	buf := bufio.NewReader(s.in)
	var l int

	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			return NewError(CodeParseError, err.Error())
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		index := strings.IndexByte(line, ':')
		if index <= 0 {
			return NewError(CodeParseError, locale.Sprintf(locale.ErrInvalidFormat))
		}

		v := strings.TrimSpace(line[index+1:])
		switch http.CanonicalHeaderKey(strings.TrimSpace(line[:index])) {
		case contentLength:
			l, err = strconv.Atoi(v)
			if err != nil {
				return NewError(CodeParseError, err.Error())
			}
		case contentType:
			if err := validContentType(v); err != nil {
				return err
			}
		default: // 忽略其它报头
		}
	}

	if l == 0 {
		return NewError(CodeParseError, locale.Sprintf(locale.ErrRequired))
	}

	data := make([]byte, l)
	n, err := io.ReadFull(buf, data)
	if err != nil {
		return NewError(CodeParseError, err.Error())
	}
	if n == 0 {
		return NewError(CodeParseError, locale.Sprintf(locale.ErrRequired))
	}

	data = data[:n]
	return json.Unmarshal(data, req)
}

func validContentType(header string) error {
	pairs := strings.Split(header, ";")

	for _, pair := range pairs {
		index := strings.IndexByte(pair, '=')
		if index > 0 &&
			strings.ToLower(strings.TrimSpace(pair[:index])) == "charset" &&
			strings.ToLower(strings.TrimSpace(pair[index+1:])) != charset {
			return NewError(CodeParseError, "TODO locale")
		}
	}

	return nil
}

func (s *stream) write(resp *Response) (int, error) {
	data, err := json.Marshal(resp)
	if err != nil {
		return 0, err
	}

	s.outMux.Lock()
	defer s.outMux.Unlock()

	n, err := fmt.Fprintf(s.out, "%s: %s\r\n%s: %d\r\n\r\n", contentType, charset, contentLength, len(data))
	if err != nil {
		return 0, err
	}

	size, err := s.out.Write(data)
	return n + size, err
}
