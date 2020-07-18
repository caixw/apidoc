// SPDX-License-Identifier: MIT

package protocol

// CompletionItemKind the kind of a completion entry.
type CompletionItemKind int

// CompletionItemKind 的各类枚举值
const (
	CompletionItemKindText CompletionItemKind = iota + 1
	CompletionItemKindMethod
	CompletionItemKindFunction
	CompletionItemKindConstructor
	CompletionItemKindField
	CompletionItemKindVariable
	CompletionItemKindClass
	CompletionItemKindInterface
	CompletionItemKindModule
	CompletionItemKindProperty
	CompletionItemKindUnit
	CompletionItemKindValue
	CompletionItemKindEnum
	CompletionItemKindKeyword
	CompletionItemKindSnippet
	CompletionItemKindColor
	CompletionItemKindFile
	CompletionItemKindReference
	CompletionItemKindFolder
	CompletionItemKindEnumMember
	CompletionItemKindConstant
	CompletionItemKindStruct
	CompletionItemKindEvent
	CompletionItemKindOperator
	CompletionItemKindTypeParameter
)

// CompletionItemTag are extra annotations that tweak the rendering of a completion item.
//
// @since 3.15.0
type CompletionItemTag int

// CompletionItemTagDeprecated render a completion as obsolete, usually using a strike-out.
const CompletionItemTagDeprecated CompletionItemTag = 1

// CompletionTriggerKind how a completion was triggered
type CompletionTriggerKind int

// CompletionTriggerKind 定义的常量
const (
	// Completion was triggered by typing an identifier (24x7 code
	// complete), manual invocation (e.g Ctrl+Space) or via API.
	CompletionTriggerKindInvoked CompletionTriggerKind = 1

	// Completion was triggered by a trigger character specified by
	// the `triggerCharacters` properties of the `CompletionRegistrationOptions`.
	CompletionTriggerKindTriggerCharacter CompletionTriggerKind = 2

	// Completion was re-triggered as the current completion list is incomplete.
	CompletionTriggerKindTriggerForIncompleteCompletions CompletionTriggerKind = 3
)

// InsertTextFormat defines whether the insert text in a completion
// item should be interpreted as plain text or a snippet.
type InsertTextFormat int

// InsertTextFormat 的可用常量
const (
	// The primary text to be inserted is treated as a plain string.
	InsertTextFormatPlainText InsertTextFormat = 1

	// The primary text to be inserted is treated as a snippet.
	//
	// A snippet can define tab stops and placeholders with `$1`, `$2`
	// and `${3:foo}`. `$0` defines the final tab stop, it defaults to
	// the end of the snippet. Placeholders with equal identifiers are linked,
	// that is typing in one will update others too.
	InsertTextFormatSnippet InsertTextFormat = 2
)

// CompletionClientCapabilities 客户端有关自动完成所支持的功能定义
type CompletionClientCapabilities struct {
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

		// Client supports the tag property on a completion item. Clients supporting
		// tags have to handle unknown tags gracefully. Clients especially need to
		// preserve unknown tags when sending a completion item back to the server in
		// a resolve call.
		//
		// @since 3.15.0
		TagSupport *struct {
			// The tags supported by the client.
			ValueSet []CompletionItemTag `json:"valueSet"`
		} `json:"tagSupport,omitempty"`

		// Client support insert replace edit to control different behavior if a
		// completion item is inserted in the text or should replace text.
		//
		// @since 3.16.0 - Proposed state
		InsertReplaceSupport bool `json:"insertReplaceSupport,omitempty"`

		// Client supports to resolve `additionalTextEdits` in the `completionItem/resolve`
		// request. So servers can postpone computing them.
		//
		// @since 3.16.0 - Proposed state
		ResolveAdditionalTextEditsSupport bool `json:"resolveAdditionalTextEditsSupport,omitempty"`
	} `json:"completionItem,omitempty"`

	CompletionItemKind struct {
		// The completion item kind values the client supports. When this
		// property exists the client also guarantees that it will
		// handle values outside its set gracefully and falls back
		// to a default value when unknown.
		//
		// If this property is not present the client only supports
		// the completion items kinds from `Text` to `Reference` as defined in
		// the initial version of the protocol.
		ValueSet []CompletionItemKind `json:"valueSet,omitempty"`
	} `json:"completionItemKind,omitempty"`

	// The client supports to send additional context information for a
	// `textDocument/completion` request.
	ContextSupport bool `json:"contextSupport,omitempty"`
}

