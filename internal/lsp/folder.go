// SPDX-License-Identifier: MIT

package lsp

import (
	"fmt"
	"os"
	"path/filepath"

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
	doc *ast.APIDoc
	h   *core.MessageHandler
	cfg *build.Config
	srv *server

	// 保存着错误和警告的信息
	errors, warns []*core.SyntaxError
}

func (f *folder) close() {
	f.srv.windowLogLogMessage(locale.CloseLSPFolder, f.Name)
	f.h.Stop()
}

func (f *folder) parseBlock(block core.Block) {
	var input *build.Input
	ext := filepath.Ext(block.Location.URI.String())
	for _, i := range f.cfg.Inputs {
		if sliceutil.Count(i.Exts, func(index int) bool { return i.Exts[index] == ext }) > 0 {
			input = i
			break
		}
	}
	if input == nil { // 无需解析
		return
	}

	f.doc.ParseBlocks(f.h, func(blocks chan core.Block) {
		input.Parse(blocks, f.h, block)
	})
}

func (f *folder) messageHandler(msg *core.Message) {
	switch msg.Type {
	case core.Erro:
		err, ok := msg.Message.(*core.SyntaxError)
		if !ok {
			f.srv.erro.Println(fmt.Sprintf("获得了非 core.SyntaxError 错误 %#v", msg.Message))
		}
		f.errors = append(f.errors, err)
	case core.Warn:
		err, ok := msg.Message.(*core.SyntaxError)
		if !ok {
			f.srv.erro.Println(fmt.Sprintf("获得了非 core.SyntaxError 错误 %#v", msg.Message))
		}
		f.warns = append(f.warns, err)
	case core.Succ, core.Info: // 仅处理错误和警告
	default:
		panic("unreached")
	}

	f.srv.textDocumentPublishDiagnostics(f.URI, f.errors, f.warns)
}

func (s *server) appendFolders(folders ...protocol.WorkspaceFolder) (err error) {
	for _, f := range folders {
		ff := &folder{
			WorkspaceFolder: f,
			doc:             &ast.APIDoc{},
			srv:             s,
		}

		ff.h = core.NewMessageHandler(ff.messageHandler)
		ff.cfg, err = build.LoadConfig(f.URI)
		if os.IsNotExist(err) {
			if ff.cfg, err = build.DetectConfig(f.URI, true); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		ff.doc.ParseBlocks(ff.h, func(blocks chan core.Block) {
			build.ParseInputs(blocks, ff.h, ff.cfg.Inputs...)
		})

		s.folders = append(s.folders, ff)
	}

	return nil
}
