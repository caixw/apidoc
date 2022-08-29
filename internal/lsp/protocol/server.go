// SPDX-License-Identifier: MIT

package protocol

// ServerCapabilities 服务端的兼容列表
type ServerCapabilities struct {
	// Defines how text documents are synced.
	//
	// Is either a detailed structure defining each notification or
	// for backwards compatibility the TextDocumentSyncKind number.
	// If omitted it defaults to `TextDocumentSyncKind.None`.
	//
	// ServerCapabilitiesTextDocumentSyncOptions | TextDocumentSyncKind;
	TextDocumentSync *ServerCapabilitiesTextDocumentSyncOptions `json:"textDocumentSync"`

	// The server provides completion support.
	CompletionProvider *CompletionOptions `json:"completionProvider,omitempty"`

	// The server provides hover support.
	HoverProvider bool `json:"hoverProvider,omitempty"`

	// The server provides goto definition support.
	DefinitionProvider bool `json:"definitionProvider,omitempty"`

	// The server provides find references support.
	ReferencesProvider bool `json:"referencesProvider,omitempty"`

	// The server provides folding provider support.
	//
	// Since 3.10.0
	FoldingRangeProvider bool `json:"foldingRangeProvider,omitempty"`

	// The server provides folding provider support.
	//
	// Since 3.16.0
	//
	// SemanticTokensOptions | SemanticTokensRegistrationOptions
	SemanticTokensProvider any `json:"semanticTokensProvider,omitempty"`

	// The server provides workspace symbol support.
	WorkspaceSymbolProvider bool `json:"workspaceSymbolProvider,omitempty"`

	// Workspace specific server capabilities
	Workspace *WorkspaceProvider `json:"workspace,omitempty"`

	// Experimental server capabilities.
	Experimental any `json:"experimental,omitempty"`
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
