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

// TextDocumentItem an item to transfer a text document from the client to the server.
type TextDocumentItem struct {
	// The text document's URI.
	URI core.URI `json:"uri"`

	// The text document's language identifier.
	LanguageID string `json:"languageId"`

	// The version number of this document (it will increase after each
	// change, including undo/redo).
	Version int `json:"version"`

	// The content of the opened text document.
	Text string `json:"text"`
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
	Synchronization struct {
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
	Hover struct {
		// Whether hover supports dynamic registration.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		// The client supports the follow content formats for the content
		// property. The order describes the preferred format of the client.
		ContentFormat []MarkupKind `json:"contentFormat,omitempty"`
	} `json:"hover,omitempty"`

	// Capabilities specific to the `textDocument/signatureHelp`
	SignatureHelp struct {
		// Whether signature help supports dynamic registration.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		// The client supports the following `SignatureInformation` specific properties.
		SignatureInformation struct {
			// The client supports the follow content formats for the documentation
			// property. The order describes the preferred format of the client.
			DocumentationFormat []MarkupKind `json:"documentationFormat,omitempty"`

			// Client capabilities specific to parameter information.
			ParameterInformation struct {
				// The client supports processing label offsets instead of a
				// simple label string.
				//
				// Since 3.14.0
				LabelOffsetSupport bool `json:"labelOffsetSupport,omitempty"`
			} `json:"parameterInformation,omitempty"`
		} `json:"signatureInformation,omitempty"`
	} `json:"signatureHelp,omitempty"`

	// Capabilities specific to the `textDocument/references`
	References DidChangeConfigurationClientCapabilities `json:"references,omitempty"`

	// Capabilities specific to the `textDocument/documentHighlight`
	DocumentHighlight DidChangeConfigurationClientCapabilities `json:"documentHighlight,omitempty"`

	// Capabilities specific to the `textDocument/documentSymbol`
	DocumentSymbol struct {
		// Whether document symbol supports dynamic registration.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		// Specific capabilities for the `SymbolKind`.
		SymbolKind struct {
			// The symbol kind values the client supports. When this
			// property exists the client also guarantees that it will
			// handle values outside its set gracefully and falls back
			// to a default value when unknown.
			//
			// If this property is not present the client only supports
			// the symbol kinds from `File` to `Array` as defined in
			// the initial version of the protocol.
			ValueSet []SymbolKind `json:"valueSet,omitempty"`
		} `json:"symbolKind,omitempty"`

		// The client supports hierarchical document symbols.
		HierarchicalDocumentSymbolSupport bool `json:"hierarchicalDocumentSymbolSupport,omitempty"`
	} `json:"documentSymbol,omitempty"`

	// Capabilities specific to the `textDocument/formatting`
	Formatting DidChangeConfigurationClientCapabilities `json:"formatting,omitempty"`

	// Capabilities specific to the `textDocument/rangeFormatting`
	RangeFormatting DidChangeConfigurationClientCapabilities `json:"rangeFormatting,omitempty"`

	// Capabilities specific to the `textDocument/onTypeFormatting`
	OnTypeFormatting DidChangeConfigurationClientCapabilities `json:"onTypeFormatting,omitempty"`

	// Capabilities specific to the `textDocument/declaration`
	Declaration struct {
		// Whether declaration supports dynamic registration. If this is set to `true`
		// the client supports the new `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
		// return value for the corresponding server capability as well.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		// The client supports additional metadata in the form of declaration links.
		//
		// Since 3.14.0
		LinkSupport bool `json:"linkSupport,omitempty"`
	} `json:"declaration,omitempty"`

	// Capabilities specific to the `textDocument/definition`.
	//
	// Since 3.14.0
	Definition struct {
		// Whether definition supports dynamic registration.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		// The client supports additional metadata in the form of definition links.
		LinkSupport bool `json:"linkSupport,omitempty"`
	} `json:"definition,omitempty"`

	// Capabilities specific to the `textDocument/typeDefinition`
	//
	// Since 3.6.0
	TypeDefinition struct {
		// Whether typeDefinition supports dynamic registration. If this is set to `true`
		// the client supports the new `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
		// return value for the corresponding server capability as well.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		// The client supports additional metadata in the form of definition links.
		//
		// Since 3.14.0
		LinkSupport bool `json:"linkSupport,omitempty"`
	} `json:"typeDefinition,omitempty"`

	// Capabilities specific to the `textDocument/implementation`.
	//
	// Since 3.6.0
	Implementation struct {
		// Whether implementation supports dynamic registration. If this is set to `true`
		// the client supports the new `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
		// return value for the corresponding server capability as well.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		// The client supports additional metadata in the form of definition links.
		//
		// Since 3.14.0
		LinkSupport bool `json:"linkSupport,omitempty"`
	} `json:"implementation,omitempty"`

	// Capabilities specific to the `textDocument/codeAction`
	CodeAction struct {
		// Whether code action supports dynamic registration.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		// The client support code action literals as a valid
		// response of the `textDocument/codeAction` request.
		//
		// Since 3.8.0
		CodeActionLiteralSupport struct {
			// The code action kind is support with the following value set.
			CodeActionKind struct {
				// The code action kind values the client supports. When this
				// property exists the client also guarantees that it will
				// handle values outside its set gracefully and falls back
				// to a default value when unknown.
				ValueSet []CodeActionKind `json:"valueSet"`
			} `json:"codeActionKind"`
		} `json:"codeActionLiteralSupport,omitempty"`
	} `json:"codeAction,omitempty"`

	// Capabilities specific to the `textDocument/codeLens`
	CodeLens DidChangeConfigurationClientCapabilities `json:"codeLens,omitempty"`

	// Capabilities specific to the `textDocument/documentLink`
	DocumentLink DidChangeConfigurationClientCapabilities `json:"documentLink,omitempty"`

	// Capabilities specific to the `textDocument/documentColor` and the
	// `textDocument/colorPresentation` request.
	//
	// Since 3.6.0
	//
	// If ColorProvider.DidChangeConfigurationClientCapabilities is set to `true`
	// the client supports the new `(ColorProviderOptions & TextDocumentRegistrationOptions & StaticRegistrationOptions)`
	// return value for the corresponding server capability as well.
	ColorProvider DidChangeConfigurationClientCapabilities `json:"colorProvider,omitempty"`

	// Capabilities specific to the `textDocument/rename`
	Rename struct {
		// Whether rename supports dynamic registration.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
		// The client supports testing for validity of rename operation before execution.
		PrepareSupport bool `json:"prepareSupport,omitempty"`
	} `json:"rename,omitempty"`

	// Capabilities specific to `textDocument/publishDiagnostics`.
	PublishDiagnostics *PublishDiagnosticsClientCapabilities `json:"publishDiagnostics,omitempty"`

	// Capabilities specific to `textDocument/foldingRange` requests.
	//
	// Since 3.10.0
	FoldingRange *FoldingRangeClientCapabilities `json:"foldingRange,omitempty"`
}

type DidChangeConfigurationClientCapabilities struct {
	// Whether formatting supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// DocumentOnTypeFormattingOptions format document on type options.
type DocumentOnTypeFormattingOptions struct {
	// A character on which formatting should be triggered, like `}`.
	FirstTriggerCharacter string `json:"firstTriggerCharacter"`

	// More trigger characters.
	MoreTriggerCharacter []string `json:"moreTriggerCharacter,omitempty"`
}

// RenameOptions Rename options
type RenameOptions struct {
	WorkDoneProgressOptions

	// Renames should be checked and tested before being executed.
	PrepareProvider bool `json:"prepareProvider,omitempty"`
}

// DocumentLinkOptions Document link options.
type DocumentLinkOptions struct {
	WorkDoneProgressOptions

	// Document links have a resolve provider as well.
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

type ServerCapabilitiesTextDocumentSyncOptions struct {
	// Open and close notifications are sent to the server.
	// If omitted open close notification should not be sent.
	OpenClose bool `json:"openClose,omitempty"`

	// Change notifications are sent to the server. See TextDocumentSyncKind.None, TextDocumentSyncKind.Full
	// and TextDocumentSyncKind.Incremental. If omitted it defaults to TextDocumentSyncKind.None.
	Change TextDocumentSyncKind `json:"change,omitempty"`
}

type TextDocumentSyncOptions struct {
	ServerCapabilitiesTextDocumentSyncOptions

	// If present will save notifications are sent to the server.
	// If omitted the notification should not be sent.
	WillSave bool `json:"willSave,omitempty"`

	// If present will save wait until requests are sent to the server.
	// If omitted the request should not be sent.
	WillSaveWaitUntil bool `json:"willSaveWaitUntil,omitempty"`
	// If present save notifications are sent to the server.
	// If omitted the notification should not be sent.
	Save SaveOptions `json:"save,omitempty"`
}

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
//  { language: 'typescript', scheme: 'file' }
//  { language: 'json', pattern: '**/package.json' }
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

// DidSaveTextDocumentParams textDocument/didSave 的参数
type DidSaveTextDocumentParams struct {
	// The document that was saved.
	TextDocument TextDocumentIdentifier `json:"textDocument"`

	// Optional the content when saved. Depends on the includeText value
	// when the save notification was requested.
	Text string `json:"text,omitempty"`
}

type DocumentHighlightOptions struct {
	WorkDoneProgressOptions
}

type DocumentSymbolOptions struct {
	WorkDoneProgressOptions
}

type DocumentColorOptions struct {
	WorkDoneProgressOptions
}

type DocumentColorRegistrationOptions struct {
	TextDocumentRegistrationOptions
	StaticRegistrationOptions
	DocumentColorOptions
}

type DocumentFormattingOptions struct {
	WorkDoneProgressOptions
}

type DocumentRangeFormattingOptions struct {
	WorkDoneProgressOptions
}

type SelectionRangeOptions struct {
	WorkDoneProgressOptions
}

type SelectionRangeRegistrationOptions struct {
	SelectionRangeOptions
	TextDocumentRegistrationOptions
	StaticRegistrationOptions
}

// DidOpenTextDocumentParams textDocument/didOpen 通知发送的参数
type DidOpenTextDocumentParams struct {
	// The document that was opened.
	TextDocument TextDocumentItem `json:"textDocument"`
}

// DidCloseTextDocumentParams textDocument/didClose 通知发送的参数
type DidCloseTextDocumentParams struct {
	// The document that was closed.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
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

// HoverParams textDocument/hover 发送的参数
type HoverParams struct {
	WorkDoneProgressParams
	TextDocumentPositionParams
}

// Hover textDocument/hover 的返回结果
type Hover struct {
	// The hover's content
	// contents MarkedString | MarkedString[] | MarkupContent;
	Contents interface{} `json:"contents"`

	// An optional range is a range inside a text document
	// that is used to visualize a hover, e.g. by changing the background color.
	Range core.Range `json:"range,omitempty"`
}
