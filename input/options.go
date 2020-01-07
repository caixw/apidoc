// SPDX-License-Identifier: MIT

package input

import (
	"os"
	"path/filepath"

	"github.com/issue9/utils"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"

	"github.com/caixw/apidoc/v6/internal/lang"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/message"
)

// Options 指定输入内容的相关信息。
type Options struct {
	// 输入的目标语言
	//
	// 取值为 lang.Language.Name
	Lang string `yaml:"lang"`

	// 源代码目录
	Dir string `yaml:"dir"`

	// 需要扫描的文件扩展名
	//
	// 若未指定，则根据 Lang 选项获取其默认的扩展名作为过滤条件。
	Exts []string `yaml:"exts,omitempty"`

	// 是否查找 Dir 的子目录
	Recursive bool `yaml:"recursive"`

	// 源文件的编码，默认为 UTF-8
	Encoding string `yaml:"encoding,omitempty"`

	blocks   []lang.Blocker    // 根据 Lang 生成
	paths    []string          // 根据 Dir、Exts 和 Recursive 生成
	encoding encoding.Encoding // 根据 Encoding 生成
}

func (opt *Options) sanitize() *message.SyntaxError {
	if opt == nil {
		return message.NewLocaleError("", "", 0, locale.ErrRequired)
	}

	if len(opt.Dir) == 0 {
		return message.NewLocaleError("", "dir", 0, locale.ErrRequired)
	}

	if !utils.FileExists(opt.Dir) {
		return message.NewLocaleError("", "dir", 0, locale.ErrDirNotExists)
	}

	if len(opt.Lang) == 0 {
		return message.NewLocaleError("", "dir", 0, locale.ErrRequired)
	}

	language := lang.Get(opt.Lang)
	if language == nil {
		return message.NewLocaleError("", "dir", 0, locale.ErrUnsupportedInputLang, opt.Lang)
	}
	opt.blocks = language.Blocks

	if len(opt.Exts) > 0 {
		exts := make([]string, 0, len(opt.Exts))
		for _, ext := range opt.Exts {
			if len(ext) == 0 {
				continue
			}

			if ext[0] != '.' {
				ext = "." + ext
			}
			exts = append(exts, ext)
		}
		opt.Exts = exts
	} else {
		opt.Exts = language.Exts
	}

	// 生成 paths
	paths, err := recursivePath(opt)
	if err != nil {
		return message.WithError("", "dir", 0, err)
	}
	if len(paths) == 0 {
		return message.NewLocaleError("", "dir", 0, locale.ErrDirIsEmpty)
	}
	opt.paths = paths

	// 生成 encoding
	if opt.Encoding != "" {
		opt.encoding, err = ianaindex.IANA.Encoding(opt.Encoding)
		if err != nil {
			return message.WithError("", "encoding", 0, err)
		}
	}

	return nil
}

// 按 Options 中的规则查找所有符合条件的文件列表。
func recursivePath(o *Options) ([]string, error) {
	var paths []string

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
		if fi.IsDir() && !o.Recursive && path != o.Dir {
			return filepath.SkipDir
		} else if extIsEnabled(filepath.Ext(path)) {
			paths = append(paths, path)
		}
		return nil
	}

	if err := filepath.Walk(o.Dir, walk); err != nil {
		return nil, err
	}

	return paths, nil
}
