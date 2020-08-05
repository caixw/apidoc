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

	s.clientParams = in

	if in.Capabilities.Workspace.WorkspaceFolders {
		out.Capabilities.Workspace = &protocol.WorkspaceProvider{
			WorkspaceFolders: &protocol.WorkspaceFoldersServerCapabilities{
				Supported:           true,
				ChangeNotifications: true,
			},
		}
	}

	out.Capabilities.TextDocumentSync = &protocol.ServerCapabilitiesTextDocumentSyncOptions{
		Change: protocol.TextDocumentSyncKindFull,
	}

	if in.Capabilities.TextDocument.Hover.ContentFormat != nil {
		out.Capabilities.HoverProvider = true
	}

	if in.Capabilities.TextDocument.FoldingRange != nil {
		out.Capabilities.FoldingRangeProvider = true
	}

	if in.Capabilities.TextDocument.Completion != nil {
		out.Capabilities.CompletionProvider = &protocol.CompletionOptions{}
	}

	if in.Capabilities.TextDocument.SemanticTokens != nil {
		out.Capabilities.SemanticTokensProvider = &protocol.SemanticTokensOptions{
			Legend: protocol.SemanticTokensLegend{
				TokenTypes:     []string{"type", "property", "variable"},
				TokenModifiers: []string{"documentation"},
			},
			Range: true,
			Full:  true,
		}
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

	s.workspaceMux.Lock()
	defer s.workspaceMux.Unlock()

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

	if s.clientParams.Capabilities.Workspace.WorkspaceFolders {
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
