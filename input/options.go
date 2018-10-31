// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"os"
	"path/filepath"

	"github.com/issue9/utils"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/internal/errors"
	"github.com/caixw/apidoc/internal/lang"
	"github.com/caixw/apidoc/internal/locale"
	"github.com/caixw/apidoc/options"
)

// Options 指定输入内容的相关信息。
type Options struct {
	options.Input

	blocks   []lang.Blocker    // 根据 Lang 生成
	paths    []string          // 根据 Dir 和 Recursive 生成
	encoding encoding.Encoding // 根据 Encoding 生成
}

func newError(field string, key message.Reference, args ...interface{}) *errors.Error {
	return &errors.Error{
		Field:       field,
		MessageKey:  key,
		MessageArgs: args,
	}
}

// Sanitize 检测 Options 变量是否符合要求
func (opt *Options) Sanitize() error {
	if len(opt.Dir) == 0 {
		return newError("dir", locale.ErrRequired)
	}

	if !utils.FileExists(opt.Dir) {
		return newError("dir", locale.ErrDirNotExists)
	}

	if len(opt.Lang) == 0 {
		return newError("dir", locale.ErrRequired)
	}

	language := lang.Get(opt.Lang)
	if language == nil {
		return newError("dir", locale.ErrUnsupportedInputLang, opt.Lang)
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
		return newError("dir", err.Error())
	}
	if len(paths) == 0 {
		return newError("dir", locale.ErrDirIsEmpty)
	}
	opt.paths = paths

	// 生成 encoding
	if opt.Encoding != "" {
		opt.encoding, err = ianaindex.IANA.Encoding(opt.Encoding)
		if err != nil {
			return locale.Errorf(locale.ErrUnsupportedEncoding, opt.Encoding)
		}
	}

	return nil
}

// 按 Options 中的规则查找所有符合条件的文件列表。
func recursivePath(o *Options) ([]string, error) {
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
