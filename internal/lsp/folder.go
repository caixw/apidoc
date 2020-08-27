// SPDX-License-Identifier: MIT

package lsp

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/issue9/sliceutil"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

// 表示项目文件夹
type folder struct {
	protocol.WorkspaceFolder
	doc       *ast.APIDoc
	h         *core.MessageHandler
	cfg       *build.Config
	srv       *server
	loadError error // 加载过程中的出错信息

	parsedMux sync.RWMutex

	// 保存着错误和警告的信息
	diagnostics map[core.URI]*protocol.PublishDiagnosticsParams
}

func (f *folder) close() {
	f.clearDiagnostics()
	f.h.Stop()
}

func (f *folder) messageHandler(msg *core.Message) {
	err, ok := msg.Message.(*core.Error)
	if !ok {
		f.srv.erro.Println(fmt.Sprintf("获得了非 core.Error 错误 %#v", msg.Message))
		return
	}

	if p, found := f.diagnostics[err.Location.URI]; found && p != nil {
		cnt := sliceutil.Count(p.Diagnostics, func(i int) bool {
			return p.Diagnostics[i].Range.Equal(err.Location.Range)
		})
		if cnt == 0 {
			p.AppendDiagnostic(err, msg.Type)
		}
		return
	}
	p := protocol.NewPublishDiagnosticsParams(err.Location.URI)
	p.AppendDiagnostic(err, msg.Type)
	f.diagnostics[err.Location.URI] = p
}

func (s *server) appendFolders(folders ...protocol.WorkspaceFolder) {
	for _, f := range folders {
		ff := s.openFolder(f)
		s.folders = append(s.folders, ff)
	}
}

// 解析 f 目录，并生成 folder
//
// 在成功加载文档之后，会通过 apidoc/outline 通知客户端新的数据。
func (s *server) openFolder(f protocol.WorkspaceFolder) (ff *folder) {
	ff = &folder{
		WorkspaceFolder: f,
		doc:             &ast.APIDoc{},
		srv:             s,
		diagnostics:     make(map[core.URI]*protocol.PublishDiagnosticsParams, 5),
	}
	ff.h = core.NewMessageHandler(ff.messageHandler)

	cfg, err := build.LoadConfig(f.URI)
	if errors.Is(err, os.ErrNotExist) { // 找不到配置文件
		if cfg, err = build.DetectConfig(f.URI, true); err != nil {
			ff.loadError = err
			s.printErr(ff.loadError)
			return ff
		}
	} else if err != nil {
		ff.loadError = err
		s.printErr(ff.loadError)
		return ff
	}
	ff.cfg = cfg

	ff.doc.ParseBlocks(ff.h, func(blocks chan core.Block) {
		build.ParseInputs(blocks, ff.h, ff.cfg.Inputs...)
	})

	if err = s.apidocOutline(ff); err != nil {
		s.printErr(err)
	}

	s.textDocumentPublishDiagnostics(ff)

	return ff
}

func (s *server) findFolder(uri core.URI) *folder {
	s.workspaceMux.RLock()
	defer s.workspaceMux.RUnlock()

	for _, f := range s.folders {
		if f.Contains(uri) {
			return f
		}
	}
	return nil
}
