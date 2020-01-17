// SPDX-License-Identifier: MIT

package lsp

import (
	"net/http"
	"sync"

	"github.com/gorilla/rpc/v2"

	"github.com/caixw/apidoc/v6/internal/locale"
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
	stateMux sync.RWMutex
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
	if s.getState() > serverInitializing {
		msg := locale.Sprintf(locale.ErrServerNotInitialized)
		return protocol.NewError(protocol.ErrServerNotInitialized, msg, nil)
	}

	// TODO notify

	s.setState(serverInitializing)

	// TODO
	return nil
}

func (s *Server) setState(state serverState) {
	s.stateMux.Lock()
	defer s.stateMux.Unlock()
	s.state = state
}

func (s *Server) getState() serverState {
	s.stateMux.RLock()
	defer s.stateMux.RUnlock()
	return s.state
}
