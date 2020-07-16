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
	CompletionProvider CompletionOptions `json:"completionProvider,omitempty"`

	// The server provides hover support.
	//
	// boolean | HoverOptions
	HoverProvider interface{} `json:"hoverProvider,omitempty"`

	// The server provides signature help support.
	SignatureHelpProvider SignatureHelpOptions `json:"signatureHelpProvider,omitempty"`

	// The server provides go to declaration support.
	//
	// Since 3.14.0
	//
	// boolean | DeclarationOptions | DeclarationRegistrationOptions
	DeclarationProvider interface{} `json:"declarationProvider,omitempty"`

	// The server provides goto definition support.
	//
	// boolean | DefinitionOptions;
	DefinitionProvider interface{} `json:"definitionProvider,omitempty"`

	// The server provides Goto Type Definition support.
	//
	// Since 3.6.0
	//
	// boolean | TypeDefinitionOptions | TypeDefinitionRegistrationOptions;
	TypeDefinitionProvider interface{} `json:"typeDefinitionProvider,omitempty"`

	// The server provides Goto Implementation support.
	//
	// Since 3.6.0
	//
	// boolean | ImplementationOptions | ImplementationRegistrationOptions;
	ImplementationProvider interface{} `json:"implementationProvider,omitempty"`

	// The server provides find references support.
	//
	// boolean | ReferenceOptions
	ReferencesProvider interface{} `json:"referencesProvider,omitempty"`

	// The server provides document highlight support.
	//
	// boolean | DocumentHighlightOptions
	DocumentHighlightProvider interface{} `json:"documentHighlightProvider,omitempty"`

	// The server provides document symbol support.
	//
	// boolean | DocumentSymbolOptions;
	DocumentSymbolProvider interface{} `json:"documentSymbolProvider,omitempty"`

	// The server provides code actions. The `CodeActionOptions` return type is only
	// valid if the client signals code action literal support via the property
	// `textDocument.codeAction.codeActionLiteralSupport`.
	//
	// boolean | CodeActionOptions;
	CodeActionProvider interface{} `json:"codeActionProvider,omitempty"`

	// The server provides code lens.
	CodeLensProvider CodeLensOptions `json:"codeLensProvider,omitempty"`

	// The server provides document link support.
	DocumentLinkProvider DocumentLinkOptions `json:"documentLinkProvider,omitempty"`

	// The server provides color provider support.
	//
	// Since 3.6.0
	//
	// boolean | DocumentColorOptions | DocumentColorRegistrationOptions;
	ColorProvider interface{} `json:"colorProvider,omitempty"`

	// The server provides document formatting.
	//
	// boolean | DocumentFormattingOptions;
	DocumentFormattingProvider interface{} `json:"documentFormattingProvider,omitempty"`

	// The server provides document range formatting.
	//
	// boolean | DocumentRangeFormattingOptions;
	DocumentRangeFormattingProvider interface{} `json:"documentRangeFormattingProvider,omitempty"`

	// The server provides document formatting on typing.
	DocumentOnTypeFormattingProvider DocumentOnTypeFormattingOptions `json:"documentOnTypeFormattingProvider,omitempty"`

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
	FoldingRangeProvider interface{} `json:"foldingRangeProvider,omitempty"`

	// The server provides execute command support.
	ExecuteCommandProvider ExecuteCommandOptions `json:"executeCommandProvider,omitempty"`

	// The server provides selection range support.
	//
	// @since 3.15.0
	//
	// boolean | SelectionRangeOptions | SelectionRangeRegistrationOptions
	SelectionRangeProvider interface{} `json:"selectionRangeProvider,omitempty"`

	// The server provides workspace symbol support.
	WorkspaceSymbolProvider bool `json:"workspaceSymbolProvider,omitempty"`

	// Workspace specific server capabilities
	Workspace struct {
		// The server supports workspace folder.
		//
		// Since 3.6.0
		WorkspaceFolders WorkspaceFoldersServerCapabilities `json:"workspaceFolders,omitempty"`
	} `json:"workspace,omitempty"`

	// Experimental server capabilities.
	Experimental interface{} `json:"experimental,omitempty"`
}

type HoverOptions struct {
	WorkDoneProgressOptions
}

type DeclarationOptions struct {
	WorkDoneProgressOptions
}

type DeclarationRegistrationOptions struct {
	DeclarationOptions
	TextDocumentRegistrationOptions
	StaticRegistrationOptions
}

type DefinitionOptions struct {
	WorkDoneProgressOptions
}

type TypeDefinitionOptions struct {
	WorkDoneProgressOptions
}

type TypeDefinitionRegistrationOptions struct {
	TextDocumentRegistrationOptions
	TypeDefinitionOptions
	StaticRegistrationOptions
}

type ImplementationOptions struct {
	WorkDoneProgressOptions
}

type ImplementationRegistrationOptions struct {
	TextDocumentRegistrationOptions
	ImplementationOptions
	StaticRegistrationOptions
}

type ReferenceOptions struct {
	WorkDoneProgressOptions
}

// SignatureHelpOptions Signature help options.
type SignatureHelpOptions struct {
	WorkDoneProgressOptions

	// The characters that trigger signature help automatically.
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`

	// List of characters that re-trigger signature help.
	//
	// These trigger characters are only active when signature help is already showing. All trigger characters
	// are also counted as re-trigger characters.
	//
	// @since 3.15.0
	RetriggerCharacters []string `json:"retriggerCharacters,omitempty"`
}

// Code Action options.
type CodeActionOptions struct {
	WorkDoneProgressOptions

	// CodeActionKinds that this server may return.
	//
	// The list of kinds may be generic, such as `CodeActionKind.Refactor`, or the server
	// may list out every specific kind they provide.
	CodeActionKinds []CodeActionKind `json:"codeActionKinds,omitempty"`
}

// Code Lens options.
type CodeLensOptions struct {
	WorkDoneProgressOptions

	// Code lens has a resolve provider as well.
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

// Execute command options.
type ExecuteCommandOptions struct {
	WorkDoneProgressOptions

	// The commands to be executed on the server
	Commands []string `json:"commands"`
}

// Save options.
type SaveOptions struct {
	// The client is supposed to include the content on save.
	IncludeText bool `json:"includeText,omitempty"`
}

// Color provider options.
type ColorProviderOptions struct{}

// StaticRegistrationOptions static registration options to be returned in the initialize request.
type StaticRegistrationOptions struct {
	// The id used to register the request. The id can be used to deregister
	// the request again. See also Registration#id.
	ID string `json:"id,omitempty"`
}
