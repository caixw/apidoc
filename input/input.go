// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 用于处理文件输入。
//
// 多行注释和单行注释在处理上会有一定区别：
//
// 单行注释，风格相同且相邻的注释会被合并成一个注释块。
// 单行注释，风格不相同且相邻的注释会被按注释风格多个注释块。
// 而多行注释，即使两个注释释块相邻也会被分成两个注释块来处理。
package input

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/logs"
)

type Options struct {
	Lang      string   `json:"lang"`           // 输入的目标语言
	Dir       string   `json:"dir"`            // 源代码目录
	Exts      []string `json:"exts,omitempty"` // 需要扫描的文件扩展名，若未指定，则使用默认值
	Recursive bool     `json:"recursive"`      // 是否查找Dir的子目录
}

// 分析源代码，获取相应的文档内容。
func Parse(o *Options) (*doc.Doc, error) {
	b, found := langs[o.Lang]
	if !found {
		return nil, errors.New("不支持该语言")
	}

	paths, err := recursivePath(o)
	if err != nil {
		return nil, err
	}

	docs := doc.New()
	wg := sync.WaitGroup{}
	defer wg.Wait()
	for _, path := range paths {
		wg.Add(1)
		go func() {
			parseFile(docs, path, b)
			wg.Done()
		}()
	}

	return docs, nil
}

// 分析 path 指向的文件，并将内容写入到 docs 中。
func parseFile(docs *doc.Doc, path string, blocks []*block) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		synerr := &doc.SyntaxError{Message: err.Error()}
		printSyntaxError(synerr)
		return
	}

	l := &lexer{data: data}
	var block *block

	wg := sync.WaitGroup{}
	defer wg.Wait()

LOOP:
	for {
		switch {
		case l.atEOF():
			return
		case block == nil:
			block = l.block(blocks)
		case block != nil:
			rs, err := block.end(l)
			if err != nil {
				err.File = path
				printSyntaxError(err)
				return
			}

			if block.Type == blockTypeString {
				block = nil
				continue LOOP
			}
			block = nil

			wg.Add(1)
			go func() {
				if err = docs.Scan(rs); err != nil {
					err.Line += l.lineNumber()
					err.File = path
					printSyntaxError(err)
				}
				wg.Done()
			}()
		} // end switch
	} // end for
}

func printSyntaxError(err *doc.SyntaxError) {
	logs.Error("[语法错误] ", err)
}

// 根据recursive值确定是否递归查找paths每个目录下的子目录。
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

// 检测 Options 变量是否符合要求
func (opt *Options) Init() error {
	if len(opt.Dir) == 0 {
		return errors.New("未指定源码目录")
	}

	if len(opt.Lang) == 0 {
		return errors.New("必须指定参数 type")
	}

	if langIsSupported(opt.Lang) {
		return fmt.Errorf("暂不支持该类型[%v]的语言", opt.Lang)
	}

	opt.Dir += string(os.PathSeparator)

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
