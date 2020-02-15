// SPDX-License-Identifier: MIT

package protocol

// InitializeParams 初始化请求的参数
type InitializeParams struct {
	WorkDoneProgressParams

	// The process Id of the parent process that started
	// the server. Is null if the process has not been started by another process.
	// If the parent process is not alive then the server should exit (see exit notification) its process.
	ProcessID int `json:"processId,omitempty"`

	// The rootPath of the workspace. Is null if no folder is open.
	//
	// @deprecated in favour of rootUri.
	RootPath string `json:"rootPath,omitempty"`

	// The rootUri of the workspace. Is null if no
	// folder is open. If both `rootPath` and `rootUri` are set `rootUri` wins.
	RootURI DocumentURI `json:"rootUri,omitempty"`

	// User provided initialization options.
	InitializationOptions interface{} `json:"initializationOptions,omitempty"`

	// The capabilities provided by the client (editor or tool)
	Capabilities ClientCapabilities `json:"capabilities"`

	// The initial trace setting. If omitted trace is disabled ('off').
	//
	// 枚举值：off, messages, verbose
	Trace string `json:"trace,omitempty"`

	// The workspace folders configured in the client when the server starts.
	// This property is only available if the client supports workspace folders.
	// It can be `null` if the client supports workspace folders but none are
	// configured.
	//
	// Since 3.6.0
	WorkspaceFolders *WorkspaceFolder `json:"workspaceFolders,omitempty"`

	// Information about the client
	//
	// @since 3.15.0
	ClientInfo *ServerInfo `json:"clientInfo,omitempty"`
}

// ServerInfo information about the client or server
//
// @since 3.15.0
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// WorkspaceRootPath 从 RootURI 和 RootPath 中获取正确的值
func (p *InitializeParams) WorkspaceRootPath() string {
	if p.RootURI != "" {
		return string(p.RootURI)
	}
	return p.RootPath
}

// InitializeResult initialize 服务的返回结构
type InitializeResult struct {
	// The capabilities the language server provides.
	Capabilities ServerCapabilities `json:"capabilities"`

	// Information about the server.
	//
	// @since 3.15.0
	ServerInfo *ServerInfo `json:"serverInfo,omitempty"`
}

// InitializedParams initialized 服务传递的参数
type InitializedParams struct{}

type ClientCapabilities struct {
	// Workspace specific client capabilities.
	Workspace struct {
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
	} `json:"workspace,omitempty"`

	// Text document specific client capabilities.
	TextDocument TextDocumentClientCapabilities `json:"textDocument,omitempty"`

	// Experimental client capabilities.
	Experimental interface{} `json:"experimental,omitempty"`
}
