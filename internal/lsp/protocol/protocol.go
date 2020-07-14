// SPDX-License-Identifier: MIT

// Package protocol 协议内容的定义
package protocol

import "github.com/issue9/jsonrpc"

// WorkDoneProgressParams a parameter literal used to pass a work done progress token.
type WorkDoneProgressParams struct {
	// An optional token that a server can use to report work done progress.
	WorkDoneToken ProgressToken `json:"workDoneToken,omitempty"`
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
