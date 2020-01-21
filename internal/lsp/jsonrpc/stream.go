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

func (s *stream) readRequest(req *Request) error {
	data, err := s.read()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, req)
}

func (s *stream) readResponse(resp *Response) error {
	data, err := s.read()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, resp)
}

// 读取内容，先验证报头，并返回实际 body 的内容
func (s *stream) read() ([]byte, error) {
	buf := bufio.NewReader(s.in)
	var l int

	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		index := strings.IndexByte(line, ':')
		if index <= 0 {
			return nil, locale.Errorf(locale.ErrInvalidHeaderFormat)
		}

		v := strings.TrimSpace(line[index+1:])
		switch http.CanonicalHeaderKey(strings.TrimSpace(line[:index])) {
		case contentLength:
			l, err = strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
		case contentType:
			if err := validContentType(v); err != nil {
				return nil, err
			}
		default: // 忽略其它报头
		}
	}

	if l <= 0 {
		return nil, locale.Errorf(locale.ErrInvalidContentLength)
	}

	data := make([]byte, l)
	n, err := io.ReadFull(buf, data)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, locale.Errorf(locale.ErrBodyIsEmpty)
	}

	return data[:n], nil
}

func validContentType(header string) error {
	pairs := strings.Split(header, ";")

	for _, pair := range pairs {
		index := strings.IndexByte(pair, '=')
		if index > 0 &&
			strings.ToLower(strings.TrimSpace(pair[:index])) == "charset" &&
			strings.ToLower(strings.TrimSpace(pair[index+1:])) != charset {
			return locale.Errorf(locale.ErrInvalidContentTypeCharset)
		}
	}

	return nil
}

func (s *stream) write(obj interface{}) (int, error) {
	data, err := json.Marshal(obj)
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
