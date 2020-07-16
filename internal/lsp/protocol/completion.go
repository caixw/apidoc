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
		TagSupport struct {
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
