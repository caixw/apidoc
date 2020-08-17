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
	"github.com/caixw/apidoc/v7/internal/locale"
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
	errors, warns []*core.Error
}

func (f *folder) close() {
	f.srv.windowLogLogMessage(locale.CloseLSPFolder, f.Name)
	f.errors = f.errors[:0]
	f.warns = f.warns[:0]
	f.h.Stop()
}

func (f *folder) messageHandler(msg *core.Message) {
	err, ok := msg.Message.(*core.Error)
	if !ok {
		f.srv.erro.Println(fmt.Sprintf("获得了非 core.Error 错误 %#v", msg.Message))
		return
	}

	switch msg.Type {
	case core.Erro:
		cnt := sliceutil.Count(f.errors, func(i int) bool {
			return f.errors[i].Location.Equal(err.Location)
		})
		if cnt == 0 {
			f.errors = append(f.errors, err)
		}
	case core.Warn:
		cnt := sliceutil.Count(f.warns, func(i int) bool {
			return f.warns[i].Location.Equal(err.Location)
		})
		if cnt == 0 {
			f.warns = append(f.warns, err)
		}
	case core.Succ, core.Info: // 仅处理错误和警告
	default:
		panic("unreached")
	}
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

	if err := s.apidocOutline(ff); err != nil {
		s.printErr(err)
	}

	return ff
}

func (s *server) printErr(err error) {
	s.erro.Println(err)
	s.windowLogMessage(protocol.MessageTypeError, err.Error())
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
