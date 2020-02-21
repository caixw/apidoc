// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/lsp/protocol"
)

// workspace/workspaceFolders
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#workspace_workspaceFolders
func (s *server) workspaceWorkspaceFolders() error {
	var folders []protocol.WorkspaceFolder
	if err := s.Send("workspace/workspaceFolders", nil, &folders); err != nil {
		return err
	}

	for _, f := range s.folders {
		if err := f.close(); err != nil {
			return err
		}
	}

	if len(folders) != 0 {
		s.appendFolders(folders...)
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

	for _, removed := range in.Event.Removed {
		for index, f := range s.folders {
			if f.Name == removed.Name && f.URI == removed.URI {
				if err := f.close(); err != nil {
					return err
				}
				s.folders = append(s.folders[:index], s.folders[index+1:]...)
			}
		}
	}

	s.appendFolders(in.Event.Added...)
	return nil
}
