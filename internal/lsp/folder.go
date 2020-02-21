// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/caixw/apidoc/v6/doc"
	"github.com/caixw/apidoc/v6/internal/lsp/protocol"
)

// 表示项目文件夹
type folder struct {
	protocol.WorkspaceFolder
	doc *doc.Doc
}

func (f *folder) close() error {
	// TODO
	return nil
}

func (s *server) appendFolders(folders ...protocol.WorkspaceFolder) {
	for _, f := range folders {
		s.folders = append(s.folders, &folder{
			WorkspaceFolder: f,
			doc:             doc.New(),
		})
	}
}
