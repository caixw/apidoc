// SPDX-License-Identifier: MIT

package lsp

import (
	"sync"

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
	state    serverState
	stateMux sync.RWMutex

	clientInfo *protocol.ServerInfo
}

func newServer() *server {
	return &server{
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
