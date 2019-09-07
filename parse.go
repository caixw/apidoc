// SPDX-License-Identifier: MIT

package apidoc

import (
	"bytes"
	"context"
	"encoding/xml"
	"sync"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/errors"
	i "github.com/caixw/apidoc/v5/internal/input"
	"github.com/caixw/apidoc/v5/options"
)

// Parse 分析从 block 中获取的代码块。并填充到 Doc 中
//
// 当所有的代码块已经放入 Block 之后，Block 会被关闭。
//
// 所有与解析有关的错误均通过 h 输出。而其它错误，比如参数问题等，通过返回参数返回。
func Parse(ctx context.Context, h *errors.Handler, input ...*options.Input) (*doc.Doc, error) {
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
	apidocBegin = []byte("<apidoc ")
	apiBegin    = []byte("<api ")
)

func parseBlock(d *doc.Doc, block i.Block, h *errors.Handler) {
	switch {
	case bytes.HasPrefix(block.Data, apidocBegin):
		err := xml.Unmarshal(block.Data, d)
		h.SyntaxError(errors.WithError(err, block.File, "", block.Line))
	case bytes.HasPrefix(block.Data, apiBegin):
		api := d.NewAPI(block.File, block.Line)
		err := xml.Unmarshal(block.Data, api)
		if err != nil {
			h.SyntaxError(errors.WithError(err, block.File, "", block.Line))
			return
		}
	}
}
