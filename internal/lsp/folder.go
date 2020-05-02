// SPDX-License-Identifier: MIT

package lsp

import (
	"fmt"
	"path/filepath"
	"strings"

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
}

func (f *folder) close() error {
	f.h.Stop()
	return nil
}

// uri 是否与属于项目匹配
func (f *folder) matchURI(uri core.URI) bool {
	return strings.HasPrefix(string(uri), string(f.URI))
}

func (f *folder) matchPosition(uri core.URI, pos core.Position) (bool, error) {
	var block *core.Block
	if f.doc.Block.Location.URI == uri {
		block = f.doc.Block
	} else {
		for _, api := range f.doc.Apis {
			if api.Block.Location.URI == uri {
				block = api.Block
				break
			}
		}
	}
	if block == nil {
		return false, nil
	}
	return pos.Line >= block.Location.Range.Start.Line &&
		pos.Line <= block.Location.Range.End.Line, nil
}

func (f *folder) openFile(uri core.URI) error {
	file, err := uri.File()
	if err != nil {
		return err
	}

	var input *build.Input
	ext := filepath.Ext(file)
	for _, i := range f.cfg.Inputs {
		if inStringSlice(i.Exts, ext) {
			input = i
			break
		}
	}
	if input == nil {
		return fmt.Errorf("xxx")
	}

	build.ParseFile(f.doc, f.h, uri, input)

	return nil
}

func (f *folder) closeFile(uri core.URI) error {
	f.doc.DeleteFile(uri)
	return nil
}

func (s *server) appendFolders(folders ...protocol.WorkspaceFolder) (err error) {
	for _, f := range folders {
		h := core.NewMessageHandler(s.messageHandler)
		cfg := build.LoadConfig(h, f.URI)
		if cfg == nil {
			cfg, err = build.DetectConfig(f.URI, true)
			if err != nil {
				return err
			}
		}

		s.folders = append(s.folders, &folder{
			WorkspaceFolder: f,
			doc:             &ast.APIDoc{},
			h:               h,
			cfg:             cfg,
		})
	}

	return nil
}

func (s *server) messageHandler(msg *core.Message) {
	switch msg.Type {
	case core.Erro:
		// TODO
	case core.Warn:
		// TODO
	case core.Succ: // 仅处理错误和警告
	case core.Info: // 仅处理错误和警告
	default:
		panic("unreached")
	}
}

func (s *server) getMatchFolder(uri core.URI) *folder {
	for _, f := range s.folders {
		if f.matchURI(uri) {
			return f
		}
	}
	return nil
}

func inStringSlice(slice []string, key string) bool {
	for _, v := range slice {
		if v == key {
			return true
		}
	}
	return false
}
