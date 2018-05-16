// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package input 用于处理文件输入，过滤代码，生成 types.Doc 数据。
//
// 多行注释和单行注释在处理上会有一定区别：
//  - 单行注释，风格相同且相邻的注释会被合并成一个注释块；
//  - 单行注释，风格不相同且相邻的注释会被按注释风格多个注释块；
//  - 多行注释，即使两个注释释块相邻也会被分成两个注释块来处理。
package input

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/caixw/apidoc/input/encoding"
	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/types"
	"github.com/caixw/apidoc/vars"
	yaml "gopkg.in/yaml.v2"
)

// 需要解析的最小代码块，小于此值，将不作解析
// 即其长度必须大于 @api 这四个字符串的长度
const miniSize = len(vars.API) + 1

// Parse 分析源代码，获取相应的文档内容。
func Parse(options ...*Options) (*types.Docs, time.Duration) {
	start := time.Now()
	docs := &types.Docs{
		Docs: make(map[string]*types.Doc, 5),
	}

	wg := &sync.WaitGroup{}
	for _, o := range options {
		wg.Add(1)
		go func(o *Options) {
			if err := parse(docs, o); err != nil {
				o.ErrorLog.Println(err)
			}
			wg.Done()
		}(o)
	}
	wg.Wait()

	return docs, time.Now().Sub(start)
}

func parse(docs *types.Docs, o *Options) error {
	blocks, found := langs[o.Lang]
	if !found {
		return errors.New(locale.Sprintf(locale.ErrUnsupportedInputLang, o.Lang))
	}

	paths, err := recursivePath(o)
	if err != nil {
		return errors.New(locale.Sprintf(locale.ErrUnsupportedInputLang, o.Lang))
	}

	wg := sync.WaitGroup{}
	defer wg.Wait()
	for _, path := range paths {
		wg.Add(1)
		go func(path string) {
			parseFile(docs, path, blocks, o)
			wg.Done()
		}(path)
	}

	return nil
}

// Encodings 返回支持的编码方式
func Encodings() []string {
	return encoding.Encodings()
}

// 分析 path 指向的文件，并将内容写入到 docs 中。
func parseFile(docs *types.Docs, path string, blocks []blocker, o *Options) {
	data, err := encoding.Transform(path, o.Encoding)
	if err != nil {
		if o.ErrorLog != nil {
			o.ErrorLog.Println(err)
		}
		return
	}

	l := &lexer{data: data, blocks: blocks}
	var block blocker

	wg := sync.WaitGroup{}
	defer wg.Wait()

	for {
		if l.atEOF() {
			return
		}

		if block == nil {
			block = l.block()
			if block == nil { // 没有匹配的 block 了
				return
			}
		}

		ln := l.lineNumber() + o.StartLineNumber // 记录当前的行号，顺便调整起始行号
		rs, ok := block.EndFunc(l)
		if !ok {
			// syntax.OutputError(o.ErrorLog, path, ln, locale.ErrNotFoundEndFlag)
			return // 没有找到结束标签，那肯定是到文件尾了，可以直接返回。
		}

		block = nil
		if len(rs) < miniSize {
			continue
		}

		wg.Add(1)
		go func(rs []rune, ln int) {
			/*i := &syntax.Input{
				File:  path,
				Line:  ln,
				Data:  rs,
				Error: o.ErrorLog,
				Warn:  o.WarnLog,
			}
			syntax.Parse(i, docs)*/
			parseData([]byte(string(rs)), docs) // TODO 减少转换

			wg.Done()
		}(rs, ln)
	} // end for
}

func parseData(data []byte, d *types.Docs) error {
	data = bytes.TrimLeft(data, " ")

	if bytes.HasPrefix([]byte(vars.API), data) {
		index := bytes.IndexByte(data, '\n')
		line := data[:index]
		data = data[index+1:]
		api := &types.API{}
		if err := yaml.Unmarshal(data, api); err != nil {
			return err
		}

		api.API = string(line)
		return d.NewAPI(api)
	}

	if bytes.HasPrefix([]byte(vars.APIDoc), data) {
		index := bytes.IndexByte(data, '\n')
		line := data[:index]
		data = data[index+1:]
		info := &types.Info{}
		if err := yaml.Unmarshal(data, info); err != nil {
			return err
		}

		info.Title = string(line)
		return d.NewInfo(info)
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
