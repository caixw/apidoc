// SPDX-License-Identifier: MIT

package protocol

import (
	"encoding/json"

	"github.com/caixw/apidoc/v7/core"
)

// HoverCapabilities 客户端有关 hover 功能的描述
type HoverCapabilities struct {
	// Whether hover supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// The client supports the follow content formats for the content
	// property. The order describes the preferred format of the client.
	ContentFormat []MarkupKind `json:"contentFormat,omitempty"`
}

// HoverParams textDocument/hover 发送的参数
type HoverParams struct {
	WorkDoneProgressParams
	TextDocumentPositionParams
}

// Hover textDocument/hover 的返回结果
type Hover struct {
	// The hover's content
	Contents MarkupContent `json:"contents"`

	// An optional range is a range inside a text document
	// that is used to visualize a hover, e.g. by changing the background color.
	Range core.Range `json:"range"`
}

// MarshalJSON 允许在 hover 为空值是返回 null
func (h *Hover) MarshalJSON() ([]byte, error) {
	if h.Contents.Kind == "" {
		return json.Marshal(nil)
	}

	type hoverShadow Hover
	shadow := (*hoverShadow)(h)
	return json.Marshal(shadow)
}
