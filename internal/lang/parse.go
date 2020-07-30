// SPDX-License-Identifier: MIT

package lang

import (
	"fmt"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/lexer"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// Parse 分析 data 的内容并输出到到 blocks
func Parse(h *core.MessageHandler, langID string, data core.Block, blocks chan core.Block) {
	l := Get(langID)
	if l == nil {
		panic(fmt.Sprintf("%s 指定的语言解析器并不存在", langID))
	}

	newParser(h, data, l.blocks).parse(blocks)
}

type parser struct {
	*lexer.Lexer
	blocks []Blocker
	h      *core.MessageHandler
}

func newParser(h *core.MessageHandler, block core.Block, blocks []Blocker) *parser {
	l, err := lexer.New(block)
	if err != nil {
		h.Error(err)
		return nil
	}

	return &parser{
		Lexer:  l,
		blocks: blocks,
		h:      h,
	}
}

// 从当前位置往后查找，直到找到第一个与 blocks 中某个相匹配的，并返回该 Blocker 。
func (l *parser) block() (Blocker, core.Position) {
	for {
		if l.AtEOF() {
			return nil, core.Position{}
		}

		pos := l.Current()
		for _, block := range l.blocks {
			if block.BeginFunc(l) {
				return block, pos.Position
			}
		}

		l.Next(1)
	}
}

// 分析 l.data 的内容并输出到 blocks
func (l *parser) parse(blocks chan core.Block) {
	var block Blocker
	var pos core.Position
	for {
		if l.AtEOF() {
			return
		}

		if block == nil {
			if block, pos = l.block(); block == nil { // 没有匹配的 block 了
				return
			}
		}

		data, ok := block.EndFunc(l)
		if !ok { // 没有找到结束标签，那肯定是到文件尾了，可以直接返回。
			loc := core.Location{
				URI: l.Location.URI,
				Range: core.Range{
					Start: pos,
					End:   l.Current().Position,
				},
			}
			l.h.Error(core.NewSyntaxError(loc, "", locale.ErrNotFoundEndFlag))
			return
		}

		block = nil // 重置 block

		if len(data) == 0 {
			continue
		}

		blocks <- core.Block{
			Location: core.Location{
				URI: l.Location.URI,
				Range: core.Range{
					Start: pos,
					End:   l.Current().Position,
				},
			},
			Data: data,
		}
	} // end for
}
