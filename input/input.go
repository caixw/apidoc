// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package input 用于处理文件输入，过滤代码，生成 doc.Doc 数据。
//
// 多行注释和单行注释在处理上会有一定区别：
//  - 单行注释，风格相同且相邻的注释会被合并成一个注释块；
//  - 单行注释，风格不相同且相邻的注释会被按注释风格多个注释块；
//  - 多行注释，即使两个注释释块相邻也会被分成两个注释块来处理。
package input

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	b "github.com/caixw/apidoc/input/block"
	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/types"

	"github.com/issue9/utils"
)

// 需要解析的最小代码块，小于此值，将不作解析
const miniSize = len("@api ")

// Options 指定输入内容的相关信息。
type Options struct {
	SyntaxLog       *log.Logger `yaml:"-"`                         // 语法错误输出通道
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

// Parse 分析源代码，获取相应的文档内容。
func Parse(docs *types.Doc, o *Options) error {
	blocks, found := langs[o.Lang]
	if !found {
		return errors.New(locale.Sprintf(locale.ErrUnsupportedInputLang, o.Lang))
	}

	paths, err := recursivePath(o)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	defer wg.Wait()
	for _, path := range paths {
		wg.Add(1)
		go func(path string) {
			parseFile(docs, path, blocks, o.SyntaxLog, o.StartLineNumber)
			wg.Done()
		}(path)
	}

	return nil
}

// 分析 path 指向的文件，并将内容写入到 docs 中。
func parseFile(docs *types.Doc, path string, blocks []blocker, synerrLog *log.Logger, startLine int) {
	data, err := ioutil.ReadFile(path)
	if err != nil && synerrLog != nil {
		synerrLog.Println(&types.SyntaxError{Message: err.Error(), File: path})
		return
	}

	l := &lexer{data: data}
	var block blocker

	wg := sync.WaitGroup{}
	defer wg.Wait()

	for {
		if l.atEOF() {
			return
		}

		if block == nil {
			block = l.block(blocks)
			if block == nil { // 没有匹配的 block 了
				return
			}
		}

		ln := l.lineNumber() + startLine // 记录当前的行号，顺便调整起始行号
		rs, ok := block.EndFunc(l)
		if !ok && synerrLog != nil {
			synerrLog.Println(&types.SyntaxError{Line: ln, File: path, Message: locale.Sprintf(locale.ErrNotFoundEndFlag)})
			return
		}

		block = nil
		if len(rs) < miniSize {
			continue
		}

		wg.Add(1)
		go func(rs []rune, ln int) {
			if err := b.Scan(docs, rs); err != nil && synerrLog != nil {
				err.Line += ln
				err.File = path
				synerrLog.Println(err)
			}
			wg.Done()
		}(rs, ln)
	} // end for
}

// 根据 recursive 值确定是否递归查找 paths 每个目录下的子目录。
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
