// SPDX-License-Identifier: MIT

package lsp

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/caixw/apidoc/v6/build"
	"github.com/caixw/apidoc/v6/internal/lsp/protocol"
	"github.com/caixw/apidoc/v6/message"
	"github.com/caixw/apidoc/v6/spec"
)

// 表示项目文件夹
type folder struct {
	protocol.WorkspaceFolder
	doc *spec.APIDoc
	h   *message.Handler
	cfg *build.Config
}

func (f *folder) close() error {
	f.h.Stop()
	return nil
}

// uri 是否与属于项目匹配
func (f *folder) matchURI(uri protocol.DocumentURI) bool {
	return strings.HasPrefix(string(uri), string(f.URI))
}

func (f *folder) matchPosition(uri protocol.DocumentURI, pos protocol.Position) (bool, error) {
	file, err := uri.File()
	if err != nil {
		return false, err
	}

	var block *spec.Block
	if f.doc.Block.File == file {
		block = f.doc.Block
	} else {
		for _, api := range f.doc.Apis {
			if api.Block.File == file {
				block = api.Block
				break
			}
		}
	}
	if block == nil {
		return false, nil
	}
	return pos.Line >= block.Line && pos.Line <= bytes.Count(block.Data, []byte{'\n'})+block.Line, nil
}

func (f *folder) openFile(uri protocol.DocumentURI) error {
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

	build.ParseFile(f.doc, f.h, file, input)

	return nil
}

func (f *folder) closeFile(uri protocol.DocumentURI) error {
	file, err := uri.File()
	if err != nil {
		return err
	}

	f.doc.DeleteFile(file)
	return nil
}

func (s *server) appendFolders(folders ...protocol.WorkspaceFolder) error {
	for _, f := range folders {
		h := message.NewHandler(s.messageHandler)
		file, err := f.URI.File()
		if err != nil {
			return err
		}

		cfg := build.LoadConfig(h, file)
		if cfg == nil {
			cfg, err = build.DetectConfig(file, true)
			if err != nil {
				return err
			}
		}

		s.folders = append(s.folders, &folder{
			WorkspaceFolder: f,
			doc:             spec.New(),
			h:               h,
			cfg:             cfg,
		})
	}

	return nil
}

func (s *server) messageHandler(msg *message.Message) {
	switch msg.Type {
	case message.Erro:
		// TODO
	case message.Warn:
		// TODO
	case message.Succ: // 仅处理错误和警告
	case message.Info: // 仅处理错误和警告
	default:
		panic("unreached")
	}
}

func (s *server) getMatchFolder(uri protocol.DocumentURI) *folder {
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
