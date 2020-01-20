// SPDX-License-Identifier: MIT

package jsonrpc

import "io"

// Conn 连接对象，json-rpc 客户端和服务端是对等的，两者都使用 conn 初始化。
type Conn struct {
	sequence int64
	stream   *stream
}

// NewConn 声明新的 Conn 实例
func NewConn(in io.Reader, out io.Writer) *Conn {
	return &Conn{
		stream: newStream(in, out),
	}
}
