// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/lsp/protocol"
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
	}

	return nil
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
		if err := f.close(); err != nil {
			return err
		}
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
