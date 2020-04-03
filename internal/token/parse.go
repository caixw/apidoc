// SPDX-License-Identifier: MIT

package token

import (
	"unicode"

	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/lexer"
	"github.com/caixw/apidoc/v6/internal/locale"
)

type parser struct {
	block *core.Block
	l     *lexer.Lexer
	err   error // 记录最后一次错误信息
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
	if p.err != nil {
		return nil, p.err
	}

	for {
		if p.l.AtEOF() {
			return nil, nil
		}

		pos := p.l.Position()
		bs := p.l.Next(1)
		if len(bs) == 0 {
			return nil, nil
		}
		if len(bs) > 1 || bs[0] != '<' {
			p.l.Rollback() // 当前字符是内容的一部分，返回给 parseContent 解析
			return p.parseContent()
		}

		var ret interface{}
		var err error
		switch {
		case p.l.Match("?"):
			ret, err = p.parseInstruction(pos)
		case p.l.Match("![CDATA["):
			ret, err = p.parseCData(pos)
		case p.l.Match("/"):
			ret, err = p.parseEndElement(pos)
		default:
			ret, err = p.parseStartElement(pos)
		}

		if err != nil {
			p.err = err
			return nil, err
		}
		return ret, nil
	}
}

func (p *parser) parseContent() (*String, error) {
	start := p.l.Position()

	data, found := p.l.Delim('<', false)
	if !found {
		return nil, p.newError(p.l.Position().Position, p.l.Position().Position, locale.ErrInvalidXML)
	}

	return &String{
		Value: string(data),
		Range: p.fixRange(start.Position, p.l.Position().Position),
	}, nil
}

func (p *parser) parseStartElement(pos lexer.Position) (*StartElement, error) {
	p.l.Spaces(0) // 跳过空格

	start := p.l.Position()
	name, found := p.l.DelimFunc(func(r rune) bool { return unicode.IsSpace(r) || r == '/' || r == '>' }, false)
	if !found || len(name) == 0 {
		return nil, p.newError(p.l.Position().Position, p.l.Position().Position, locale.ErrInvalidXML)
	}

	elem := &StartElement{
		Name: String{
			Range: p.fixRange(start.Position, p.l.Position().Position),
			Value: string(name),
		},
	}

	attrs, err := p.parseAttributes()
	if err != nil {
		return nil, err
	}
	elem.Attributes = attrs

	if p.l.Match("/>") {
		elem.Range = p.fixRange(pos.Position, p.l.Position().Position)
		elem.Close = true
		return elem, nil
	}
	if p.l.Match(">") {
		elem.Range = p.fixRange(pos.Position, p.l.Position().Position)
		return elem, nil
	}

	return nil, p.newError(p.l.Position().Position, p.l.Position().Position, locale.ErrInvalidXML)
}

// pos 表示当前元素的起始位置，包含了 < 元素
func (p *parser) parseEndElement(pos lexer.Position) (*EndElement, error) {
	// 名称开始的定位，传递过来的 pos 表示的 < 起始位置
	start := p.l.Position()

	name, found := p.l.Delim('>', false)
	if !found || len(name) == 0 {
		return nil, p.newError(p.l.Position().Position, p.l.Position().Position, locale.ErrInvalidXML)
	}
	end := p.l.Position()
	p.l.Next(1) // 去掉 > 符号

	return &EndElement{
		Range: p.fixRange(pos.Position, p.l.Position().Position),
		Name: String{
			Range: p.fixRange(start.Position, end.Position),
			Value: string(name),
		},
	}, nil
}

// pos 表示当前元素的起始位置，包含了 < 元素
func (p *parser) parseCData(pos lexer.Position) (*CData, error) {
	start := p.l.Position()

	value, found := p.l.DelimString("]]>", false)
	if !found {
		return nil, p.newError(pos.Position, p.l.Position().Position, locale.ErrInvalidXML)
	}
	end := p.l.Position()
	p.l.Next(3) // 去掉 ]]>

	return &CData{
		Range: p.fixRange(pos.Position, p.l.Position().Position),
		Value: String{
			Range: p.fixRange(start.Position, end.Position),
			Value: string(value),
		},
	}, nil
}

// pos 表示当前元素的起始位置，包含了 < 元素
func (p *parser) parseInstruction(pos lexer.Position) (*Instruction, error) {
	name, nameRange := p.getName()
	if len(name) == 0 {
		return nil, p.newError(p.l.Position().Position, p.l.Position().Position, locale.ErrInvalidXML)
	}
	elem := &Instruction{
		Name: String{
			Range: nameRange,
			Value: string(name),
		},
	}

	attrs, err := p.parseAttributes()
	if err != nil {
		return nil, err
	}
	elem.Attributes = attrs

	if p.l.Match("?>") {
		elem.Range = p.fixRange(pos.Position, p.l.Position().Position)
		return elem, nil
	}

	return nil, p.newError(p.l.Position().Position, p.l.Position().Position, locale.ErrInvalidXML)
}

func (p *parser) parseAttributes() (attrs []*Attribute, err error) {
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
	p.l.Spaces(0)

	return attrs, nil
}

func (p *parser) parseAttribute() (*Attribute, error) {
	p.l.Spaces(0) // 忽略空格
	pos, start := p.l.Position(), p.l.Position()

	name, nameRange := p.getName()
	if len(name) == 0 {
		return nil, nil
	}
	attr := &Attribute{
		Name: String{
			Range: nameRange,
			Value: string(name),
		},
	}

	p.l.Spaces(0)
	if !p.l.Match("=") {
		return nil, p.newError(p.l.Position().Position, p.l.Position().Position, locale.ErrInvalidXML)
	}

	p.l.Spaces(0)
	if !p.l.Match("\"") {
		return nil, p.newError(p.l.Position().Position, p.l.Position().Position, locale.ErrInvalidXML)
	}

	pos = p.l.Position()
	value, found := p.l.Delim('"', true)
	if !found || len(value) == 0 {
		return nil, p.newError(p.l.Position().Position, p.l.Position().Position, locale.ErrInvalidXML)
	}
	end := p.l.Position().SubRune('"') // 不包含 " 符号
	attr.Value = String{
		Range: p.fixRange(pos.Position, end.Position),
		Value: string(p.l.Bytes(pos.Offset, end.Offset)),
	}

	attr.Range = p.fixRange(start.Position, p.l.Position().Position)

	return attr, nil
}

func (p *parser) getName() ([]byte, core.Range) {
	start := p.l.Position()

	for {
		if p.l.AtEOF() {
			break
		}

		rs := p.l.Next(1)
		if len(rs) != 1 {
			continue
		}

		b := rs[0]
		if b == '"' || b == '=' || b == '<' || b == '>' || b == '?' || b == '/' || unicode.IsSpace(rune(b)) {
			p.l.Rollback()
			break
		}
	}

	end := p.l.Position()

	return p.l.Bytes(start.Offset, end.Offset), p.fixRange(start.Position, end.Position)
}

func (p *parser) fixRange(start, end core.Position) core.Range {
	l := p.block.Location

	if start.Line == 0 {
		start.Character += l.Range.Start.Character
	}

	if end.Line == 0 {
		end.Character += l.Range.Start.Character
	}

	start.Line += l.Range.Start.Line
	end.Line += l.Range.Start.Line

	return core.Range{Start: start, End: end}
}

func (p *parser) newError(start, end core.Position, key message.Reference, v ...interface{}) error {
	return core.NewLocaleError(core.Location{
		URI:   p.block.Location.URI,
		Range: p.fixRange(start, end),
	}, "", key, v...)
}
