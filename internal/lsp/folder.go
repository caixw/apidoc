// SPDX-License-Identifier: MIT

package lsp

import (
	"fmt"
	"path/filepath"

	"github.com/issue9/sliceutil"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
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

func (f *folder) close() error {
	f.h.Stop()
	return nil
}

func (f *folder) openFile(uri core.URI) error {
	file, err := uri.File()
	if err != nil {
		return err
	}

	var input *build.Input
	ext := filepath.Ext(file)
	for _, i := range f.cfg.Inputs {
		if sliceutil.Count(i.Exts, func(index int) bool { return i.Exts[index] == ext }) > 0 {
			input = i
			break
		}
	}
	if input == nil { // 无需解析
		return nil
	}

	f.parseFile(uri, input)
	return nil
}

// 分析 path 的内容，并将其中的文档解析至 doc
func (f *folder) parseFile(uri core.URI, i *build.Input) {
	f.doc.ParseBlocks(f.h, func(blocks chan core.Block) {
		i.ParseFile(blocks, f.h, uri)
	})
}

func (f *folder) closeFile(uri core.URI) error {
	f.doc.DeleteURI(uri)
	return nil
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
		ff.cfg = build.LoadConfig(ff.h, f.URI)
		if ff.cfg == nil {
			if ff.cfg, err = build.DetectConfig(f.URI, true); err != nil {
				return err
			}
		}

		ff.doc.ParseBlocks(ff.h, func(blocks chan core.Block) {
			build.ParseInputs(blocks, ff.h, ff.cfg.Inputs...)
		})

		s.folders = append(s.folders, ff)
	}

	return nil
}
