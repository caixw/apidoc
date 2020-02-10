// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/lsp/protocol"
)

// The workspace/workspaceFolders request is sent from the server to the client to fetch the current open
// list of workspace folders. Returns null in the response if only a single file is open in the tool.
// Returns an empty array if a workspace is open but no folders are configured.
func (s *server) workspaceWorkspaceFolders() error {
	var folders []protocol.WorkspaceFolder
	if err := s.Send("workspace/workspaceFolders", nil, &folders); err != nil {
		return err
	}

	if len(folders) != 0 {
		s.workspaceFolders = folders
	}
	return nil
}

// workspace/didChangeWorkspaceFolders
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#workspace_didChangeWorkspaceFolders
func (s *server) workspaceDidChangeWorkspaceFolders(notify bool, in *protocol.DidChangeWorkspaceFoldersParams, out *interface{}) error {
	if s.getState() != serverInitialized {
		return newError(ErrInvalidRequest, locale.ErrInvalidLSPState)
	}

	for _, folder := range in.Event.Removed {
		for index, f2 := range s.workspaceFolders {
			if f2.Name == folder.Name && f2.URI == folder.URI {
				s.workspaceFolders = append(s.workspaceFolders[:index], s.workspaceFolders[index+1:]...)
			}
		}
	}

	s.workspaceFolders = append(s.workspaceFolders, in.Event.Added...)
	return nil
}
