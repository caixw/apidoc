// SPDX-License-Identifier: MIT

package jsonrpc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

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

func (s *stream) read() ([]byte, int, error) {
	data, err := ioutil.ReadAll(s.in)
	if err != nil {
		return nil, 0, err
	}

	index := bytes.IndexByte(data, '{')

	headers := data[:index]
	// TODO 验证报头正确性

	return data[index:], len(data), nil

}

func (s *stream) write(data []byte) (int, error) {
	s.outMux.Lock()
	defer s.outMux.Unlock()

	n, err := fmt.Fprintf(s.out, "%s: %s\r\n%s: %d\r\n\r\n", contentType, "utf-8", contentLength, len(data))
	if err != nil {
		return 0, err
	}

	size, err := s.out.Write(data)
	return n + size, err
}
