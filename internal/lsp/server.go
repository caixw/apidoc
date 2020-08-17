// SPDX-License-Identifier: MIT

package lsp

import (
	"context"
	"log"
	"sync"

	"github.com/issue9/jsonrpc"
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/internal/locale"
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
	state        serverState
	stateMux     sync.RWMutex
	workspaceMux sync.RWMutex

	folders []*folder

	clientParams *protocol.InitializeParams
	serverResult *protocol.InitializeResult
	info, erro   *log.Logger
	cancelFunc   context.CancelFunc
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

// $/cancelRequest
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#cancelRequest
func (s *server) cancel(notify bool, in *protocol.CancelParams, out *interface{}) error {
	return nil
}

// 所有以 $/ 开头且未处理的服务由此函数处理
//
// $ Notifications and Requests
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#dollarRequests
func (s *server) dollarHandler(notify bool, in, out *interface{}) error {
	if !notify {
		return newError(ErrMethodNotFound, locale.UnimplementedRPC, "$/***")
	}
	return nil
}

// window/logMessage
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#window_logMessage
func (s *server) windowLogMessage(t protocol.MessageType, message string) {
	err := s.Notify("window/logMessage", &protocol.LogMessageParams{
		Type:    t,
		Message: message,
	})
	if err != nil {
		s.erro.Println(err)
	}
}

func (s *server) windowLogInfoMessage(key message.Reference, v ...interface{}) {
	s.windowLogMessage(protocol.MessageTypeInfo, locale.Sprintf(key, v...))
}

func (s *server) windowLogLogMessage(key message.Reference, v ...interface{}) {
	s.windowLogMessage(protocol.MessageTypeLog, locale.Sprintf(key, v...))
}

func (s *server) windowLogWarnMessage(key message.Reference, v ...interface{}) {
	s.windowLogMessage(protocol.MessageTypeWarning, locale.Sprintf(key, v...))
}

func (s *server) windowLogErrorMessage(key message.Reference, v ...interface{}) {
	s.windowLogMessage(protocol.MessageTypeError, locale.Sprintf(key, v...))
}
