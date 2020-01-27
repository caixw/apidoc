// SPDX-License-Identifier: MIT

package jsonrpc

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

var _ Streamer = &httpStream{}

func TestHTTPStream_Read(t *testing.T) {
	a := assert.New(t)

	r := new(bytes.Buffer)
	s := NewHTTPStream(r, nil)
	a.NotNil(s)
	r.WriteString(`Content-Type: text/json;charset=utf-8
	Content-Length:26

{"jsonrpc":"2.0","id":"1"}`)
	rr := &Request{}
	a.NotError(s.Read(rr))
	a.Equal(rr.Version, Version).Equal(rr.ID, "1")

	// 无效的 content-length
	r = new(bytes.Buffer)
	s = NewHTTPStream(r, nil)
	a.NotNil(s)
	r.WriteString(`Content-Type: text/json;charset=utf-8
	Content-Length:0

{"jsonrpc":"2.0","id":"1"}`)
	rr = &Request{}
	a.Error(s.Read(rr))

	// content-type 中未指定 charset
	r = new(bytes.Buffer)
	s = NewHTTPStream(r, nil)
	a.NotNil(s)
	r.WriteString(`Content-Type: text/json;charset-xx=utf-8
	Content-Length:26

{"jsonrpc":"2.0","id":"1"}`)
	rr = &Request{}
	a.NotError(s.Read(rr))

	// content-length 格式无效
	r = new(bytes.Buffer)
	s = NewHTTPStream(r, nil)
	a.NotNil(s)
	r.WriteString(`Content-Type: text/json;charset-xx=utf-8
	Content-Length:26xx

{"jsonrpc":"2.0","id":"1"}`)
	rr = &Request{}
	a.Error(s.Read(rr))

	// content-type 是指定了非 utf-8 编码
	r = new(bytes.Buffer)
	s = NewHTTPStream(r, nil)
	a.NotNil(s)
	r.WriteString(`Content-Type: text/json;charset-xx=utf-7
	Content-Length:26xx

{"jsonrpc":"2.0","id":"1"}`)
	rr = &Request{}
	a.Error(s.Read(rr))
}

func TestHTTPStream_Write(t *testing.T) {
	a := assert.New(t)
	w := new(bytes.Buffer)
	s := NewHTTPStream(nil, w)
	a.NotNil(s)

	err := s.Write(&Response{
		Version: "1.0.1",
		Error: &Error{
			Code:    CodeParseError,
			Message: "message",
		},
		ID: "1",
	})
	a.NotError(err)
	a.NotEmpty(w.Bytes())
}

func TestValidContentType(t *testing.T) {
	a := assert.New(t)

	a.NotError(validContentType("text/xml"))
	a.NotError(validContentType(""))
	a.NotError(validContentType("charset=utf-8"))
	a.NotError(validContentType(";charset=utf-8"))
	a.NotError(validContentType("text/xml;charset=utf-8"))
	a.NotError(validContentType("text/xml;"))
	a.Error(validContentType("text/xml;charset="))
	a.Error(validContentType("text/xml;charset=utf8"))
}
