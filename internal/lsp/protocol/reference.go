// SPDX-License-Identifier: MIT

package protocol

// ReferenceClientCapabilities 客户端对 textDocument/references 的支持情况
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

// DefinitionClientCapabilities 客户端对 textDocument/definition 的支持情况
type DefinitionClientCapabilities struct {
	// Whether definition supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

	// The client supports additional metadata in the form of definition links.
	//
	// @since 3.14.0
	LinkSupport bool `json:"linkSupport,omitempty"`
}

// DefinitionParams textDocument/definition 服务的客户端请求参数
type DefinitionParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}
