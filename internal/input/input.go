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
	"context"
	"io/ioutil"
	"log"
	"sync"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"

	"github.com/caixw/apidoc/internal/errors"
	"github.com/caixw/apidoc/internal/lang"
	opt "github.com/caixw/apidoc/options"
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
// 只有处理文本内容的错误信息会被输出到 errolog 和 warnlog，
// 其中 errlog 用于错误信息，而 warnlog 表示的一些警告信息。
// 普通错误依然通过返回值返回。
func Parse(ctx context.Context, errlog, warnlog *log.Logger, o ...*opt.Input) (chan Block, error) {
	if len(o) == 0 {
		return nil, &errors.Error{
			// TODO
		}
	}

	opts := make([]*options, 0, len(o))
	for _, item := range o {
		opt, err := buildOptions(item)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	data := make(chan Block, 500)

	go func() {
		wg := &sync.WaitGroup{}
		for _, opt := range opts {
			select {
			case <-ctx.Done():
				return
			default:
				parse(ctx, data, errlog, warnlog, wg, opt)
			}
		}
		wg.Wait()

		close(data)
	}()

	return data, nil
}

func parse(ctx context.Context, data chan Block, errlog, warnlog *log.Logger, wg *sync.WaitGroup, o *options) {
	for _, path := range o.paths {
		select {
		case <-ctx.Done():
			return
		default:
			wg.Add(1)
			go func(path string) {
				parseFile(data, errlog, warnlog, path, o)
				wg.Done()
			}(path)
		}
	} // end for
}

// 分析 path 指向的文件。
//
// NOTE: parseFile 内部不能有协程处理代码。
func parseFile(channel chan Block, errlog, warnlog *log.Logger, path string, o *options) {
	data, err := readFile(path, o.encoding)
	if err != nil {
		if errlog != nil {
			errlog.Println(err)
		}
		return
	}

	ret, err := lang.Parse(data, o.blocks)
	if err != nil {
		if serr, ok := err.(*errors.Error); ok {
			serr.File = path
		}
		errlog.Println(err)
	}
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
