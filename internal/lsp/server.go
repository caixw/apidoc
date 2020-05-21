// SPDX-License-Identifier: MIT

package lsp

import (
	"context"
	"log"
	"sync"

	"github.com/issue9/jsonrpc"

	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
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

	folders []*folder

	clientInfo         *protocol.ServerInfo
	clientCapabilities *protocol.ClientCapabilities

	info, erro *log.Logger
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
