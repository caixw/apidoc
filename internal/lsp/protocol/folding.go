// SPDX-License-Identifier: MIT

package protocol

import "github.com/caixw/apidoc/v7/internal/xmlenc"

// 代码关折叠块的种类
const (
	FoldingRangeKindComment = "comment" // Folding range for a comment
	FoldingRangeKindImports = "imports" // Folding range for a imports or includes
	FoldingRangeKindRegion  = "region"  // Folding range for a region (e.g. `#region`)
)

// FoldingRangeClientCapabilities 定义客户对代码拆叠功能的支持情况
type FoldingRangeClientCapabilities struct {
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
}

// FoldingRangeParams 由用户传递的 textDocument/foldingRange 参数
type FoldingRangeParams struct {
	WorkDoneProgressParams
	PartialResultParams

	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// FoldingRange represents a folding range
type FoldingRange struct {
	// The zero-based line number from where the folded range starts.
	StartLine int `json:"startLine"`

	// The zero-based character offset from where the folded range starts. If not defined, defaults to the length of the start line.
	//
	// 0 是一个合法的值，所以只能采用指针类型表示空值。
	StartCharacter *int `json:"startCharacter,omitempty"`

	// The zero-based line number where the folded range ends.
	EndLine int `json:"endLine"`

	// The zero-based character offset before the folded range ends. If not defined, defaults to the length of the end line.
	//
	// 0 是一个合法的值，所以只能采用指针类型表示空值。
	EndCharacter *int `json:"endCharacter,omitempty"`

	// Describes the kind of the folding range such as `comment` or `region`. The kind
	// is used to categorize folding ranges and used by commands like 'Fold all comments'. See
	// [FoldingRangeKind](#FoldingRangeKind) for an enumeration of standardized kinds.
	Kind string `json:"kind,omitempty"`
}

// BuildFoldingRange 根据参数构建 FoldingRange 实例
func BuildFoldingRange(base xmlenc.Base, lineFoldingOnly bool) FoldingRange {
	item := FoldingRange{
		StartLine: base.Location.Range.Start.Line,
		EndLine:   base.Location.Range.End.Line,
		Kind:      FoldingRangeKindComment,
	}

	if lineFoldingOnly {
		item.StartCharacter = &base.Location.Range.Start.Character
		item.EndCharacter = &base.Location.Range.End.Character
	}

	return item
}
