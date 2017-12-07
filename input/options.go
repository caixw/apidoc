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

	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/types"
	"github.com/issue9/utils"
)

// Options 指定输入内容的相关信息。
type Options struct {
	SyntaxErrorLog  *log.Logger `yaml:"-"`                         // 语法错误输出通道
	StartLineNumber int         `yaml:"startLineNumber,omitempty"` // 代码的超始行号，默认为 0
	Lang            string      `yaml:"lang"`                      // 输入的目标语言
	Dir             string      `yaml:"dir"`                       // 源代码目录
	Exts            []string    `yaml:"exts,omitempty"`            // 需要扫描的文件扩展名，若未指定，则使用默认值
	Recursive       bool        `yaml:"recursive"`                 // 是否查找 Dir 的子目录
}

// Sanitize 检测 Options 变量是否符合要求
func (opt *Options) Sanitize() *types.OptionsError {
	if len(opt.Dir) == 0 {
		return &types.OptionsError{Field: "dir", Message: locale.Sprintf(locale.ErrRequired)}
	}

	if !utils.FileExists(opt.Dir) {
		return &types.OptionsError{Field: "dir", Message: locale.Sprintf(locale.ErrDirNotExists)}
	}

	if len(opt.Lang) == 0 {
		return &types.OptionsError{Field: "lang", Message: locale.Sprintf(locale.ErrRequired)}
	}

	if !langIsSupported(opt.Lang) {
		return &types.OptionsError{Field: "lang", Message: locale.Sprintf(locale.ErrUnsupportedInputLang, opt.Lang)}
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

	return nil
}

// Detect 检测指定目录下的内容，并为其生成一个合适的 Options 实例。
//
// 检测依据为根据扩展名来做统计，数量最大且被支持的获胜。
// 不会分析子目录。
func Detect(dir string, recursive bool) (*Options, error) {
	paths, err := recursiveDir(dir, recursive)
	if err != nil {
		return nil, err
	}

	// langsMap 记录每个支持的语言对应的文件数量
	langsMap := make(map[string]int, len(paths))
	for _, f := range paths { // 遍历所有的文件
		ext := strings.ToLower(filepath.Ext(f))
		lang := getLangByExt(ext)
		if len(lang) > 0 {
			langsMap[lang]++
		}
	}

	if len(langsMap) == 0 {
		return nil, errors.New(locale.Sprintf(locale.ErrNotFoundSupportedLang))
	}

	lang := ""
	cnt := 0
	for k, v := range langsMap {
		if v >= cnt {
			lang = k
			cnt = v
		}
	}

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

// 返回 dir 目录下所有的文件集合。
// recursive 表示是否查找子目录下的文件。
func recursiveDir(dir string, recursive bool) ([]string, error) {
	paths := []string{}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if !recursive && dir != path {
				return filepath.SkipDir
			}
		} else {
			paths = append(paths, path)
		}

		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return nil, err
	}

	return paths, nil
}
