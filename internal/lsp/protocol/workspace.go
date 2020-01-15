// SPDX-License-Identifier: MIT

package protocol

type WorkspaceFolder struct {
	// The associated URI for this workspace folder.
	URI DocumentURI `json:"uri"`

	// The name of the workspace folder. Used to refer to this
	// workspace folder in the user interface.
	Name string `json:"name"`
}

// Workspace specific client capabilities.
type WorkspaceClientCapabilities struct {
	// The client supports applying batch edits to the workspace by supporting
	// the request 'workspace/applyEdit'
	ApplyEdit bool `json:"applyEdit,omitempty"`

	// Capabilities specific to `WorkspaceEdit`s
	WorkspaceEdit struct {
		// The client supports versioned document changes in `WorkspaceEdit`s
		DocumentChanges bool `json:"documentChanges,omitempty"`

		// The resource operations the client supports. Clients should at least
		// support 'create', 'rename' and 'delete' files and folders.
		ResourceOperations []ResourceOperationKind `json:"resourceOperations,omitempty"`

		// The failure handling strategy of a client if applying the workspace edit fails.
		FailureHandling FailureHandlingKind `json:"failureHandling,omitempty"`
	} `json:"workspaceEdit,omitempty"`

	// Capabilities specific to the `workspace/didChangeConfiguration` notification.
	DidChangeConfiguration DynamicRegistration `json:"didChangeConfiguration,omitempty"`

	// Capabilities specific to the `workspace/didChangeWatchedFiles` notification.
	//
	// DidChangeWatchedFiles.DynamicRegistration:
	// Did change watched files notification supports dynamic registration. Please note
	// that the current protocol doesn't support static configuration for file changes
	// from the server side.
	DidChangeWatchedFiles DynamicRegistration `json:"didChangeWatchedFiles,omitempty"`

	// Capabilities specific to the `workspace/symbol` request.
	Symbol struct {
		// Symbol request supports dynamic registration.
		DynamicRegistration bool `json:"dynamicRegistration,omitempty"`

		// Specific capabilities for the `SymbolKind` in the `workspace/symbol` request.
		SymbolKind struct {
			// The symbol kind values the client supports. When this
			// property exists the client also guarantees that it will
			// handle values outside its set gracefully and falls back
			// to a default value when unknown.
			//
			// If this property is not present the client only supports
			// the symbol kinds from `File` to `Array` as defined in
			// the initial version of the protocol.
			ValueSet SymbolKind `json:"valueSet,omitempty"`
		} `json:"symbolKind,omitempty"`
	} `json:"symbol,omitempty"`

	// Capabilities specific to the `workspace/executeCommand` request.
	ExecuteCommand DynamicRegistration `json:"executeCommand,omitempty"`

	// The client has support for workspace folders.
	//
	// Since 3.6.0
	WorkspaceFolders bool `json:"workspaceFolders,omitempty"`

	// The client supports `workspace/configuration` requests.
	//
	// Since 3.6.0
	Configuration bool `json:"configuration,omitempty"`
}
