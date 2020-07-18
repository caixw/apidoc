// SPDX-License-Identifier: MIT

// Package protocol 协议内容的定义
package protocol

import (
	"github.com/issue9/jsonrpc"

	"github.com/caixw/apidoc/v7/core"
)

// WorkDoneProgressParams a parameter literal used to pass a work done progress token.
type WorkDoneProgressParams struct {
	// An optional token that a server can use to report work done progress.
	WorkDoneToken ProgressToken `json:"workDoneToken,omitempty"`
}

// PartialResultParams a parameter literal used to pass a partial result token
type PartialResultParams struct {
	// An optional token that a server can use to report
	// partial results (e.g. streaming) to the client
	PartialResultToken ProgressToken `json:"partialResultToken,omitempty"`
}

// ProgressToken type ProgressToken = number | string;
type ProgressToken interface{}

// WorkDoneProgressOptions options to signal work done progress support in server capabilities.
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#workDoneProgressOptions
type WorkDoneProgressOptions struct {
	WorkDoneProgress bool `json:"workDoneProgress,omitempty"`
}

// CancelParams The base protocol offers support for request cancellation
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#cancelRequest
type CancelParams struct {
	// The request id to cancel.
	ID *jsonrpc.ID
}

// MarkupContent literal represents a string value which content is interpreted base on its
// kind flag. Currently the protocol supports `plaintext` and `markdown` as markup kinds.
//
// If the kind is `markdown` then the value can contain fenced code blocks like in GitHub issues.
// See https://help.github.com/articles/creating-and-highlighting-code-blocks/#syntax-highlighting
//
// Here is an example how such a string can be constructed using JavaScript / TypeScript:
// ```typescript
// let markdown: MarkdownContent = {
//  kind: MarkupKind.Markdown,
//	value: [
//		'# Header',
//		'Some text',
//		'```typescript',
//		'someCode();',
//		'```'
//	].join('\n')
// };
// ```
//
// *Please Note* that clients might sanitize the return markdown. A client could decide to
// remove HTML from the markdown to avoid script execution.
type MarkupContent struct {
	// The type of the Markup
	Kind MarkupKind `json:"kind"`

	// The content itself
	Value string `json:"value"`
}

// TextEdit a textual edit applicable to a text document.
type TextEdit struct {
	// The range of the text document to be manipulated. To insert
	// text into a document create a range where start === end.
	Range core.Range `json:"range"`

	// The string to be inserted. For delete operations use an empty string.
	NewText string `json:"newText"`
}

// Command represents a reference to a command
//
// Provides a title which will be used to represent a command in the UI.
// Commands are identified by a string identifier.
// The recommended way to handle commands is to implement their execution
// on the server side if the client and server provides the corresponding capabilities.
// Alternatively the tool extension code could handle the command.
// The protocol currently doesn’t specify a set of well-known commands.
type Command struct {
	// Title of the command, like `save`.
	Title string `json:"title"`

	// The identifier of the actual command handler.
	Command string `json:"command"`

	// Arguments that the command handler should be invoked with.
	Arguments []interface{} `json:"arguments,omitempty"`
}
