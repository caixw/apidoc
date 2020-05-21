// SPDX-License-Identifier: MIT

package ast

import (
	"io"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/token"
)

// ParseBlocks 从多个 core.Block 实例中解析文档内容
//
// g 必须是一个阻塞函数，直到所有代码块都写入参数之后，才能返回。
func (doc *APIDoc) ParseBlocks(h *core.MessageHandler, g func(chan core.Block)) {
	done := make(chan struct{})
	blocks := make(chan core.Block, 50)

	go func() {
		for block := range blocks {
			doc.Parse(h, block)
		}
		done <- struct{}{}
	}()

	g(blocks)
	close(blocks)
	<-done
}

// Parse 将注释块的内容添加到当前文档
//
// 分析注释块内容，如果正确，则添加到当前文档中。
//
// 如果内容不是文档内容，刚将返回 false
func (doc *APIDoc) Parse(h *core.MessageHandler, b core.Block) {
	if len(b.Data) < minSize {
		return
	}

	p, err := token.NewParser(b)
	if err != nil {
		h.Error(err)
		return
	}

	name, err := getTagName(p)
	if err == io.EOF {
		return
	} else if err != nil {
		h.Error(err)
		return
	}
	switch name {
	case "api":
		if doc.Apis == nil {
			doc.Apis = make([]*API, 0, 100)
		}

		api := &API{doc: doc}
		token.Decode(h, p, api)
		doc.Apis = append(doc.Apis, api)
	case "apidoc":
		if doc.Title != nil { // 多个 apidoc 标签
			h.Error(p.NewError(b.Location.Range.Start, b.Location.Range.End, "apidoc", locale.ErrDuplicateValue))
			return
		}
		token.Decode(h, p, doc)
	default:
		return
	}

	doc.sortAPIs()
}

// 获取根标签的名称
func getTagName(p *token.Parser) (string, error) {
	start := p.Position()
	for {
		t, r, err := p.Token()
		if err != nil {
			return "", err
		}

		switch elem := t.(type) {
		case *token.StartElement:
			p.Move(start)
			return elem.Name.Value, nil
		case *token.EndElement, *token.CData:
			return "", p.NewError(r.Start, r.End, "", locale.ErrInvalidXML)
		default: // 其它标签忽略
		}
	}
}
