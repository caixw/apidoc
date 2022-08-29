// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

// 自定义的服务端下发通知 apidoc/outline
func (s *server) apidocOutline(f *folder) error {
	if f.loadError != nil {
		return s.Notify("apidoc/outline", &protocol.APIDocOutline{Err: f.loadError.Error(), NoConfig: f.noConfig})
	}

	if outline := protocol.BuildAPIDocOutline(f.WorkspaceFolder, f.doc); outline != nil {
		outline.NoConfig = f.noConfig
		return s.Notify("apidoc/outline", outline)
	}
	return nil
}

// 由客户端发给服务端的刷新通知 apidoc/refreshOutline
//
// 之后刷新的内容会通过 apidoc/outline 通知客户端。
func (s *server) apidocRefreshOutline(notify bool, in *protocol.WorkspaceFolder, out *any) error {
	if f := s.findFolder(in.URI); f != nil {
		f.parsedMux.RLock()
		defer f.parsedMux.RUnlock()
		f.refresh(true)
	}
	return nil
}

// 由客户端发给服务端用于创建项目的配置文件 apidoc/detect
//
// 之后新内容会通过 apidoc/outline 通知客户端。
func (s *server) apidocDetect(notify bool, in *protocol.APIDocDetectParams, out *protocol.APIDocDetectResult) error {
	f := s.findFolder(in.TextDocument.URI)
	if f == nil {
		return nil
	}

	cfg, err := build.DetectConfig(in.TextDocument.URI, in.Recursive)
	if err != nil {
		out.Error = err.Error()
		return nil
	}
	if err := cfg.Save(in.TextDocument.URI); err != nil {
		return err
	}
	return s.apidocOutline(f)
}
