// SPDX-License-Identifier: MIT

// Package protocol 协议内容的定义
package protocol

// Position in a text document expressed as zero-based line and zero-based character offset
//
// A position is between two characters like an ‘insert’ cursor in a editor.
// Special values like for example -1 to denote the end of a line are not supported.
type Position struct {
	// Line position in a document (zero-based).
	Line int `json:"line"`

	// Character offset on a line in a document (zero-based). Assuming that the line is
	// represented as a string, the `character` value represents the gap between the
	// `character` and `character + 1`.
	//
	// If the character value is greater than the line length it defaults back to the
	// line length.
	Character int `json:"character"`
}

// Range a range in a text document expressed as (zero-based) start and end positions
//
// A range is comparable to a selection in an editor. Therefore the end position is exclusive.
// If you want to specify a range that contains a line including the line ending character(s)
// then use an end position denoting the start of the next line. For example:
//  {
//     start: { line: 5, character: 23 },
//     end : { line 6, character : 0 }
//  }
type Range struct {
	// The range's start position.
	Start Position `json:"start"`

	// The range's end position.
	End Position `json:"end"`
}

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
