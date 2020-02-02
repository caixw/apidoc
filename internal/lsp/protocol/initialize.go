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
	Workspace WorkspaceClientCapabilities `json:"workspace,omitempty"`

	// Text document specific client capabilities.
	TextDocument TextDocumentClientCapabilities `json:"textDocument,omitempty"`

	// Experimental client capabilities.
	Experimental interface{} `json:"experimental,omitempty"`
}
