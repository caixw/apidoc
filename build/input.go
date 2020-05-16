// SPDX-License-Identifier: MIT

package build

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/issue9/utils"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/lang"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// Input 指定输入内容的相关信息。
type Input struct {
	// 输入的目标语言
	//
	// 取值为 lang.Language.Name
	Lang string `yaml:"lang"`

	// 源代码目录
	Dir core.URI `yaml:"dir"`

	// 需要扫描的文件扩展名
	//
	// 若未指定，则根据 Lang 选项获取其默认的扩展名作为过滤条件。
	Exts []string `yaml:"exts,omitempty"`

	// 是否查找 Dir 的子目录
	Recursive bool `yaml:"recursive"`

	// 源文件的编码，默认为 UTF-8
	Encoding string `yaml:"encoding,omitempty"`

	blocks   []lang.Blocker    // 根据 Lang 生成
	paths    []core.URI        // 根据 Dir、Exts 和 Recursive 生成
	encoding encoding.Encoding // 根据 Encoding 生成
}

// Sanitize 验证参数正确性
func (o *Input) Sanitize() error {
	if o == nil {
		return core.NewSyntaxError(core.Location{}, "", locale.ErrRequired)
	}

	if len(o.Dir) == 0 {
		return core.NewSyntaxError(core.Location{}, "dir", locale.ErrRequired)
	}

	local, err := o.Dir.File()
	if err != nil {
		return core.NewSyntaxErrorWithError(core.Location{}, "dir", err)
	}

	if !utils.FileExists(local) {
		return core.NewSyntaxError(core.Location{}, "dir", locale.ErrDirNotExists)
	}

	if len(o.Lang) == 0 {
		return core.NewSyntaxError(core.Location{}, "lang", locale.ErrRequired)
	}

	language := lang.Get(o.Lang)
	if language == nil {
		return core.NewSyntaxError(core.Location{}, "lang", locale.ErrInvalidValue)
	}
	o.blocks = language.Blocks

	if len(o.Exts) > 0 {
		exts := make([]string, 0, len(o.Exts))
		for _, ext := range o.Exts {
			if len(ext) == 0 {
				continue
			}

			if ext[0] != '.' {
				ext = "." + ext
			}
			exts = append(exts, ext)
		}
		o.Exts = exts
	} else {
		o.Exts = language.Exts
	}

	// 生成 paths
	paths, err := recursivePath(o)
	if err != nil {
		return core.NewSyntaxErrorWithError(core.Location{}, "dir", err)
	}
	if len(paths) == 0 {
		return core.NewSyntaxError(core.Location{}, "dir", locale.ErrDirIsEmpty)
	}
	o.paths = paths

	// 生成 encoding
	if o.Encoding != "" {
		o.encoding, err = ianaindex.IANA.Encoding(o.Encoding)
		if err != nil {
			return core.NewSyntaxErrorWithError(core.Location{}, "encoding", err)
		}
	}

	return nil
}

// 按 Input 中的规则查找所有符合条件的文件列表。
func recursivePath(o *Input) ([]core.URI, error) {
	var uris []core.URI

	local, err := o.Dir.File()
	if err != nil {
		return nil, err
	}

	extIsEnabled := func(ext string) bool {
		for _, v := range o.Exts {
			if ext == v {
				return true
			}
		}
		return false
	}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() && !o.Recursive && path != local {
			return filepath.SkipDir
		} else if extIsEnabled(filepath.Ext(path)) {
			uris = append(uris, core.FileURI(path))
		}
		return nil
	}

	if err := filepath.Walk(local, walk); err != nil {
		return nil, err
	}

	return uris, nil
}

// 分析 opt 中所指定的内容
//
// 分析后的内容推送至 blocks 中。
func parseInputs(blocks chan core.Block, h *core.MessageHandler, opt ...*Input) {
	wg := &sync.WaitGroup{}
	for _, o := range opt {
		for _, path := range o.paths {
			wg.Add(1)
			go func(path core.URI, o *Input) {
				o.ParseFile(blocks, h, path)
				wg.Done()
			}(path, o)
		}
	}
	wg.Wait()
}

// ParseFile 分析 uri 指向的文件
func (o *Input) ParseFile(blocks chan core.Block, h *core.MessageHandler, uri core.URI) {
	data, err := uri.ReadAll(o.encoding)
	if err != nil {
		h.Error(core.Erro, core.NewSyntaxErrorWithError(core.Location{URI: uri}, "", err))
		return
	}

	l, err := lang.NewLexer(data, o.blocks)
	if err != nil {
		h.Error(core.Erro, core.NewSyntaxErrorWithError(core.Location{URI: uri}, "", err))
		return
	}

	l.Parse(blocks, h, uri)
}
