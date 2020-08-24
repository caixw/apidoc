// SPDX-License-Identifier: MIT

package ast

import (
	"bytes"
	"errors"
	"io"
	"sort"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
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
func (doc *APIDoc) Parse(h *core.MessageHandler, b core.Block) {
	if !isValid(b) {
		return
	}

	p, err := xmlenc.NewParser(h, b)
	if err != nil {
		h.Error(err)
		return
	}

	switch getTagName(p) {
	case "api":
		if doc.APIs == nil {
			doc.APIs = make([]*API, 0, 100)
		}

		api := &API{
			doc: doc,
			URI: b.Location.URI,
		}
		xmlenc.Decode(p, api, core.XMLNamespace)
		doc.APIs = append(doc.APIs, api)

		if doc.Title.V() != "" { // apidoc 已经初始化，检测依赖于 apidoc 的字段
			api.sanitizeTags(p)
		}
	case "apidoc":
		if doc.Title != nil { // 多个 apidoc 标签
			h.Error(p.NewError(b.Location.Range.Start, b.Location.Range.End, "apidoc", locale.ErrDuplicateValue))
			return
		}
		xmlenc.Decode(p, doc, core.XMLNamespace)
	default:
		return
	}

	// api 进入 doc 的顺序是未知的，进行排序可以保证文档的顺序一致。
	doc.sortAPIs()
}

// 简单预判是否是一个合规的 apidoc 内容
func isValid(b core.Block) bool {
	bs := bytes.TrimSpace(b.Data)
	if len(bs) < minSize {
		return false
	}

	// 去除空格之后，必须保证以 < 开头，且不能以 </ 开关。
	return bs[0] == '<' && bs[1] != '/'
}

// 获取根标签的名称
func getTagName(p *xmlenc.Parser) string {
	start := p.Current()
	for {
		t, _, err := p.Token()
		if errors.Is(err, io.EOF) {
			return ""
		} else if err != nil {
			p.Error(err)
			return ""
		}

		switch elem := t.(type) {
		case *xmlenc.StartElement:
			p.Move(start)
			return elem.Name.Local.Value
		case *xmlenc.EndElement, *xmlenc.CData: // 表示是一个非法的 XML，忽略！
			return ""
		default: // 其它标签忽略
		}
	}
}

func (doc *APIDoc) sortAPIs() {
	sort.SliceStable(doc.APIs, func(i, j int) bool {
		ii := doc.APIs[i]
		jj := doc.APIs[j]

		var iip string
		if ii.Path != nil && ii.Path.Path != nil {
			iip = ii.Path.Path.V()
		}

		var jjp string
		if jj.Path != nil && jj.Path.Path != nil {
			jjp = jj.Path.Path.V()
		}

		var iim string
		if ii.Method != nil {
			iim = ii.Method.V()
		}

		var jjm string
		if jj.Method != nil {
			jjm = jj.Method.V()
		}

		if iip == jjp {
			return iim < jjm
		}
		return iip < jjp
	})
}
