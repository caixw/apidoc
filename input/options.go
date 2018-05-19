// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/issue9/utils"

	"github.com/caixw/apidoc/input/encoding"
	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/types/openapi"
)

// Options 指定输入内容的相关信息。
type Options struct {
	ErrorLog *log.Logger `yaml:"-"` // 错误输出通道
	WarnLog  *log.Logger `yaml:"-"` // 警告输出通道

	StartLineNumber int      `yaml:"startLineNumber,omitempty"` // 代码的超始行号，默认为 0
	Lang            string   `yaml:"lang"`                      // 输入的目标语言
	Dir             string   `yaml:"dir"`                       // 源代码目录
	Exts            []string `yaml:"exts,omitempty"`            // 需要扫描的文件扩展名，若未指定，则使用默认值
	Recursive       bool     `yaml:"recursive"`                 // 是否查找 Dir 的子目录
	Encoding        string   `yaml:"encoding,omitempty"`        // 文件的编码

	blocks []blocker
	paths  []string
}

// Sanitize 检测 Options 变量是否符合要求
func (opt *Options) Sanitize() *openapi.Error {
	if len(opt.Dir) == 0 {
		return &openapi.Error{Field: "dir", Message: locale.Sprintf(locale.ErrRequired)}
	}

	if !utils.FileExists(opt.Dir) {
		return &openapi.Error{Field: "dir", Message: locale.Sprintf(locale.ErrDirNotExists)}
	}

	if len(opt.Lang) == 0 {
		return &openapi.Error{Field: "lang", Message: locale.Sprintf(locale.ErrRequired)}
	}
	blocks, found := langs[opt.Lang]
	if !found {
		return &openapi.Error{Field: "lang", Message: locale.Sprintf(locale.ErrUnsupportedInputLang, opt.Lang)}
	}
	opt.blocks = blocks

	if !langIsSupported(opt.Lang) {
		return &openapi.Error{Field: "lang", Message: locale.Sprintf(locale.ErrUnsupportedInputLang, opt.Lang)}
	}

	if len(opt.Encoding) == 0 {
		opt.Encoding = encoding.DefaultEncoding
	}

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
		opt.Exts = langExts[opt.Lang]
	}

	paths, err := recursivePath(opt)
	if err != nil {
		return &openapi.Error{Field: "dir", Message: err.Error()}
	}
	if len(paths) == 0 {
		return &openapi.Error{Field: "dir", Message: locale.Sprintf(locale.ErrDirIsEmpty)}
	}
	opt.paths = paths

	return nil
}

// Detect 检测指定目录下的内容，并为其生成一个合适的 Options 实例。
//
// 检测依据为根据扩展名来做统计，数量最大且被支持的获胜。
func Detect(dir string, recursive bool) (*Options, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	exts, err := detectExts(dir, recursive)
	if err != nil {
		return nil, err
	}

	// 删除不支持的扩展名
	for ext := range exts {
		lang := getLangByExt(ext)
		if len(lang) <= 0 {
			delete(exts, ext)
		}
	}

	if len(exts) == 0 {
		return nil, errors.New(locale.Sprintf(locale.ErrNotFoundSupportedLang))
	}

	ext := ""
	cnt := 0
	for k, v := range exts {
		if v >= cnt {
			ext = k
			cnt = v
		}
	}
	if len(ext) == 0 {
		return nil, errors.New(locale.Sprintf(locale.ErrNotFoundSupportedLang))
	}

	lang := getLangByExt(ext)
	if len(lang) == 0 {
		return nil, errors.New(locale.Sprintf(locale.ErrNotFoundSupportedLang))
	}

	return &Options{
		StartLineNumber: 0,
		Lang:            lang,
		Dir:             dir,
		Exts:            langExts[lang],
		Recursive:       recursive,
	}, nil
}

// 返回 dir 目录下文件类型及对应的文件数量的一个集合。
// recursive 表示是否查找子目录。
func detectExts(dir string, recursive bool) (map[string]int, error) {
	exts := map[string]int{}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if !recursive && dir != path {
				return filepath.SkipDir
			}
		} else {
			ext := strings.ToLower(filepath.Ext(path))
			exts[ext]++
		}

		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return nil, err
	}

	return exts, nil
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
