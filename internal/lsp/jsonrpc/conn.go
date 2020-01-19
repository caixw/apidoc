// SPDX-License-Identifier: MIT

package jsonrpc

// Conn 连接对象，json-rpc 客户端和服务端是对等的，两者都使用 conn 初始化。
type Conn struct {
	sequence int64
	stream   *stream
}
