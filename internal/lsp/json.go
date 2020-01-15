// SPDX-License-Identifier: MIT

package lsp

import (
	"bytes"
	"net/http"
	"unicode"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
)

type jsonCodec struct {
	*json2.Codec
}

func newJSONCodec() *jsonCodec {
	return &jsonCodec{
		Codec: json2.NewCodec(),
	}
}

// NewRequest 返回新的 rpc.CodecRequest 对象
func (c *jsonCodec) NewRequest(r *http.Request) rpc.CodecRequest {
	return &jsonCodecRequest{
		CodecRequest: c.Codec.NewRequest(r).(*json2.CodecRequest),
	}
}

type jsonCodecRequest struct {
	*json2.CodecRequest
}

// Method 返回 LSP 请求的 method 字段
func (r *jsonCodecRequest) Method() (string, error) {
	m, err := r.CodecRequest.Method()
	return convertMethod(m), err
}

func convertMethod(method string) string {
	m := []byte(method)
	dotIndex := bytes.IndexByte(m, '/')
	if dotIndex > 0 {
		m[dotIndex] = '.'
		m[dotIndex+1] = byte(unicode.ToUpper(rune(m[dotIndex+1])))
	}
	return string(m)
}
