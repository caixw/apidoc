// SPDX-License-Identifier: MIT

package protocol

import (
	"path"

	"github.com/caixw/apidoc/v7/core"
)

// InitializationOptions 用户需要提交的自定义初始化参数
type InitializationOptions struct {
	// 客户端的本地化 ID
	//
	// 服务端会根据此值决定提示内容如何翻译。
	// 如果提交的 Locale 无法识别或是服务端不支持，
	// 则会采用服务端的默认值，即 locale.DefaultLocaleID。
	Locale string `json:"locale,omitempty"`
}

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
	RootURI core.URI `json:"rootUri,omitempty"`

	// User provided initialization options.
	InitializationOptions *InitializationOptions `json:"initializationOptions,omitempty"`

	// The capabilities provided by the client (editor or tool)
	Capabilities ClientCapabilities `json:"capabilities"`

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

// Folders 获取客户端当前打开的所有项目
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
			URI:  core.URI(p.RootPath),
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

// ClientCapabilities 客户端的兼容列表
type ClientCapabilities struct {
	// Workspace specific client capabilities.
	Workspace *WorkspaceClientCapabilities `json:"workspace,omitempty"`

	// Text document specific client capabilities.
	TextDocument TextDocumentClientCapabilities `json:"textDocument,omitempty"`

	// Experimental client capabilities.
	Experimental any `json:"experimental,omitempty"`
}
