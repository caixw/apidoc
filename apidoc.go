// SPDX-License-Identifier: MIT

// Package apidoc RESTful API 文档生成工具。
package apidoc

import (
	"bytes"
	"context"
	"sync"

	"github.com/caixw/apidoc/v5/doc"
	i "github.com/caixw/apidoc/v5/input"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
	o "github.com/caixw/apidoc/v5/output"
)

// Version 获取当前程序的版本号
func Version() string {
	return vars.Version()
}

// Output 按 output 的要求输出内容。
func Output(doc *doc.Doc, output *o.Options) error {
	return o.Render(doc, output)
}

// Parse 分析从 input 中获取的代码块
//
// 所有与解析有关的错误均通过 h 输出。
// 如果 input 参数有误，会通过 error 参数返回。
func Parse(ctx context.Context, h *message.Handler, input ...*i.Options) (*doc.Doc, error) {
	block, err := i.Parse(ctx, h, input...)
	if err != nil {
		return nil, err
	}

	doc := doc.New()
	wg := sync.WaitGroup{}

LOOP:
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case blk, ok := <-block:
			if !ok {
				break LOOP
			}

			wg.Add(1)
			go func(b i.Block) {
				parseBlock(doc, b, h)
				wg.Done()
			}(blk)
		}
	}

	wg.Wait()

	if err := doc.Sanitize(); err != nil {
		h.Error(message.Erro, err)
	}

	return doc, nil
}

var (
	apidocBegin = []byte("<apidoc")
	apiBegin    = []byte("<api")
)

func parseBlock(d *doc.Doc, block i.Block, h *message.Handler) {
	var err error
	switch {
	case bytes.HasPrefix(block.Data, apidocBegin):
		err = d.FromXML(block.Data)
	case bytes.HasPrefix(block.Data, apiBegin):
		err = d.NewAPI(block.File, block.Line).FromXML(block.Data)
	}

	h.Error(message.Erro, message.WithError(block.File, "", block.Line, err))
}
