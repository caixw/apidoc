// SPDX-License-Identifier: MIT

package lsp

import (
	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

// initialize
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialize
func (s *server) initialize(notify bool, in *protocol.InitializeParams, out *protocol.InitializeResult) error {
	if s.getState() != serverCreated {
		return newError(ErrServerNotInitialized, locale.ErrInvalidLSPState)
	}

	s.setState(serverInitializing)

	s.clientInfo = in.ClientInfo
	s.clientCapabilities = &in.Capabilities

	if in.Capabilities.Workspace.WorkspaceFolders {
		out.Capabilities.Workspace.WorkspaceFolders.Supported = true
		out.Capabilities.Workspace.WorkspaceFolders.ChangeNotifications = true
	}

	out.Capabilities.TextDocumentSync = &protocol.ServerCapabilitiesTextDocumentSyncOptions{
		OpenClose: true,
		Change:    protocol.TextDocumentSyncKindFull,
	}

	if in.Capabilities.TextDocument.Hover.ContentFormat != nil {
		out.Capabilities.HoverProvider = true
	}

	if in.InitializationOptions != nil && in.InitializationOptions.Locale != "" {
		tag, err := language.Parse(in.InitializationOptions.Locale)
		if err != nil {
			s.erro.Println(err) // 输出错误信息，但是不中断执行
		}
		locale.SetTag(tag)
	}

	out.ServerInfo = &protocol.ServerInfo{
		Name:    core.Name,
		Version: core.Version,
	}

	return s.appendFolders(in.Folders()...)
}

// initialized
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialized
func (s *server) initialized(bool, *protocol.InitializedParams, *interface{}) error {
	if s.getState() != serverInitializing {
		return newError(ErrInvalidRequest, locale.ErrInvalidLSPState)
	}
	s.setState(serverInitialized)

	if s.clientCapabilities.Workspace.WorkspaceFolders {
		if err := s.workspaceWorkspaceFolders(); err != nil {
			return err
		}
	}

	return nil
}

// shutdown
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#shutdown
func (s *server) shutdown(bool, *interface{}, *interface{}) error {
	if s.getState() != serverInitialized {
		return newError(ErrInvalidRequest, locale.ErrInvalidLSPState)
	}

	for _, f := range s.folders {
		f.close()
	}

	if s.cancelFunc != nil {
		s.cancelFunc()
	}

	s.setState(serverShutdown)

	return nil
}

// exit
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#exit
func (s *server) exit(bool, *interface{}, *interface{}) error {
	if s.getState() != serverShutdown {
		return newError(ErrInvalidRequest, locale.ErrInvalidLSPState)
	}

	return nil
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