// CompletionOptions 服务端返回的有关对自动完成支持性的描述
type CompletionOptions struct {
	WorkDoneProgressOptions

	// The server provides support to resolve additional
	// information for a completion item.
	ResolveProvider bool `json:"resolveProvider,omitempty"`

	// The characters that trigger completion automatically.
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`

	// The list of all possible characters that commit a completion. This field can be used
	// if clients don't support individual commit characters per completion item. See
	// `ClientCapabilities.textDocument.completion.completionItem.commitCharactersSupport`.
	//
	// If a server provides both `allCommitCharacters` and commit characters on an individual
	// completion item the ones on the completion item win.
	//
	// @since 3.2.0
	AllCommitCharacters []string `json:"allCommitCharacters,omitempty"`
}

// CompletionParams textDocument/completion 的参数
type CompletionParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
	// The completion context. This is only available if the client specifies
	// to send this using `ClientCapabilities.textDocument.completion.contextSupport === true`
	Context *CompletionContext `json:"context,omitempty"`
}

// CompletionContext contains additional information about the context in which a completion request is triggered.
type CompletionContext struct {
	// How the completion was triggered.
	TriggerKind CompletionTriggerKind `json:"triggerKind"`

	// The trigger character (a single character) that has trigger code complete.
	// Is undefined if `triggerKind !== CompletionTriggerKind.TriggerCharacter`
	TriggerCharacter string `json:"triggerCharacter,omitempty"`
}

// CompletionList represents a collection of [completion items](#CompletionItem) to be presented
// in the editor.
type CompletionList struct {
	// This list it not complete. Further typing should result in recomputing this list.
	IsIncomplete bool `json:"isIncomplete"`

	// The completion items.
	Items []CompletionItem `json:"items"`
}

// CompletionItem completion items
type CompletionItem struct {
	// The label of this completion item. By default also the text that is
	// inserted when selecting this completion.
	Label string `json:"label"`

	// The kind of this completion item. Based of the kind
	// an icon is chosen by the editor. The standardized set
	// of available values is defined in `CompletionItemKind`.
	Kind CompletionItemKind `json:"kind,omitempty"`

	// Tags for this completion item.
	//
	// @since 3.15.0
	Tags []CompletionItemTag `json:"tags,omitempty"`

	// A human-readable string with additional information
	// about this item, like type or symbol information.
	Detail string `json:"detail,omitempty"`

	// A human-readable string that represents a doc-comment.
	Documentation MarkupContent `json:"documentation,omitempty"`

	// Select this item when showing.
	//
	// *Note* that only one completion item can be selected and that the
	// tool / client decides which item that is. The rule is that the *first*
	// item of those that match best is selected.
	Preselect bool `json:"preselect,omitempty"`

	// A string that should be used when comparing this item
	// with other items. When `falsy` the label is used.
	SortText string `json:"sortText,omitempty"`

	// A string that should be used when filtering a set of
	// completion items. When `falsy` the label is used.
	FilterText string `json:"filterText,omitempty"`

	// A string that should be inserted into a document when selecting
	// this completion. When `falsy` the label is used.
	//
	// The `insertText` is subject to interpretation by the client side.
	// Some tools might not take the string literally. For example
	// VS Code when code complete is requested in this example `con<cursor position>`
	// and a completion item with an `insertText` of `console` is provided it
	// will only insert `sole`. Therefore it is recommended to use `textEdit` instead
	// since it avoids additional client side interpretation.
	InsertText string `json:"insertText,omitempty"`

	// The format of the insert text. The format applies to both the `insertText` property
	// and the `newText` property of a provided `textEdit`. If omitted defaults to
	// `InsertTextFormat.PlainText`.
	InsertTextFormat InsertTextFormat `json:"insertTextFormat,omitempty"`

	// An edit which is applied to a document when selecting this completion. When an edit is provided the value of
	// `insertText` is ignored.
	//
	// *Note:* The range of the edit must be a single line range and it must contain the position at which completion
	// has been requested.
	TextEdit *TextEdit `json:"textEdit,omitempty"`

	// An optional array of additional text edits that are applied when
	// selecting this completion. Edits must not overlap (including the same insert position)
	// with the main edit nor with themselves.
	//
	// Additional text edits should be used to change text unrelated to the current cursor position
	// (for example adding an import statement at the top of the file if the completion item will
	// insert an unqualified type).
	AdditionalTextEdits []TextEdit `json:"additionalTextEdits,omitempty"`

	// An optional set of characters that when pressed while this completion is active will accept it first and
	// then type that character. *Note* that all commit characters should have `length=1` and that superfluous
	// characters will be ignored.
	CommitCharacters []string `json:"commitCharacters,omitempty"`

	// An optional command that is executed *after* inserting this completion. *Note* that
	// additional modifications to the current document should be described with the
	// additionalTextEdits-property.
	Command *Command `json:"command,omitempty"`

	// A data entry field that is preserved on a completion item between
	// a completion and a completion resolve request.
	Data interface{} `json:"data,omitempty"`
}
