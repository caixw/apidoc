// SPDX-License-Identifier: MIT

// Package input 用于处理输入的文件，从代码中提取基本的注释内容。
//
// 多行注释和单行注释在处理上会有一定区别：
//  - 单行注释，风格相同且相邻的注释会被合并成一个注释块；
//  - 单行注释，风格不相同且相邻的注释会被按注释风格多个注释块；
//  - 多行注释，即使两个注释释块相邻也会被分成两个注释块来处理。
package input

import (
	"io/ioutil"
	"os"
	"sync"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"

	"github.com/caixw/apidoc/v6/internal/lang"
	"github.com/caixw/apidoc/v6/message"
)

// Block 表示原始的注释代码块
type Block struct {
	File string
	Line int
	Data []byte
}

// Parse 分析 opt 中所指定的内容
//
// 分析后的内容推送至 blocks 中。
func Parse(blocks chan Block, h *message.Handler, opt ...*Options) {
	wg := &sync.WaitGroup{}
	for _, o := range opt {
		for _, path := range o.paths {
			wg.Add(1)
			go func(path string) {
				ParseFile(blocks, h, path, o)
				wg.Done()
			}(path)
		}
	}
	wg.Wait()
}

// ParseFile 分析 path 指向的文件。
func ParseFile(blocks chan Block, h *message.Handler, path string, o *Options) {
	data, err := readFile(path, o.encoding)
	if err != nil {
		h.Error(message.Erro, message.WithError(path, "", 0, err))
		return
	}

	ret := lang.Parse(path, data, o.blocks, h)
	for line, data := range ret {
		blocks <- Block{
			File: path,
			Line: line,
			Data: data,
		}
	}
}

// 以指定的编码方式读取内容。
func readFile(path string, enc encoding.Encoding) ([]byte, error) {
	if enc == nil || enc == encoding.Nop {
		return ioutil.ReadFile(path)
	}

	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	reader := transform.NewReader(r, enc.NewDecoder())
	return ioutil.ReadAll(reader)
}
