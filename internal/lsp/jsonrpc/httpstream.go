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

type httpStream struct {
	in     io.Reader
	out    io.Writer
	outMux sync.Mutex
}

// NewHTTPStream 声明基于 HTTP 的 streamer 实例
func NewHTTPStream(in io.Reader, out io.Writer) Streamer {
	return &httpStream{
		in:  in,
		out: out,
	}
}

// Read 读取内容，先验证报头，并返回实际 body 的内容
func (s *httpStream) Read(v interface{}) error {
	buf := bufio.NewReader(s.in)
	var l int

	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		index := strings.IndexByte(line, ':')
		if index <= 0 {
			return locale.Errorf(locale.ErrInvalidHeaderFormat)
		}

		v := strings.TrimSpace(line[index+1:])
		switch http.CanonicalHeaderKey(strings.TrimSpace(line[:index])) {
		case contentLength:
			l, err = strconv.Atoi(v)
			if err != nil {
				return err
			}
		case contentType:
			if err := validContentType(v); err != nil {
				return err
			}
		default: // 忽略其它报头
		}
	}

	if l <= 0 {
		return locale.Errorf(locale.ErrInvalidContentLength)
	}

	data := make([]byte, l)
	n, err := io.ReadFull(buf, data)
	if err != nil {
		return err
	}
	if n == 0 {
		return locale.Errorf(locale.ErrBodyIsEmpty)
	}

	return json.Unmarshal(data[:n], v)
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

func (s *httpStream) Write(obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	s.outMux.Lock()
	defer s.outMux.Unlock()

	_, err = fmt.Fprintf(s.out, "%s: %s\r\n%s: %d\r\n\r\n", contentType, charset, contentLength, len(data))
	if err != nil {
		return err
	}

	_, err = s.out.Write(data)
	return err
}
