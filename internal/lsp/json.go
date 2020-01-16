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

// 将 Server/method 转换成 Server.Method
// 将 method 转换成 Server.Method
// 即未指定服务名称的，自动映身到 Server 这个服务名称
func convertMethod(method string) string {
	if len(method) == 0 {
		return method
	}

	m := []byte(method)
	dotIndex := bytes.IndexByte(m, '/')
	if dotIndex > 0 {
		m[dotIndex] = '.'
		m[dotIndex+1] = byte(unicode.ToUpper(rune(m[dotIndex+1])))
	} else if bytes.IndexByte(m, '.') < 0 {
		m = append([]byte("Server."), m...)
	}

	return string(m)
}
