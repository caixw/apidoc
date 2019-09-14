// SPDX-License-Identifier: MIT

package input

import (
	"os"
	"path/filepath"

	"github.com/issue9/utils"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"

	"github.com/caixw/apidoc/v5/internal/lang"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
	opt "github.com/caixw/apidoc/v5/options"
)

type options struct {
	opt.Input
	blocks   []lang.Blocker    // 根据 Lang 生成
	paths    []string          // 根据 Dir 和 Recursive 生成
	encoding encoding.Encoding // 根据 Encoding 生成
}

func buildOptions(opt *opt.Input) (*options, *message.SyntaxError) {
	o := &options{}

	if len(opt.Dir) == 0 {
		return nil, message.NewError("", "dir", 0, locale.ErrRequired)
	}

	if !utils.FileExists(opt.Dir) {
		return nil, message.NewError("", "dir", 0, locale.ErrDirNotExists)
	}

	if len(opt.Lang) == 0 {
		return nil, message.NewError("", "dir", 0, locale.ErrRequired)
	}

	language := lang.Get(opt.Lang)
	if language == nil {
		return nil, message.NewError("", "dir", 0, locale.ErrUnsupportedInputLang, opt.Lang)
	}
	o.blocks = language.Blocks

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
		return nil, message.NewError("", "dir", 0, err.Error())
	}
	if len(paths) == 0 {
		return nil, message.NewError("", "dir", 0, locale.ErrDirIsEmpty)
	}
	o.paths = paths

	// 生成 encoding
	if opt.Encoding != "" {
		o.encoding, err = ianaindex.IANA.Encoding(opt.Encoding)
		if err != nil {
			return nil, message.WithError(err, "", "encoding", 0)
		}
	}

	o.Input = *opt
	return o, nil
}

// 按 Options 中的规则查找所有符合条件的文件列表。
func recursivePath(o *opt.Input) ([]string, error) {
	paths := []string{}

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
