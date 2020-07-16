// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/token"
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

	// 未指定 Range
	p = &DidChangeTextDocumentParams{
		TextDocument: VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: TextDocumentIdentifier{URI: core.FileURI("test.go")},
		},
		ContentChanges: []TextDocumentContentChangeEvent{
			{
				Text: "text",
			},
		},
	}
	a.Equal(p.Blocks(), []core.Block{
		{
			Data: []byte("text"),
			Location: core.Location{
				URI: core.FileURI("test.go"),
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

func TestBuildFoldingRange(t *testing.T) {
	a := assert.New(t)

	base := token.Base{}
	a.Equal(BuildFoldingRange(base, false), FoldingRange{Kind: FoldingRangeKindComment})

	base = token.Base{Range: core.Range{
		Start: core.Position{Line: 1, Character: 11},
		End:   core.Position{Line: 2, Character: 11},
	}}
	a.Equal(BuildFoldingRange(base, false), FoldingRange{
		StartLine: 1,
		Kind:      FoldingRangeKindComment,
		EndLine:   2,
	})
	a.Equal(BuildFoldingRange(base, true), FoldingRange{
		StartLine:      1,
		StartCharacter: &base.Start.Character,
		EndLine:        2,
		EndCharacter:   &base.End.Character,
		Kind:           FoldingRangeKindComment,
	})
}
