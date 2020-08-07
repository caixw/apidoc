// SPDX-License-Identifier: MIT

package protocol

// ServerCapabilities 服务端的兼容列表
type ServerCapabilities struct {
	// Defines how text documents are synced. Is either a detailed structure defining each notification or
	// for backwards compatibility the TextDocumentSyncKind number.
	// If omitted it defaults to `TextDocumentSyncKind.None`.
	//
	// ServerCapabilitiesTextDocumentSyncOptions | TextDocumentSyncKind;
	TextDocumentSync interface{} `json:"textDocumentSync"`

	// The server provides completion support.
	CompletionProvider *CompletionOptions `json:"completionProvider,omitempty"`

	// The server provides hover support.
	//
	// boolean | HoverOptions
	HoverProvider bool `json:"hoverProvider,omitempty"`

	// The server provides goto definition support.
	//
	// boolean | DefinitionOptions;
	DefinitionProvider interface{} `json:"definitionProvider,omitempty"`

	// The server provides code actions. The `CodeActionOptions` return type is only
	// valid if the client signals code action literal support via the property
	// `textDocument.codeAction.codeActionLiteralSupport`.
	//
	// boolean | CodeActionOptions;
	CodeActionProvider interface{} `json:"codeActionProvider,omitempty"`

	// The server provides document range formatting.
	//
	// boolean | DocumentRangeFormattingOptions;
	DocumentRangeFormattingProvider interface{} `json:"documentRangeFormattingProvider,omitempty"`

	// The server provides rename support. RenameOptions may only be
	// specified if the client states that it supports
	// `prepareSupport` in its initial `initialize` request.
	//
	// boolean | RenameOptions;
	RenameProvider interface{} `json:"renameProvider,omitempty"`

	// The server provides folding provider support.
	//
	// Since 3.10.0
	//
	// boolean | FoldingRangeOptions | FoldingRangeRegistrationOptions;
	FoldingRangeProvider bool `json:"foldingRangeProvider,omitempty"`

	// The server provides folding provider support.
	//
	// Since 3.16.0
	//
	// SemanticTokensOptions | SemanticTokensRegistrationOptions
	SemanticTokensProvider interface{} `json:"semanticTokensProvider,omitempty"`

	// The server provides workspace symbol support.
	WorkspaceSymbolProvider bool `json:"workspaceSymbolProvider,omitempty"`

	// Workspace specific server capabilities
	Workspace *WorkspaceProvider `json:"workspace,omitempty"`

	// Experimental server capabilities.
	Experimental interface{} `json:"experimental,omitempty"`
}

// CodeActionOptions Code Action options.
type CodeActionOptions struct {
	WorkDoneProgressOptions

	// CodeActionKinds that this server may return.
	//
	// The list of kinds may be generic, such as `CodeActionKind.Refactor`, or the server
	// may list out every specific kind they provide.
	CodeActionKinds []CodeActionKind `json:"codeActionKinds,omitempty"`
}

// SaveOptions Save options.
type SaveOptions struct {
	// The client is supposed to include the content on save.
	IncludeText bool `json:"includeText,omitempty"`
}

// StaticRegistrationOptions static registration options to be returned in the initialize request.
type StaticRegistrationOptions struct {
	// The id used to register the request. The id can be used to deregister
	// the request again. See also Registration#id.
	ID string `json:"id,omitempty"`
}
