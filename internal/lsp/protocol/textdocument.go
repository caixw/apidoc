// SPDX-License-Identifier: MIT

package protocol

// Text document specific client capabilities.
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
	Completion struct {
		// Whether completion supports dynamic registration.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		// The client supports the following `CompletionItem` specific
		// capabilities.
		CompletionItem struct {
			// The client supports snippets as insert text.
			//
			// A snippet can define tab stops and placeholders with `$1`, `$2`
			// and `${3:foo}`. `$0` defines the final tab stop, it defaults to
			// the end of the snippet. Placeholders with equal identifiers are linked,
			// that is typing in one will update others too.
			SnippetSupport bool `json:"snippetSupport,omitempty"`

			// The client supports commit characters on a completion item.
			CommitCharactersSupport bool `json:"commitCharactersSupport,omitempty"`

			// The client supports the following content formats for the documentation
			// property. The order describes the preferred format of the client.
			DocumentationFormat []MarkupKind `json:"documentationFormat,omitempty"`

			// The client supports the deprecated property on a completion item.
			DeprecatedSupport bool `json:"deprecatedSupport,omitempty"`

			// The client supports the preselect property on a completion item.
			PreselectSupport bool `json:"preselectSupport,omitempty"`
		} `json:"completionItem,omitempty"`

		CompletionItemKind struct {
			/**
			 * The completion item kind values the client supports. When this
			 * property exists the client also guarantees that it will
			 * handle values outside its set gracefully and falls back
			 * to a default value when unknown.
			 *
			 * If this property is not present the client only supports
			 * the completion items kinds from `Text` to `Reference` as defined in
			 * the initial version of the protocol.
			 */
			ValueSet []CompletionItemKind `json:"valueSet,omitempty"`
		} `json:"completionItemKind,omitempty"`

		/**
		 * The client supports to send additional context information for a
		 * `textDocument/completion` request.
		 */
		ContextSupport bool `json:"contextSupport,omitempty"`
	} `json:"completion,omitempty"`

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
				LabelOffsetSupport bool `json:"labelOffsetSupport,omitmepty"`
			} `json:"parameterInformation,omitempty"`
		} `json:"signatureInformation,omitempty"`
	} `json:"signatureHelp,omitempty"`

	// Capabilities specific to the `textDocument/references`
	References DynamicRegistration `json:"references,omitempty"`

	// Capabilities specific to the `textDocument/documentHighlight`
	DocumentHighlight DynamicRegistration `json:"documentHighlight,omitempty"`

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
	Formatting DynamicRegistration `json:"formatting,omitempty"`

	// Capabilities specific to the `textDocument/rangeFormatting`
	RangeFormatting DynamicRegistration `json:"rangeFormatting,omitempty"`

	// Capabilities specific to the `textDocument/onTypeFormatting`
	OnTypeFormatting DynamicRegistration `json:"onTypeFormatting,omitempty"`

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
	CodeLens DynamicRegistration `json:"codeLens,omitempty"`

	// Capabilities specific to the `textDocument/documentLink`
	DocumentLink DynamicRegistration `json:"documentLink,omitempty"`

	// Capabilities specific to the `textDocument/documentColor` and the
	// `textDocument/colorPresentation` request.
	//
	// Since 3.6.0
	//
	// If ColorProvider.DynamicRegistration is set to `true`
	// the client supports the new `(ColorProviderOptions & TextDocumentRegistrationOptions & StaticRegistrationOptions)`
	// return value for the corresponding server capability as well.
	ColorProvider DynamicRegistration `json:"colorProvider,omitempty"`

	// Capabilities specific to the `textDocument/rename`
	Rename struct {
		// Whether rename supports dynamic registration.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
		// The client supports testing for validity of rename operation before execution.
		PrepareSupport bool `json:"prepareSupport,omitempty"`
	} `json:"rename,omitempty"`

	// Capabilities specific to `textDocument/publishDiagnostics`.
	PublishDiagnostics struct {
		// Whether the clients accepts diagnostics with related information.
		RelatedInformation bool `json:"relatedInformation,omitempty"`
	} `json:"publishDiagnostics,omitempty"`
	// Capabilities specific to `textDocument/foldingRange` requests.
	//
	// Since 3.10.0
	FoldingRange struct {
		// Whether implementation supports dynamic registration for folding range providers. If this is set to `true`
		// the client supports the new `(FoldingRangeProviderOptions & TextDocumentRegistrationOptions & StaticRegistrationOptions)`
		// return value for the corresponding server capability as well.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		// The maximum number of folding ranges that the client prefers to receive per document. The value serves as a
		// hint, servers are free to follow the limit.
		RangeLimit int `json:"rangeLimit,omitempty"`

		// If set, the client signals that it only supports folding complete lines. If set, client will
		// ignore specified `startCharacter` and `endCharacter` properties in a FoldingRange.
		LineFoldingOnly bool `json:"lineFoldingOnly,omitempty"`
	} `json:"foldingRange,omitempty"`
}

type DynamicRegistration struct {
	// Whether formatting supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}
