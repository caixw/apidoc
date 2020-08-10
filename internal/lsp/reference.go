// SPDX-License-Identifier: MIT

package lsp

import (
	"reflect"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

var (
	referencerType   = reflect.TypeOf((*ast.Referencer)(nil)).Elem()
	definitionerType = reflect.TypeOf((*ast.Definitioner)(nil)).Elem()
)

// textDocument/references
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_references
func (s *server) textDocumentReferences(notify bool, in *protocol.ReferenceParams, out *[]core.Location) error {
	f := s.findFolder(in.TextDocument.URI)
	if f == nil {
		return nil
	}

	f.parsedMux.RLock()
	defer f.parsedMux.RUnlock()

	*out = references(f.doc, in.TextDocument.URI, in.Position, in.Context.IncludeDeclaration)
	return nil
}

// textDocument/definition
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_definition
func (s *server) textDocumentDefinition(notify bool, in *protocol.DefinitionParams, out *[]core.Location) error {
	// NOTE: LSP 允许 out 的值是 null，而 jsonrpc 模块默认情况下是空值，而不是 nil，
	// 所以在可能的情况下，都尽量将其返回类型改为数组，
	// 或是像 protocol.Hover 一样为返回类型实现 json.Marshaler 接口。
	f := s.findFolder(in.TextDocument.URI)
	if f == nil {
		return nil
	}

	f.parsedMux.RLock()
	defer f.parsedMux.RUnlock()

	*out = definition(f.doc, in.TextDocument.URI, in.Position)
	return nil
}

func references(doc *ast.APIDoc, uri core.URI, pos core.Position, include bool) (locations []core.Location) {
	r := doc.Search(uri, pos, referencerType)
	if r == nil {
		return
	}

	referencer := r.(ast.Referencer)
	if include {
		locations = append(locations, core.Location{
			URI:   uri,
			Range: referencer.R(),
		})
	}

	for _, ref := range referencer.References() {
		locations = append(locations, ref.Location)
	}

	return
}

func definition(doc *ast.APIDoc, uri core.URI, pos core.Position) []core.Location {
	r := doc.Search(uri, pos, definitionerType)
	if r == nil {
		return []core.Location{}
	}

	return []core.Location{r.(ast.Definitioner).Definition().Location}
}
