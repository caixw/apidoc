// SPDX-License-Identifier: MIT

package lsp

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/issue9/assert/v2"
	"github.com/issue9/jsonrpc"

	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

func TestServer_workspaceDidChangeWorkspaceFolders(t *testing.T) {
	a := assert.New(t, false)

	s := newTestServer(true, log.New(ioutil.Discard, "", 0), log.New(ioutil.Discard, "", 0))
	in := &protocol.DidChangeWorkspaceFoldersParams{}
	err := s.workspaceDidChangeWorkspaceFolders(false, in, nil)
	a.Error(err)
	jerr, ok := err.(*jsonrpc.Error)
	a.True(ok).Equal(jerr.Code, ErrInvalidRequest)

	s = newTestServer(true, log.New(ioutil.Discard, "", 0), log.New(ioutil.Discard, "", 0))
	s.setState(serverInitialized)
	in = &protocol.DidChangeWorkspaceFoldersParams{}
	a.NotError(s.workspaceDidChangeWorkspaceFolders(false, in, nil))
	a.Equal(0, len(s.folders))

	s = newTestServer(true, log.New(ioutil.Discard, "", 0), log.New(ioutil.Discard, "", 0))
	s.setState(serverInitialized)
	s.folders = []*folder{
		{
			srv:             s,
			WorkspaceFolder: protocol.WorkspaceFolder{Name: "n0", URI: "file:///n0"},
			h:               messagetest.NewMessageHandler().Handler,
		},
		{
			srv:             s,
			WorkspaceFolder: protocol.WorkspaceFolder{Name: "n1", URI: "file:///n1"},
			h:               messagetest.NewMessageHandler().Handler,
		},
	}

	in = &protocol.DidChangeWorkspaceFoldersParams{
		Event: protocol.WorkspaceFoldersChangeEvent{
			Added: []protocol.WorkspaceFolder{
				{Name: "n3", URI: "file:///n3"},
			},
			Removed: []protocol.WorkspaceFolder{
				{Name: "n1", URI: "file:///n1"},
			},
		},
	}
	a.NotError(s.workspaceDidChangeWorkspaceFolders(false, in, nil))
	a.Equal(2, len(s.folders))
}
