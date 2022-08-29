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
func (s *server) workspaceDidChangeWorkspaceFolders(notify bool, in *protocol.DidChangeWorkspaceFoldersParams, out *any) error {
	if s.getState() != serverInitialized {
		return newError(ErrInvalidRequest, locale.ErrInvalidLSPState)
	}

	s.workspaceMux.Lock()
	defer s.workspaceMux.Unlock()

	old := s.folders
	s.folders = sliceutil.Delete(s.folders, func(i *folder) bool {
		return sliceutil.Count(in.Event.Removed, func(j protocol.WorkspaceFolder) bool {
			return j.URI == i.URI
		}) > 0
	})

	deleted := old[len(s.folders):]
	for _, f := range deleted {
		f.close()
	}

	s.appendFolders(in.Event.Added...)

	return nil
}
