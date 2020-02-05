// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/issue9/jsonrpc"

	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/lsp/protocol"
)

func (s *server) initialize(notify bool, in *protocol.InitializeParams, out *protocol.InitializeResult) error {
	if s.getState() > serverInitializing {
		msg := locale.Sprintf(locale.ErrServerNotInitialized)
		return jsonrpc.NewError(ErrServerNotInitialized, msg)
	}

	s.setState(serverInitializing)

	s.clientInfo = in.ClientInfo

	if in.Capabilities.Workspace.WorkspaceFolders {
		out.Capabilities.Workspace.WorkspaceFolders.Supported = true
		if err := s.workspaceWorkspaceFolders(); err != nil {
			return err
		}
	}

	return nil
}

// initialized
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialized
func (s *server) initialized(bool, *protocol.InitializedParams, interface{}) error {
	s.setState(serverInitialized)
	return nil
}

// shutdown
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#shutdown
func (s *server) shutdown(bool, interface{}, interface{}) error {
	s.setState(serverShutDown)
	return nil
}

func (s *server) exit(bool, interface{}, interface{}) error {
	// TODO
	return nil
}
