// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/issue9/sliceutil"

	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

// 自定义的服务端下发通知 apidoc/outline
func (s *server) apidocOutline(f *folder) error {
	if f.loadError != nil {
		return s.Notify("apidoc/outline", &protocol.APIDocOutline{Err: f.loadError.Error()})
	}

	if outline := protocol.BuildAPIDocOutline(f.WorkspaceFolder, f.doc); outline != nil {
		return s.Notify("apidoc/outline", outline)
	}
	return nil
}

// 由客户端发给服务端的刷新通知 apidoc/refreshOutline
//
// 之后刷新的内容会通过 apidoc/outline 通知客户端。
func (s *server) apidocRefreshOutline(notify bool, in *protocol.WorkspaceFolder, out *interface{}) error {
	sliceutil.Delete(s.folders, func(i int) bool {
		return s.folders[i].URI == in.URI
	})

	s.appendFolders(*in)
	return nil
}
