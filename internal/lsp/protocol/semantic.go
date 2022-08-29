// SPDX-License-Identifier: MIT

package protocol

// TokenFormat the protocol defines an additional token format capability to allow future extensions of the format
//
// The only format that is currently specified is relative expressing that the tokens are described using relative positions
type TokenFormat string

// TokenFormat支持的常量
const (
	TokenFormatRelative = "relative"
)

// SemanticTokensClientCapabilities 客户端的支持情况
type SemanticTokensClientCapabilities struct {
	// Whether implementation supports dynamic registration. If this is set to `true`
	// the client supports the new `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
	// return value for the corresponding server capability as well.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// Which requests the client supports and might send to the server.
	Requests struct {
		// The client will send the `textDocument/semanticTokens/range` request if
		// the server provides a corresponding handler.
		//
		// bool | {}
		Range any `json:"range,omitempty"`

		// The client will send the `textDocument/semanticTokens/full` request if
		// the server provides a corresponding handler.
		//
		// bool | SemanticTokensOptions
		Full any `json:"full,omitempty"`
	} `json:"requests"`

	// The token types that the client supports.
	TokenTypes []string `json:"tokenTypes"`

	// The token modifiers that the client supports.
	TokenModifiers []string `json:"tokenModifiers"`

	// The formats the clients supports.
	Formats []TokenFormat `json:"formats"`
}

// SemanticTokensOptions 服务端有关 SemanticTokens 的支持情况
type SemanticTokensOptions struct {
	WorkDoneProgressOptions
	// The legend used by the server
	Legend SemanticTokensLegend `json:"legend"`

	// Server supports providing semantic tokens for a specific range of a document.
	//
	// bool | {}
	Range any `json:"range,omitempty"`

	// Server supports providing semantic tokens for a full document.
	//
	// bool | SemanticTokensOptions
	Full any `json:"full,omitempty"`
}

type SemanticTokensLegend struct {
	// The token types a server uses.
	TokenTypes []string `json:"tokenTypes"`

	// The token modifiers a server uses.
	TokenModifiers []string `json:"tokenModifiers"`
}

// SemanticTokensParams textDocument/semanticTokens 入口参数
type SemanticTokensParams struct {
	WorkDoneProgressParams
	PartialResultParams

	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// SemanticTokens textDocument/semanticTokens 返回参数
type SemanticTokens struct {
	// An optional result id. If provided and clients support delta updating
	// the client will include the result id in the next semantic token request.
	// A server can then instead of computing all semantic tokens again simply
	// send a delta.
	ResultID string `json:"resultId,omitempty"`

	// The actual tokens.
	Data []int `json:"data"`
}

// SemanticTokensRegistrationOptions textDocument/semanticTokens/full 返回参数
type SemanticTokensRegistrationOptions struct {
	TextDocumentRegistrationOptions
	SemanticTokensOptions
	StaticRegistrationOptions
}
