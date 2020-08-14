// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/issue9/sliceutil"

	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

// workspace/workspaceFolders
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#workspace_workspaceFolders
func (s *server) workspaceWorkspaceFolders() error {
	return s.Send("workspace/workspaceFolders", nil, func(folders *[]protocol.WorkspaceFolder) error {
		s.workspaceMux.Lock()
		defer s.workspaceMux.Unlock()

		for _, f := range s.folders {
			f.close()
		}
		s.folders = s.folders[:0]

		s.appendFolders(*folders...)
		return nil
	})
}

// workspace/didChangeWorkspaceFolders
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#workspace_didChangeWorkspaceFolders
func (s *server) workspaceDidChangeWorkspaceFolders(notify bool, in *protocol.DidChangeWorkspaceFoldersParams, out *interface{}) error {
	if s.getState() != serverInitialized {
		return newError(ErrInvalidRequest, locale.ErrInvalidLSPState)
	}

	s.workspaceMux.Lock()
	defer s.workspaceMux.Unlock()

	for _, removed := range in.Event.Removed {
		sliceutil.Delete(s.folders, func(i int) bool {
			if f := s.folders[i]; f.Name == removed.Name && f.URI == removed.URI {
				f.close()
				return true
			}
			return false
		})
	}

	s.appendFolders(in.Event.Added...)

	return nil
}
