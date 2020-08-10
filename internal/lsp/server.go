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

// window/showMessage
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#window_showMessage
func (s *server) windowShowMessage(t protocol.MessageType, messag string) error {
	return s.Notify("window/showMessage", &protocol.ShowMessageParams{
		Type:    t,
		Message: messag,
	})
}

func (s *server) windowShowInfoMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeInfo, locale.Sprintf(key, v...))
}

func (s *server) windowShowLogMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeLog, locale.Sprintf(key, v...))
}

func (s *server) windowShowWarnMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeWarning, locale.Sprintf(key, v...))
}

func (s *server) windowShowErrorMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeError, locale.Sprintf(key, v...))
}

// window/showMessageRequest
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#window_showMessageRequest
func (s *server) windowShowMessageRequest(t protocol.MessageType, actions []protocol.MessageActionItem, message string) (*protocol.MessageActionItem, error) {
	out := &protocol.MessageActionItem{}
	in := &protocol.ShowMessageRequestParams{
		Type:    t,
		Message: message,
		Actions: actions,
	}

	err := s.Send("window/showMessageRequest", in, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *server) windowShowInfoMessageRequest(actions []protocol.MessageActionItem, key message.Reference, v ...interface{}) (*protocol.MessageActionItem, error) {
	return s.windowShowMessageRequest(protocol.MessageTypeInfo, actions, locale.Sprintf(key, v...))
}

func (s *server) windowShowLogMessageRequest(actions []protocol.MessageActionItem, key message.Reference, v ...interface{}) (*protocol.MessageActionItem, error) {
	return s.windowShowMessageRequest(protocol.MessageTypeLog, actions, locale.Sprintf(key, v...))
}

func (s *server) windowShowWarnMessageRequest(actions []protocol.MessageActionItem, key message.Reference, v ...interface{}) (*protocol.MessageActionItem, error) {
	return s.windowShowMessageRequest(protocol.MessageTypeWarning, actions, locale.Sprintf(key, v...))
}

func (s *server) windowShowErrorMessageRequest(actions []protocol.MessageActionItem, key message.Reference, v ...interface{}) (*protocol.MessageActionItem, error) {
	return s.windowShowMessageRequest(protocol.MessageTypeError, actions, locale.Sprintf(key, v...))
}

// window/logMessage
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#window_logMessage
func (s *server) windowLogMessage(t protocol.MessageType, message string) error {
	return s.Notify("window/logMessage", &protocol.LogMessageParams{
		Type:    t,
		Message: message,
	})
}

func (s *server) windowLogInfoMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeInfo, locale.Sprintf(key, v...))
}

func (s *server) windowLogLogMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeLog, locale.Sprintf(key, v...))
}

func (s *server) windowLogWarnMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeWarning, locale.Sprintf(key, v...))
}

func (s *server) windowLogErrorMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeError, locale.Sprintf(key, v...))
}
