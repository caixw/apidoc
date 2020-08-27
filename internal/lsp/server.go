// SPDX-License-Identifier: MIT

package lsp

import (
	"context"
	"log"
	"strings"
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
	trace        string
	stateMux     sync.RWMutex
	workspaceMux sync.RWMutex

	folders []*folder

	clientParams *protocol.InitializeParams
	serverResult *protocol.InitializeResult
	info, erro   *log.Logger
	cancelFunc   context.CancelFunc
}

func newServe(t jsonrpc.Transport, infolog, errlog *log.Logger) *server {
	jsonrpcServer := jsonrpc.NewServer()

	srv := &server{
		Conn:  jsonrpcServer.NewConn(t, errlog),
		state: serverCreated,
		trace: protocol.TraceValueOff,
		info:  infolog,
		erro:  errlog,
	}

	jsonrpcServer.Registers(map[string]interface{}{
		"initialize":      srv.initialize,
		"initialized":     srv.initialized,
		"shutdown":        srv.shutdown,
		"exit":            srv.exit,
		"$/cancelRequest": srv.cancel,
		"$/setTrace":      srv.setTrace,

		// workspace
		"workspace/didChangeWorkspaceFolders": srv.workspaceDidChangeWorkspaceFolders,

		// textDocument
		"textDocument/didChange":      srv.textDocumentDidChange,
		"textDocument/hover":          srv.textDocumentHover,
		"textDocument/foldingRange":   srv.textDocumentFoldingRange,
		"textDocument/completion":     srv.textDocumentCompletion,
		"textDocument/semanticTokens": srv.textDocumentSemanticTokens,
		"textDocument/references":     srv.textDocumentReferences,
		"textDocument/definition":     srv.textDocumentDefinition,

		// apidoc 自定义的接口
		"apidoc/refreshOutline": srv.apidocRefreshOutline,
	})

	jsonrpcServer.RegisterMatcher(func(method string) bool {
		return strings.HasPrefix(method, "$/")
	}, srv.dollarHandler)

	return srv
}

func (s *server) serve() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	return s.Serve(ctx)
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

// $/setTrace
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#setTrace
func (s *server) setTrace(notify bool, in *protocol.SetTraceParams, out *interface{}) error {
	if protocol.IsValidTraceValue(in.Value) {
		s.trace = in.Value
	}
	s.trace = protocol.TraceValueOff
	return nil
}

// $/logTrace
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#logTrace
func (s *server) logTrace(message, verbose string) {
	if p := protocol.BuildLogTrace(s.trace, message, verbose); p != nil {
		if err := s.Notify("$/logTrace", p); err != nil {
			s.erro.Println(err)
		}
	}
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

func (s *server) printErr(err error) {
	s.erro.Println(err)
	s.logTrace(err.Error(), "")
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
