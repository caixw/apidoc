// SPDX-License-Identifier: MIT

package protocol

import "github.com/caixw/apidoc/v7/core"

// ResourceOperationKind the kind of resource operations supported by the client.
type ResourceOperationKind string

const (
	// ResourceOperationKindCreate supports creating new files and folders.
	ResourceOperationKindCreate ResourceOperationKind = "create"

	// ResourceOperationKindRename supports renaming existing files and folders.
	ResourceOperationKindRename ResourceOperationKind = "rename"

	// ResourceOperationKindDelete supports deleting existing files and folders.
	ResourceOperationKindDelete ResourceOperationKind = "delete"
)

// FailureHandlingKind 定义出错后的处理方式
type FailureHandlingKind string

const (

	// FailureHandlingKindAbort applying the workspace change is simply aborted
	// if one of the changes provided fails.
	// All operations executed before the failing operation stay executed.
	FailureHandlingKindAbort FailureHandlingKind = "abort"

	// FailureHandlingKindTransactional all operations are executed transactionally.
	// That means they either all succeed or no changes at all are applied to the workspace.
	FailureHandlingKindTransactional FailureHandlingKind = "transactional"

	// FailureHandlingKindTextOnlyTransactional if the workspace edit contains only textual file changes they are executed transactionally.
	// If resource changes (create, rename or delete file) are part of the change the failure
	// handling strategy is abort.
	FailureHandlingKindTextOnlyTransactional FailureHandlingKind = "textOnlyTransactional"

	// FailureHandlingKindUndo The client tries to undo the operations already executed.
	// But there is no guarantee that this succeeds.
	FailureHandlingKindUndo FailureHandlingKind = "undo"
)

// WorkspaceClientCapabilities 客户有关 workspace 的支持情况
type WorkspaceClientCapabilities struct {
	// The client supports applying batch edits to the workspace by supporting
	// the request 'workspace/applyEdit'
	ApplyEdit bool `json:"applyEdit,omitempty"`

	// Capabilities specific to `WorkspaceEdit`s
	WorkspaceEdit *WorkspaceEditClientCapabilities `json:"workspaceEdit,omitempty"`

	// Capabilities specific to the `workspace/didChangeConfiguration` notification.
	DidChangeConfiguration *DidChangeConfigurationClientCapabilities `json:"didChangeConfiguration,omitempty"`

	// Capabilities specific to the `workspace/didChangeWatchedFiles` notification.
	DidChangeWatchedFiles *DidChangeConfigurationClientCapabilities `json:"didChangeWatchedFiles,omitempty"`

	// The client has support for workspace folders.
	//
	// Since 3.6.0
	WorkspaceFolders bool `json:"workspaceFolders,omitempty"`

	// The client supports `workspace/configuration` requests.
	//
	// Since 3.6.0
	Configuration bool `json:"configuration,omitempty"`
}

// WorkspaceProvider 服务端有关 workspace 的支持情况
type WorkspaceProvider struct {
	// The server supports workspace folder.
	//
	// Since 3.6.0
	WorkspaceFolders *WorkspaceFoldersServerCapabilities `json:"workspaceFolders,omitempty"`
}

// WorkspaceFolder 项目文件夹
type WorkspaceFolder struct {
	// The associated URI for this workspace folder.
	URI core.URI `json:"uri"`

	// The name of the workspace folder. Used to refer to this
	// workspace folder in the user interface.
	Name string `json:"name"`
}

// DidChangeWorkspaceFoldersParams workspace/didChangeWorkspaceFolders 参数
type DidChangeWorkspaceFoldersParams struct {
	// The actual workspace folder change event.
	Event WorkspaceFoldersChangeEvent `json:"event"`
}

// WorkspaceFoldersChangeEvent the workspace folder change event.
type WorkspaceFoldersChangeEvent struct {
	// The array of added workspace folders
	Added []WorkspaceFolder `json:"added"`

	// The array of the removed workspace folders
	Removed []WorkspaceFolder `json:"removed"`
}

// WorkspaceFoldersServerCapabilities 服务端有关项目文件夹的支持情况
type WorkspaceFoldersServerCapabilities struct {
	// The server has support for workspace folders
	Supported bool `json:"supported,omitempty"`

	// Whether the server wants to receive workspace folder
	// change notifications.
	//
	// If a string is provided, the string is treated as an ID
	// under which the notification is registered on the client
	// side. The ID can be used to unregister for these events
	// using the `client/unregisterCapability` request.
	//
	// string | boolean;
	ChangeNotifications interface{} `json:"changeNotifications,omitempty"`
}

// WorkspaceEditClientCapabilities the capabilities of a workspace edit has evolved over the time. Clients can describe their support using the following client capability
type WorkspaceEditClientCapabilities struct {
	// The client supports versioned document changes in `WorkspaceEdit`s
	DocumentChanges bool `json:"documentChanges,omitempty"`

	// The resource operations the client supports. Clients should at least
	// support 'create', 'rename' and 'delete' files and folders.
	//
	// @since 3.13.0
	ResourceOperations []ResourceOperationKind `json:"resourceOperations,omitempty"`

	// The failure handling strategy of a client if applying the workspace edit fails.
	//
	// @since 3.13.0
	FailureHandling FailureHandlingKind `json:"failureHandling,omitempty"`
}
