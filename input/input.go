// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package input 用于处理输入的文件，从代码中提取基本的注释内容。
//
// 多行注释和单行注释在处理上会有一定区别：
//  - 单行注释，风格相同且相邻的注释会被合并成一个注释块；
//  - 单行注释，风格不相同且相邻的注释会被按注释风格多个注释块；
//  - 多行注释，即使两个注释释块相邻也会被分成两个注释块来处理。
package input

import (
	"bytes"
	"io/ioutil"
	"log"
	"sync"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"

	"github.com/caixw/apidoc/internal/lang"
)

// Block 解析出来的注释块
type Block struct {
	File string
	Line int
	Data []byte
}

// Parse 分析源代码，获取相应的文档内容。
//
// 当所有的代码块已经放入 Block 之后，Block 会被关闭。
//
// errlog 表示错误内容输出通道，若不需要，则使用 nil 代替。
func Parse(errlog *log.Logger, o ...*Options) chan Block {
	data := make(chan Block, 500)

	go func() {
		wg := &sync.WaitGroup{}
		for _, opt := range o {
			parse(data, errlog, wg, opt)
		}
		wg.Wait()

		close(data)
	}()

	return data
}

func parse(data chan Block, errlog *log.Logger, wg *sync.WaitGroup, o *Options) {
	for _, path := range o.paths {
		wg.Add(1)
		go func(path string) {
			parseFile(data, errlog, path, o)
			wg.Done()
		}(path)
	}
}

// 分析 path 指向的文件。
//
// NOTE: parseFile 内部不能有协程处理代码。
func parseFile(channel chan Block, errlog *log.Logger, path string, o *Options) {
	data, err := readFile(path, o.encoding)
	if err != nil {
		if errlog != nil {
			errlog.Println(err)
		}
		return
	}

	ret := lang.Parse(errlog, data, o.blocks)
	for line, data := range ret {
		channel <- Block{
			File: path,
			Line: line,
			Data: data,
		}
	}
}

// 以指定的编码方式读取内容。
func readFile(path string, encoding encoding.Encoding) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if encoding == nil {
		return data, nil
	}

	reader := transform.NewReader(bytes.NewReader(data), encoding.NewDecoder())
	return ioutil.ReadAll(reader)
}
