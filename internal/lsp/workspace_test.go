// SPDX-License-Identifier: MIT

package lsp

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/jsonrpc"

	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

func TestServer_workspaceDidChangeWorkspaceFolders(t *testing.T) {
	a := assert.New(t)

	s := &server{}
	in := &protocol.DidChangeWorkspaceFoldersParams{}
	err := s.workspaceDidChangeWorkspaceFolders(false, in, nil)
	a.Error(err)
	jerr, ok := err.(*jsonrpc.Error)
	a.True(ok).Equal(jerr.Code, ErrInvalidRequest)

	s = &server{
		Conn:  jsonrpc.NewServer().NewConn(jsonrpc.NewStreamTransport(true, os.Stdin, os.Stdout, nil), nil),
		state: serverInitialized,
	}
	in = &protocol.DidChangeWorkspaceFoldersParams{}
	a.NotError(s.workspaceDidChangeWorkspaceFolders(false, in, nil))
	a.Equal(0, len(s.folders))

	s = &server{
		Conn:  jsonrpc.NewServer().NewConn(jsonrpc.NewStreamTransport(true, os.Stdin, os.Stdout, nil), nil),
		state: serverInitialized,
		erro:  log.New(ioutil.Discard, "", 0),
	}
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
