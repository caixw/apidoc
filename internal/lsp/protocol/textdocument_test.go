// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/caixw/apidoc/v7/core"
	"github.com/issue9/assert"
)

func TestDidChangeTextDocumentParams_Blocks(t *testing.T) {
	a := assert.New(t)

	p := &DidChangeTextDocumentParams{}
	a.Empty(p.Blocks())

	p = &DidChangeTextDocumentParams{
		TextDocument: VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: TextDocumentIdentifier{URI: core.FileURI("test.go")},
		},
		ContentChanges: []TextDocumentContentChangeEvent{
			{
				Range: &core.Range{End: core.Position{Line: 1, Character: 5}},
				Text:  "text",
			},
		},
	}
	a.Equal(p.Blocks(), []core.Block{
		{
			Data: []byte("text"),
			Location: core.Location{
				URI:   core.FileURI("test.go"),
				Range: core.Range{End: core.Position{Line: 1, Character: 5}},
			},
		},
	})

	// 多个元素
	p = &DidChangeTextDocumentParams{
		TextDocument: VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: TextDocumentIdentifier{URI: core.FileURI("test.go")},
		},
		ContentChanges: []TextDocumentContentChangeEvent{
			{
				Range: &core.Range{End: core.Position{Line: 1, Character: 5}},
				Text:  "text",
			},

			{
				Range: &core.Range{End: core.Position{Line: 2, Character: 5}},
				Text:  "text2",
			},
		},
	}
	a.Equal(p.Blocks(), []core.Block{
		{
			Data: []byte("text"),
			Location: core.Location{
				URI:   core.FileURI("test.go"),
				Range: core.Range{End: core.Position{Line: 1, Character: 5}},
			},
		},
		{
			Data: []byte("text2"),
			Location: core.Location{
				URI:   core.FileURI("test.go"),
				Range: core.Range{End: core.Position{Line: 2, Character: 5}},
			},
		},
	})
}
