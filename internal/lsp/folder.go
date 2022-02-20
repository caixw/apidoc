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
	cfg       *build.Config
	srv       *server
	loadError error // 加载过程中的出错信息
	noConfig  bool
	h         *core.MessageHandler

	parsedMux sync.RWMutex // 解析 doc 时需要的锁

	// 保存着错误和警告的信息
	diagnostics map[core.URI]*protocol.PublishDiagnosticsParams
}

func (f *folder) close() {
	f.clearDiagnostics()
	if f.h != nil {
		f.h.Stop()
	}
}

func (f *folder) messageHandler(msg *core.Message) {
	err, ok := msg.Message.(*core.Error)
	if !ok {
		f.srv.erro.Println(fmt.Sprintf("获得了非 core.Error 错误 %#v", msg.Message))
		return
	}

	if p, found := f.diagnostics[err.Location.URI]; found && p != nil {
		cnt := sliceutil.Count(p.Diagnostics, func(i protocol.Diagnostic) bool {
			return i.Range.Equal(err.Location.Range)
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
	for _, ff := range folders {
		f := &folder{
			WorkspaceFolder: ff,
			doc:             &ast.APIDoc{},
			srv:             s,
			diagnostics:     make(map[core.URI]*protocol.PublishDiagnosticsParams, 5),
		}
		f.refresh(false)
		s.folders = append(s.folders, f)
	}
}

// 刷新项目
//
// 默认情况下，没有配置文件不会解析项目，但是在 force 为 true 时，会强制解析项目内容。
func (f *folder) refresh(force bool) {
	f.loadError = nil
	cfg, err := build.LoadConfig(f.URI)
	if errors.Is(err, os.ErrNotExist) { // 找不到配置文件
		f.noConfig = true
		if force {
			if cfg, err = build.DetectConfig(f.URI, true); err != nil {
				f.loadError = err
				f.cfg = cfg
				f.srv.printErr(f.loadError)
				return
			}
		} else {
			f.loadError = err // 仅在不强制刷新项目的情况下，才会将错误保存至 f.loadError
			return
		}
	} else if err != nil {
		f.loadError = err
		f.srv.printErr(f.loadError)
		return
	}
	f.cfg = cfg

	f.doc = &ast.APIDoc{}
	f.clearDiagnostics()

	if f.h == nil {
		f.h = core.NewMessageHandler(f.messageHandler)
	}

	f.doc.ParseBlocks(f.h, func(blocks chan core.Block) {
		build.ParseInputs(blocks, f.h, f.cfg.Inputs...)
	})

	if err = f.srv.apidocOutline(f); err != nil {
		f.srv.printErr(err)
	}

	f.srv.textDocumentPublishDiagnostics(f)
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
