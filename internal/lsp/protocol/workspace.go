// SPDX-License-Identifier: MIT

package protocol

type WorkspaceFolder struct {
	// The associated URI for this workspace folder.
	URI DocumentURI `json:"uri"`

	// The name of the workspace folder. Used to refer to this
	// workspace folder in the user interface.
	Name string `json:"name"`
}

// DidChangeWorkspaceFoldersParams workspace/didChangeWorkspaceFolders 参数
type DidChangeWorkspaceFoldersParams struct {
	// The actual workspace folder change event.
	Event WorkspaceFoldersChangeEvent `json:"event"`
}

// The workspace folder change event.
type WorkspaceFoldersChangeEvent struct {
	// The array of added workspace folders
	Added []WorkspaceFolder `json:"added"`

	// The array of the removed workspace folders
	Removed []WorkspaceFolder `json:"removed"`
}

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
