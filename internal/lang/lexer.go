// SPDX-License-Identifier: MIT

package lang

import (
	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/lexer"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/spec"
)

// Lexer 是对一个文本内容的包装，方便 blocker 等接口操作。
type Lexer struct {
	*lexer.Lexer
	blocks []Blocker
}

// NewLexer 声明 Lexer 实例
func NewLexer(data []byte, blocks []Blocker) (*Lexer, error) {
	l, err := lexer.New(data)
	if err != nil {
		return nil, err
	}

	return &Lexer{
		Lexer:  l,
		blocks: blocks,
	}, nil
}

// Block 从当前位置往后查找，直到找到第一个与 blocks 中某个相匹配的，并返回该 Blocker 。
func (l *Lexer) block() (Blocker, core.Position) {
	for {
		if l.AtEOF() {
			return nil, core.Position{}
		}

		pos := l.Position()
		for _, block := range l.blocks {
			if block.BeginFunc(l) {
				return block, pos.Position
			}
		}

		l.Next(1)
	}
}

// Parse 分析 l.data 的内容
//
// uri 表示在出错时，其返回的错误信息包含的定位信息。
func (l *Lexer) Parse(blocks chan spec.Block, h *core.MessageHandler, uri core.URI) {
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

		raw, data, ok := block.EndFunc(l)
		if !ok { // 没有找到结束标签，那肯定是到文件尾了，可以直接返回。
			loc := core.Location{
				URI: uri,
				Range: core.Range{
					Start: pos,
					End:   l.Position().Position,
				},
			}
			h.Error(core.Erro, core.NewLocaleError(loc, "", locale.ErrNotFoundEndFlag))
			return
		}

		block = nil // 重置 block

		if len(raw) == 0 {
			continue
		}

		blocks <- spec.Block{
			Location: core.Location{
				URI: uri,
				Range: core.Range{
					Start: pos,
					End:   l.Position().Position,
				},
			},
			Data: data,
			Raw:  raw,
		}
	} // end for
}
