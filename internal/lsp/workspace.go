// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

// workspace/workspaceFolders
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#workspace_workspaceFolders
func (s *server) workspaceWorkspaceFolders() error {
	err := s.Send("workspace/workspaceFolders", nil, func(folders *[]protocol.WorkspaceFolder) error {
		for _, f := range s.folders {
			f.close()
		}

		if len(*folders) != 0 {
			if err := s.appendFolders(*folders...); err != nil {
				return err
			}
		}
		return nil
	})
	return err

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
				f.close()
				s.folders = append(s.folders[:index], s.folders[index+1:]...)
			}
		}
	}

	return s.appendFolders(in.Event.Added...)
}
