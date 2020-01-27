// SPDX-License-Identifier: MIT

package jsonrpc

import (
	"encoding/json"
	"io"
	"net"
	"sync"
)

type tcpStream struct {
	in *json.Decoder

	out    io.Writer
	outMux sync.Mutex
}

// NewTCPStream 声明基于 TCP 通讯的 Streamer 实例
func NewTCPStream(conn *net.TCPConn) Streamer {
	return &tcpStream{
		in:  json.NewDecoder(conn),
		out: conn,
	}
}

func (s *tcpStream) Read(v interface{}) error {
	return s.in.Decode(v)
}

func (s *tcpStream) Write(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	s.outMux.Lock()
	defer s.outMux.Unlock()
	_, err = s.out.Write(data)
	return err
}
