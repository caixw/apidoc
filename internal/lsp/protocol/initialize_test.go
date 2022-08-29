// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/issue9/assert/v3"
)

func TestInitializeParams_Folders(t *testing.T) {
	a := assert.New(t, false)

	p := &InitializeParams{}
	a.Nil(p.Folders())

	p.RootPath = "../../lsp/protocol"
	a.Equal(p.Folders(), []WorkspaceFolder{
		{
			Name: "protocol",
			URI:  "../../lsp/protocol",
		},
	})

	p.RootPath = "/../lsp/protocol2"
	a.Equal(p.Folders(), []WorkspaceFolder{
		{
			Name: "protocol2",
			URI:  "/../lsp/protocol2",
		},
	})

	p.RootURI = "file:///lsp/protocol"
	a.Equal(p.Folders(), []WorkspaceFolder{
		{
			Name: "protocol",
			URI:  "file:///lsp/protocol",
		},
	})

	p.WorkspaceFolders = []WorkspaceFolder{
		{
			Name: "f1",
			URI:  "file:///f1",
		},
		{
			Name: "f2",
			URI:  "https://example.com/f1",
		},
	}
	a.Equal(p.Folders(), []WorkspaceFolder{
		{
			Name: "f1",
			URI:  "file:///f1",
		},
		{
			Name: "f2",
			URI:  "https://example.com/f1",
		},
	})
}
