// SPDX-License-Identifier: MIT

package protocol

type ServerCapabilities struct {
	// Defines how text documents are synced. Is either a detailed structure defining each notification or
	// for backwards compatibility the TextDocumentSyncKind number.
	// If omitted it defaults to `TextDocumentSyncKind.None`.
	//
	// TextDocumentSyncOptions | number;
	TextDocumentSync interface{} `json:"textDocumentSync"`

	// The server provides hover support.
	HoverProvider bool `json:"hoverProvider,omitempty"`

	// The server provides completion support.
	CompletionProvider CompletionOptions `json:"completionProvider,omitempty"`

	// The server provides signature help support.
	SignatureHelpProvider SignatureHelpOptions `json:"signatureHelpProvider,omitempty"`

	// The server provides goto definition support.
	DefinitionProvider bool `json:"definitionProvider,omitempty"`

	// The server provides Goto Type Definition support.
	//
	// Since 3.6.0
	//
	// boolean | (TextDocumentRegistrationOptions & StaticRegistrationOptions);
	TypeDefinitionProvider interface{} `json:"typeDefinitionProvider,omitempty"`

	// The server provides Goto Implementation support.
	//
	// Since 3.6.0
	//
	// boolean | (TextDocumentRegistrationOptions & StaticRegistrationOptions);
	ImplementationProvider interface{} `json:"implementationProvider,omitempty"`

	// The server provides find references support.
	ReferencesProvider bool `json:"referencesProvider,omitempty"`

	// The server provides document highlight support.
	DocumentHighlightProvider bool `json:"documentHighlightProvider,omitempty"`

	// The server provides document symbol support.
	DocumentSymbolProvider bool `json:"documentSymbolProvider,omitempty"`

	// The server provides workspace symbol support.
	WorkspaceSymbolProvider bool `json:"workspaceSymbolProvider,omitempty"`

	// The server provides code actions. The `CodeActionOptions` return type is only
	// valid if the client signals code action literal support via the property
	// `textDocument.codeAction.codeActionLiteralSupport`.
	//
	// boolean | CodeActionOptions;
	CodeActionProvider interface{} `json:"codeActionProvider,omitempty"`

	// The server provides code lens.
	CodeLensProvider CodeLensOptions `json:"codeLensProvider,omitempty"`

	// The server provides document formatting.
	DocumentFormattingProvider bool `json:"documentFormattingProvider,omitempty"`

	// The server provides document range formatting.
	DocumentRangeFormattingProvider bool `json:"documentRangeFormattingProvider,omitempty"`

	// The server provides document formatting on typing.
	DocumentOnTypeFormattingProvider DocumentOnTypeFormattingOptions `json:"documentOnTypeFormattingProvider,omitempty"`

	// The server provides rename support. RenameOptions may only be
	// specified if the client states that it supports
	// `prepareSupport` in its initial `initialize` request.
	RenameProvider interface{} `json:"renameProvider,omitempty"` // boolean | RenameOptions;

	// The server provides document link support.
	DocumentLinkProvider DocumentLinkOptions `json:"documentLinkProvider,omitempty"`

	// The server provides color provider support.
	//
	// Since 3.6.0
	//
	// boolean | ColorProviderOptions | (ColorProviderOptions & TextDocumentRegistrationOptions & StaticRegistrationOptions);
	ColorProvider interface{} `json:"colorProvider,omitempty"`

	// The server provides folding provider support.
	//
	// Since 3.10.0
	//
	// boolean | FoldingRangeProviderOptions | (FoldingRangeProviderOptions & TextDocumentRegistrationOptions & StaticRegistrationOptions);
	FoldingRangeProvider interface{} `json:"foldingRangeProvider,omitempty"`

	// The server provides go to declaration support.
	//
	// Since 3.14.0
	//
	// boolean | (TextDocumentRegistrationOptions & StaticRegistrationOptions);
	DeclarationProvider interface{} `json:"declarationProvider,omitempty"`

	// The server provides execute command support.
	ExecuteCommandProvider ExecuteCommandOptions `json:"executeCommandProvider,omitempty"`

	// Workspace specific server capabilities
	Workspace struct {
		// The server supports workspace folder.
		//
		// Since 3.6.0
		WorkspaceFolders struct {
			// The server has support for workspace folders
			Supported bool `json:"supported,omitempty"`

			// Whether the server wants to receive workspace folder
			// change notifications.
			//
			// If a strings is provided the string is treated as a ID
			// under which the notification is registered on the client
			// side. The ID can be used to unregister for these events
			// using the `client/unregisterCapability` request.
			ChangeNotifications string `json:"changeNotifications,omitempty"` // string | boolean
		} `json:"workspaceFolders,omitempty"`
	} `json:"workspace,omitempty"`

	// Experimental server capabilities.
	Experimental interface{} `json:"experimental,omitempty"`
}

// Completion options.
type CompletionOptions struct {
	// The server provides support to resolve additional
	// information for a completion item.
	ResolveProvider bool `json:"resolveProvider,omitempty"`

	// The characters that trigger completion automatically.
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

// Signature help options.
type SignatureHelpOptions struct {
	// The characters that trigger signature help automatically.
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

// Code Action options.
type CodeActionOptions struct {
	// CodeActionKinds that this server may return.
	//
	// The list of kinds may be generic, such as `CodeActionKind.Refactor`, or the server
	// may list out every specific kind they provide.
	CodeActionKinds []CodeActionKind `json:"codeActionKinds,omitempty"`
}

// Code Lens options.
type CodeLensOptions struct {
	// Code lens has a resolve provider as well.
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

// Execute command options.
type ExecuteCommandOptions struct {
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

// Folding range provider options.
type FoldingRangeProviderOptions struct{}

// Static registration options to be returned in the initialize request.
type StaticRegistrationOptions struct {
	// The id used to register the request. The id can be used to deregister
	// the request again. See also Registration#id.
	ID string `json:"id,omitempty"`
}
