// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"log"

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
