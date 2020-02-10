// SPDX-License-Identifier: MIT

package lsp

import (
	"context"
	"sync"

	"github.com/issue9/jsonrpc"

	"github.com/caixw/apidoc/v6/internal/lsp/protocol"
)

type serverState int

const (
	serverCreated serverState = iota
	serverInitializing
	serverInitialized
	serverShutdown
)

// server LSP 服务实例
type server struct {
	*jsonrpc.Conn

	state    serverState
	stateMux sync.RWMutex

	cancelFunc context.CancelFunc

	workspaceFolders []protocol.WorkspaceFolder

	clientInfo         *protocol.ServerInfo
	clientCapabilities *protocol.ClientCapabilities
}

func newServer(conn *jsonrpc.Conn, cancel context.CancelFunc) *server {
	return &server{
		Conn:       conn,
		state:      serverCreated,
		cancelFunc: cancel,
	}
}

func (s *server) setState(state serverState) {
	s.stateMux.Lock()
	defer s.stateMux.Unlock()
	s.state = state
}

func (s *server) getState() serverState {
	s.stateMux.RLock()
	defer s.stateMux.RUnlock()
	return s.state
}

func (s *server) close() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}
