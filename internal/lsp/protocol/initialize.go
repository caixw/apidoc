// SPDX-License-Identifier: MIT

package protocol

import "path"

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
	Trace string `json:"trace,omitempty"`

	// The workspace folders configured in the client when the server starts.
	// This property is only available if the client supports workspace folders.
	// It can be `null` if the client supports workspace folders but none are
	// configured.
	//
	// Since 3.6.0
	//
	// 在客户端支持工作区的情况下，RootURI 和 RootPath 的内容为 WorkspaceFolders 的第一个元素，
	// 否则 WorkspaceFolders 为空。
	// RootURI 和 RootPath 的区别在于，RootURI 会带协议，比如 file:///path 而 RootPath 为 /path
	WorkspaceFolders []WorkspaceFolder `json:"workspaceFolders,omitempty"`

	// Information about the client
	//
	// @since 3.15.0
	ClientInfo *ServerInfo `json:"clientInfo,omitempty"`
}

func (p *InitializeParams) Folders() []WorkspaceFolder {
	if len(p.WorkspaceFolders) > 0 {
		return p.WorkspaceFolders
	}

	if p.RootURI != "" {
		return []WorkspaceFolder{{
			Name: path.Base(string(p.RootURI)),
			URI:  p.RootURI,
		}}
	}
	if p.RootPath != "" {
		return []WorkspaceFolder{{
			Name: path.Base(p.RootPath),
			URI:  DocumentURI(p.RootPath),
		}}
	}

	return nil
}

// ServerInfo 终端的信息，同时用于描述服务和客户端。
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
		WorkspaceEdit WorkspaceEditClientCapabilities `json:"workspaceEdit,omitempty"`

		// Capabilities specific to the `workspace/didChangeConfiguration` notification.
		DidChangeConfiguration DidChangeConfigurationClientCapabilities `json:"didChangeConfiguration,omitempty"`

		// Capabilities specific to the `workspace/didChangeWatchedFiles` notification.
		DidChangeWatchedFiles DidChangeConfigurationClientCapabilities `json:"didChangeWatchedFiles,omitempty"`

		// Capabilities specific to the `workspace/symbol` request.
		Symbol WorkspaceSymbolClientCapabilities `json:"symbol,omitempty"`

		// Capabilities specific to the `workspace/executeCommand` request.
		ExecuteCommand ExecuteCommandClientCapabilities `json:"executeCommand,omitempty"`

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

type ExecuteCommandClientCapabilities struct {
	// Execute command supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}
