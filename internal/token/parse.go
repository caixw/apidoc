// SPDX-License-Identifier: MIT

package token

import (
	"unicode"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/lexer"
	"github.com/caixw/apidoc/v6/internal/locale"
)

type parser struct {
	block *core.Block
	l     *lexer.Lexer
}

func newParser(b *core.Block) (*parser, error) {
	l, err := lexer.New(b.Data)
	if err != nil {
		return nil, err
	}

	return &parser{
		block: b,
		l:     l,
	}, nil
}

func (p *parser) token() (interface{}, error) {
	pos := p.l.Position()
	for {
		if p.l.AtEOF() {
			return nil, nil
		}

		bs := p.l.Next(1)
		if len(bs) == 0 {
			return nil, nil
		}
		if len(bs) > 1 || bs[0] != '<' {
			continue
		}

		switch {
		case p.l.Match("?"):
			return p.parseInstruction(pos)
		case p.l.Match("![CDATA["):
			return p.parseCData(pos)
		case p.l.Match("/"):
			return p.parseEndElement(pos)
		default:
			return p.parseStartElement(pos)
		}
	}
}

func (p *parser) parseStartElement(pos lexer.Position) (*StartElement, error) {
	p.l.Spaces() // 跳过空格

	start := p.l.Position()
	name := p.l.DelimFunc(func(r rune) bool { return unicode.IsSpace(r) }, false)
	elem := &StartElement{
		Range: core.Range{
			Start: pos.Position,
		},
		Name: String{
			Range: core.Range{Start: start.Position, End: p.l.Position().Position},
			Value: string(name),
		},
	}

	attrs, err := p.parseAttributes()
	if err != nil {
		return nil, err
	}
	elem.Attributes = attrs

	if p.l.Match("/>") {
		elem.End = p.l.Position().Position
		elem.Close = true
		return elem, nil
	}
	if p.l.Match(">") {
		elem.End = p.l.Position().Position
		return elem, nil
	}

	loc := core.Location{
		URI:   p.block.Location.URI,
		Range: core.Range{Start: p.l.Position().Position, End: p.l.Position().Position},
	}
	return nil, core.NewLocaleError(loc, "", locale.ErrInvalidXML)
}

// pos 表示当前元素的起始位置，包含了 < 元素
func (p *parser) parseEndElement(pos lexer.Position) (*EndElement, error) {
	// 名称开始的定位，传递过来的 pos 表示的 < 起始位置
	start := p.l.Position()

	name := p.l.Delim('>', true)
	if len(name) == 0 {
		loc := core.Location{
			URI:   p.block.Location.URI,
			Range: core.Range{Start: p.l.Position().Position, End: p.l.Position().Position},
		}
		return nil, core.NewLocaleError(loc, "", locale.ErrInvalidXML)
	}

	return &EndElement{
		Range: core.Range{
			Start: pos.Position,
			End:   p.l.Position().Position,
		},
		Name: String{
			Range: core.Range{
				Start: start.Position,
				End:   p.l.Position().SubRune('>').Position,
			},
			Value: string(name),
		},
	}, nil
}

// pos 表示当前元素的起始位置，包含了 < 元素
func (p *parser) parseCData(pos lexer.Position) (*CData, error) {
	start := p.l.Position()

	value := p.l.DelimString("]]>")
	if len(value) == 0 {
		loc := core.Location{
			URI: p.block.Location.URI,
			Range: core.Range{
				Start: pos.Position,
				End:   p.l.Position().Position,
			},
		}
		return nil, core.NewLocaleError(loc, "", locale.ErrInvalidXML)
	}

	return &CData{
		Range: core.Range{
			Start: pos.Position,
			End:   p.l.Position().Position,
		},
		Value: String{
			Range: core.Range{
				Start: start.Position,
				End:   p.l.Position().SubRune('>').SubRune(']').SubRune(']').Position,
			},
			Value: string(value),
		},
	}, nil
}

// pos 表示当前元素的起始位置，包含了 < 元素
func (p *parser) parseInstruction(pos lexer.Position) (*Instruction, error) {
	start := p.l.Position()
	name := p.l.DelimFunc(func(r rune) bool { return unicode.IsSpace(r) }, false)
	elem := &Instruction{
		Range: core.Range{
			Start: pos.Position,
		},
		Name: String{
			Range: core.Range{Start: start.Position, End: p.l.Position().Position},
			Value: string(name),
		},
	}

	attrs, err := p.parseAttributes()
	if err != nil {
		return nil, err
	}
	elem.Attributes = attrs

	if p.l.Match("?>") {
		elem.End = p.l.Position().Position
		return elem, nil
	}

	loc := core.Location{
		URI:   p.block.Location.URI,
		Range: core.Range{Start: p.l.Position().Position, End: p.l.Position().Position},
	}
	return nil, core.NewLocaleError(loc, "", locale.ErrInvalidXML)
}

func (p *parser) parseAttributes() ([]*Attribute, error) {
	attrs := make([]*Attribute, 0, 10)

	for {
		attr, err := p.parseAttribute()
		if err != nil {
			return nil, err
		}
		if attr == nil {
			break
		}

		attrs = append(attrs, attr)
	}

	p.l.Spaces()

	return attrs, nil
}

func (p *parser) parseAttribute() (*Attribute, error) {
	p.l.Spaces() // 忽略空格

	pos := p.l.Position()
	name := p.l.DelimFunc(func(r rune) bool { return unicode.IsSpace(r) || r == '=' }, false)
	if len(name) == 0 {
		return nil, nil
	}

	attr := &Attribute{Name: String{
		Range: core.Range{Start: pos.Position, End: p.l.Position().Position},
		Value: string(name),
	}}

	p.l.Spaces()
	pos = p.l.Position()
	if !p.l.Match("=") {
		loc := core.Location{
			URI:   p.block.Location.URI,
			Range: core.Range{Start: pos.Position, End: p.l.Position().Position},
		}
		return nil, core.NewLocaleError(loc, "", locale.ErrInvalidXML)
	}

	p.l.Spaces()
	pos = p.l.Position()
	if !p.l.Match("\"") {
		loc := core.Location{
			URI:   p.block.Location.URI,
			Range: core.Range{Start: pos.Position, End: p.l.Position().Position},
		}
		return nil, core.NewLocaleError(loc, "", locale.ErrInvalidXML)
	}

	pos = p.l.Position().SubRune('"') // 回滚 " 符号，作为属性值的一部分
	value := p.l.Delim('"', true)
	if len(value) == 0 {
		loc := core.Location{
			URI:   p.block.Location.URI,
			Range: core.Range{Start: pos.Position, End: p.l.Position().Position},
		}
		return nil, core.NewLocaleError(loc, "", locale.ErrInvalidXML)
	}
	attr.Value = String{
		Range: core.Range{Start: pos.Position, End: p.l.Position().Position},
		Value: string(p.l.Bytes(pos.Offset, p.l.Position().Offset)),
	}

	return attr, nil
}
