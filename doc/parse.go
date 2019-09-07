// SPDX-License-Identifier: MIT

package doc

import (
	"bytes"
	"context"
	"encoding/xml"
	"sort"
	"sync"

	"github.com/caixw/apidoc/v5/errors"
	i "github.com/caixw/apidoc/v5/internal/input"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/options"
)

// Parse 分析从 block 中获取的代码块。并填充到 Doc 中
//
// 当所有的代码块已经放入 Block 之后，Block 会被关闭。
//
// 所有与解析有关的错误均通过 h 输出。而其它错误，比如参数问题等，通过返回参数返回。
func Parse(ctx context.Context, h *errors.Handler, input ...*options.Input) (*Doc, error) {
	block, err := i.Parse(ctx, h, input...)
	if err != nil {
		return nil, err
	}

	doc := &Doc{
		APIDoc: vars.Version(),
	}
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
				doc.parseBlock(b, h)
				wg.Done()
			}(blk)
		}
	}

	wg.Wait()

	doc.check(h)

	return doc, nil
}

func (doc *Doc) check(h *errors.Handler) {
	// Tag.Name 查重
	sort.SliceStable(doc.Tags, func(i, j int) bool {
		return doc.Tags[i].Name > doc.Tags[j].Name
	})
	for i := 1; i < len(doc.Tags); i++ {
		if doc.Tags[i].Name == doc.Tags[i-1].Name {
			h.SyntaxError(errors.New(doc.file, "", doc.line, locale.ErrDuplicateTag))
			return
		}
	}

	// Server.Name 查重
	sort.SliceStable(doc.Servers, func(i, j int) bool {
		return doc.Servers[i].Name > doc.Servers[j].Name
	})
	for i := 1; i < len(doc.Tags); i++ {
		if doc.Servers[i].Name == doc.Servers[i-1].Name {
			h.SyntaxError(errors.New(doc.file, "", doc.line, locale.ErrDuplicateTag))
			return
		}
	}

	// Server.URL 查重
	sort.SliceStable(doc.Servers, func(i, j int) bool {
		return doc.Servers[i].URL > doc.Servers[j].URL
	})
	for i := 1; i < len(doc.Tags); i++ {
		if doc.Servers[i].URL == doc.Servers[i-1].URL {
			h.SyntaxError(errors.New(doc.file, "", doc.line, locale.ErrDuplicateTag))
			return
		}
	}

	for _, api := range doc.Apis {
		for _, tag := range api.Tags {
			if !doc.tagExists(tag) {
				h.SyntaxError(errors.New(api.file, "", api.line, locale.ErrInvalidValue))
			}
		}

		for _, srv := range api.Servers {
			if !doc.serverExists(srv) {
				h.SyntaxError(errors.New(api.file, "", api.line, locale.ErrInvalidValue))
			}
		}
	} // end doc.Apis
}

var (
	apidocBegin = []byte("<apidoc ")
	apiBegin    = []byte("<api ")
)

func (doc *Doc) parseBlock(block i.Block, h *errors.Handler) {
	switch {
	case bytes.HasPrefix(block.Data, apidocBegin):
		err := xml.Unmarshal(block.Data, doc)
		if serr, ok := err.(*xml.SyntaxError); ok {
			h.SyntaxError(errors.New(block.File, "", block.Line+serr.Line, serr.Msg))
		}
	case bytes.HasPrefix(block.Data, apiBegin):
		api := &API{}
		err := xml.Unmarshal(block.Data, api)
		if err != nil {
			if serr, ok := err.(*xml.SyntaxError); ok {
				h.SyntaxError(errors.New(block.File, "", block.Line+serr.Line, serr.Msg))
			}
			return
		}
		api.line = block.Line
		api.file = block.File
		doc.Apis = append(doc.Apis, api)
	}
}

func (doc *Doc) tagExists(tag string) bool {
	for _, t := range doc.Tags {
		if t.Name == tag {
			return true
		}
	}

	return false
}

func (doc *Doc) serverExists(srv string) bool {
	for _, s := range doc.Servers {
		if s.Name == srv {
			return true
		}
	}

	return false
}
