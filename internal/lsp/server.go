// SPDX-License-Identifier: MIT

package lsp

import (
	"net/http"
	"sync"

	"github.com/gorilla/rpc/v2"

	"github.com/caixw/apidoc/v6/internal/lsp/protocol"
)

type serverState int

const (
	serverCreated serverState = iota
	serverInitializing
	serverInitialized
	serverShutDown
)

// Server LSP 服务实例
type Server struct {
	server *rpc.Server

	state    serverState
	stateMux sync.Mutex
}

// NewServer 新的 Server 实例
func NewServer() (*Server, error) {
	srv := &Server{
		server: rpc.NewServer(),
		state:  serverCreated,
	}

	srv.server.RegisterCodec(newJSONCodec(), "application/json")

	if err := srv.server.RegisterService(hello, "hello"); err != nil {
		return nil, err
	}

	return srv, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
}

// Initialize 执行初始化服务
func (s *Server) Initialize(r *http.Request, args *protocol.InitializeParams, reply *protocol.InitializeResult) error {
	s.stateMux.Lock()
	defer s.stateMux.Unlock()
	s.state = serverInitializing

	// TODO
	return nil
}
