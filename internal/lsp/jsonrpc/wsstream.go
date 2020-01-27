// SPDX-License-Identifier: MIT

package jsonrpc

import (
	"github.com/gorilla/websocket"
)

type websocketStream struct {
	conn *websocket.Conn
}

// NewWebsocketStream 声明基于 websocket 的 streamer 实例
func NewWebsocketStream(conn *websocket.Conn) Streamer {
	return &websocketStream{conn: conn}
}

func (s *websocketStream) Read(v interface{}) error {
	return s.conn.ReadJSON(v)
}

func (s *websocketStream) Write(v interface{}) error {
	return s.conn.WriteJSON(v)
}
