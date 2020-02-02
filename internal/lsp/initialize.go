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
	}

	return nil
}

func (s *server) initialized(notify bool, in *protocol.InitializedParams, out interface{}) error {
	// TODO
	return nil
}
