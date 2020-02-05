// SPDX-License-Identifier: MIT

package lsp

import (
	"sync"

	"github.com/issue9/jsonrpc"

	"github.com/caixw/apidoc/v6/internal/lsp/protocol"
)

type serverState int

const (
	serverCreated serverState = iota
	serverInitializing
	serverInitialized
	serverShutDown
)

// server LSP 服务实例
type server struct {
	*jsonrpc.Conn
	state    serverState
	stateMux sync.RWMutex

	workspaceFolders []protocol.WorkspaceFolder

	clientInfo *protocol.ServerInfo
}

func newServer(conn *jsonrpc.Conn) *server {
	return &server{
		Conn:  conn,
		state: serverCreated,
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
