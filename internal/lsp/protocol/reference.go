// SPDX-License-Identifier: MIT

package protocol

// ReferenceClientCapabilities 客户端对与 textDocument/references 的支持情况
type ReferenceClientCapabilities struct {
	// Whether references supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// ReferenceParams textDocument/references 的请求参数
type ReferenceParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
	Context struct {
		// Include the declaration of the current symbol.
		IncludeDeclaration bool `json:"includeDeclaration"`
	} `json:"context"`
}
