// SPDX-License-Identifier: MIT

package protocol

import "github.com/caixw/apidoc/v7/core"

// TextDocumentIdentifier text documents are identified using a URI.
// On the protocol level, URIs are passed as strings.
// The corresponding JSON structure looks like this:
type TextDocumentIdentifier struct {
	// The text document's URI.
	URI core.URI `json:"uri"`
}

// VersionedTextDocumentIdentifier an identifier to denote a specific version of a text document.
type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier

	// The version number of this document. If a versioned text document identifier
	// is sent from the server to the client and the file is not open in the editor
	// (the server has not received an open notification before) the server can send
	// `null` to indicate that the version is known and the content on disk is the
	// truth (as speced with document content ownership).
	//
	// The version number of a document will increase after each change, including
	// undo/redo. The number doesn't need to be consecutive.
	Version int `json:"version,omitempty"`
}

// TextDocumentPositionParams a parameter literal used in requests to pass a text document and a position inside that document.
type TextDocumentPositionParams struct {
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`

	// The position inside the text document.
	Position core.Position `json:"position"`
}

// TextDocumentClientCapabilities text document specific client capabilities.
type TextDocumentClientCapabilities struct {
	Synchronization *struct {
		// Whether text document synchronization supports dynamic registration.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		// The client supports sending will save notifications.
		WillSave bool `json:"willSave,omitempty"`

		// The client supports sending a will save request and
		// waits for a response providing text edits which will
		// be applied to the document before it is saved.
		WillSaveWaitUntil bool `json:"willSaveWaitUntil,omitempty"`

		// The client supports did save notifications.
		DidSave bool `json:"didSave,omitempty"`
	} `json:"synchronization,omitempty"`

	// Capabilities specific to the `textDocument/completion`
	Completion *CompletionClientCapabilities `json:"completion,omitempty"`

	// Capabilities specific to the `textDocument/hover`
	Hover *HoverCapabilities `json:"hover,omitempty"`

	// Capabilities specific to the `textDocument/textDocument/semanticTokens/*`
	SemanticTokens *SemanticTokensClientCapabilities `json:"semanticTokens,omitempty"`

	// Capabilities specific to the `textDocument/references` request.
	References *ReferenceClientCapabilities `json:"references,omitempty"`

	// Capabilities specific to the `textDocument/definition`.
	//
	// Since 3.14.0
	Definition *DefinitionClientCapabilities `json:"definition,omitempty"`

	// Capabilities specific to `textDocument/publishDiagnostics`.
	PublishDiagnostics *PublishDiagnosticsClientCapabilities `json:"publishDiagnostics,omitempty"`

	// Capabilities specific to `textDocument/foldingRange` requests.
	//
	// Since 3.10.0
	FoldingRange *FoldingRangeClientCapabilities `json:"foldingRange,omitempty"`
}

// ServerCapabilitiesTextDocumentSyncOptions 服务端对文档同步的支持项
type ServerCapabilitiesTextDocumentSyncOptions struct {
	// Open and close notifications are sent to the server.
	// If omitted open close notification should not be sent.
	OpenClose bool `json:"openClose,omitempty"`

	// Change notifications are sent to the server. See TextDocumentSyncKind.None, TextDocumentSyncKind.Full
	// and TextDocumentSyncKind.Incremental. If omitted it defaults to TextDocumentSyncKind.None.
	Change TextDocumentSyncKind `json:"change,omitempty"`
}

// TextDocumentRegistrationOptions General text document registration options
type TextDocumentRegistrationOptions struct {
	// A document selector to identify the scope of the registration. If set to null
	// the document selector provided on the client side will be used.
	DocumentSelector DocumentSelector `json:"documentSelector,omitempty"`
}

// DocumentSelector is the combination of one or more document filters
type DocumentSelector []DocumentFilter

// DocumentFilter denotes a document through properties like language,
// scheme or pattern. An example is a filter that applies to TypeScript files on disk.
// Another example is a filter the applies to JSON files with name package.json:
//
//	{ language: 'typescript', scheme: 'file' }
//	{ language: 'json', pattern: '**/package.json' }
type DocumentFilter struct {
	// A language id, like `typescript`.
	Language string `json:"language,omitempty"`

	// A Uri [scheme](#Uri.scheme), like `file` or `untitled`.
	Scheme string `json:"scheme,omitempty"`

	// A glob pattern, like `*.{ts,js}`.
	//
	// Glob patterns can have the following syntax:
	// - `*` to match one or more characters in a path segment
	// - `?` to match on one character in a path segment
	// - `**` to match any number of path segments, including none
	// - `{}` to group conditions (e.g. `**​/*.{ts,js}` matches all TypeScript and JavaScript files)
	// - `[]` to declare a range of characters to match in a path segment (e.g., `example.[0-9]` to match on `example.0`, `example.1`, …)
	// - `[!...]` to negate a range of characters to match in a path segment (e.g., `example.[!0-9]` to match on `example.a`, `example.b`, but not `example.0`)
	Pattern string `json:"pattern,omitempty"`
}

// DidChangeTextDocumentParams textDocument/didChange 的参数
type DidChangeTextDocumentParams struct {
	// The document that did change. The version number points
	// to the version after all provided content changes have been applied.
	TextDocument VersionedTextDocumentIdentifier `json:"textDocument"`

	// The actual content changes. The content changes describe single state changes
	// to the document. So if there are two content changes c1 (at array index 0) and
	// c2 (at array index 1) for a document in state S then c1 moves the document from
	// S to S' and c2 from S' to S''. So c1 is computed on the state S and c2 is computed
	// on the state S'.
	//
	// To mirror the content of a document using change events use the following approach:
	// - start with the same initial content
	// - apply the 'textDocument/didChange' notifications in the order you recevie them.
	// - apply the `TextDocumentContentChangeEvent`s in a single notification in the order
	//   you receive them.
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// Blocks 返回 core.Block 的列表
func (p *DidChangeTextDocumentParams) Blocks() []core.Block {
	blocks := make([]core.Block, 0, len(p.ContentChanges))
	for _, c := range p.ContentChanges {
		blk := core.Block{
			Data: []byte(c.Text),
			Location: core.Location{
				URI: p.TextDocument.URI,
			},
		}
		if c.Range != nil {
			blk.Location.Range = *c.Range
		}
		blocks = append(blocks, blk)
	}
	return blocks
}

// TextDocumentContentChangeEvent an event describing a change to a text document.
// If range and rangeLength are omitted the new text is considered to be
// the full content of the document.
type TextDocumentContentChangeEvent struct {
	// The range of the document that changed.
	Range *core.Range `json:"range,omitempty"`

	// The new text for the provided range.
	// The new text of the whole document.
	//
	// 如果没有 Range 内容，表示整个文档内容；否则表示 Range 表示的内容
	Text string `json:"text"`

	// The optional length of the range that got replaced.
	//
	// @deprecated use range instead.
	RangeLength int `json:"rangeLength,omitempty"`
}
