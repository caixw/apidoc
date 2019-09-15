// SPDX-License-Identifier: MIT

package apidoc

import (
	"bytes"
	"context"
	"sync"

	"github.com/caixw/apidoc/v5/doc"
	i "github.com/caixw/apidoc/v5/internal/input"
	"github.com/caixw/apidoc/v5/message"
	"github.com/caixw/apidoc/v5/options"
)

// Parse 分析从 block 中获取的代码块。并填充到 Doc 中
//
// 当所有的代码块已经放入 Block 之后，Block 会被关闭。
//
// 所有与解析有关的错误均通过 h 输出。而其它错误，比如参数问题等，通过返回参数返回。
func Parse(ctx context.Context, h *message.Handler, input ...*options.Input) (*doc.Doc, error) {
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

	h.Error(message.WithError(block.File, "", block.Line, err))
}
