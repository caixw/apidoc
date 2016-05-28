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
	"io/ioutil"
	"sync"

	"github.com/caixw/apidoc/doc"
	"github.com/issue9/term/colors"
)

func Parse(paths []string, langID string) (*doc.Doc, error) {
	docs := doc.New()

	b, found := langs[langID]
	if !found {
		return nil, errors.New("不支持该语言")
	}

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
	colors.Print(colors.Stderr, colors.Red, colors.Default, "SyntaxError:")
	colors.Println(colors.Stderr, colors.Default, colors.Default, err)
}
