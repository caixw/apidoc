// SPDX-License-Identifier: MIT

package lsp

import "github.com/caixw/apidoc/v6/internal/lsp/protocol"

type workspace struct {
	server *server
}

func newWorkspace(srv *server) *workspace {
	return &workspace{server: srv}
}

// The workspace/workspaceFolders request is sent from the server to the client to fetch the current open
// list of workspace folders. Returns null in the response if only a single file is open in the tool.
// Returns an empty array if a workspace is open but no folders are configured.
func (w *workspace) workspaceFolders(notify bool, in interface{}, out *[]protocol.WorkspaceFolder) error {
	// TODO
	return nil
}
