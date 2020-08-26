// SPDX-License-Identifier: MIT

package xmlenc

import (
	"bytes"
	"errors"
	"io"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/lexer"
	"github.com/caixw/apidoc/v7/internal/locale"
)

const (
	cdataStart  = "<![CDATA["
	cdataEnd    = "]]>"
	cdataEscape = "]]]]><![CDATA[>"
)

// Parser 代码块的解析器
type Parser struct {
	*lexer.Lexer
	*core.MessageHandler
}

// NewParser 声明新的 Parser 实例
func NewParser(h *core.MessageHandler, b core.Block) (*Parser, error) {
	l, err := lexer.New(b)
	if err != nil {
		return nil, err
	}

	return &Parser{
		Lexer:          l,
		MessageHandler: h,
	}, nil
}

// Token 返回下一个 token 对象
//
// token 可能的类型为 *StartElement、*EndElement、*Instruction、*Attribute、*CData、*Comment 和 *String。
// 其中 *String 用于表示 XML 元素的内容。
//
// loc 表示返回的 token 所占的范围；
// 当返回 nil, {}, io.EOF 时，表示已经结束
func (p *Parser) Token() (token interface{}, loc core.Location, err error) {
	for {
		if p.AtEOF() {
			return nil, core.Location{}, io.EOF
		}

		pos := p.Current() // 记录元素的开始位置

		bs := p.Next(1)
		if len(bs) == 0 {
			return nil, core.Location{}, io.EOF
		}
		if len(bs) > 1 || bs[0] != '<' {
			p.Rollback() // 当前字符是内容的一部分，返回给 parseContent 解析
			return p.parseContent()
		}

		switch {
		case p.Match("?"):
			return p.parseInstruction(pos)
		case p.Match("![CDATA["):
			return p.parseCData(pos)
		case p.Match("/"):
			return p.parseEndElement(pos)
		case p.Match("!--"):
			return p.parseComment(pos)
		default:
			return p.parseStartElement(pos)
		}
	}
}

func (p *Parser) parseContent() (*String, core.Location, error) {
	start := p.Current()

	data, found := p.Delim('<', false)
	if !found {
		data = p.All()
		if len(data) == 0 {
			return nil, core.Location{}, io.EOF
		}
	}

	loc := core.Location{URI: p.Location.URI, Range: core.Range{Start: start.Position, End: p.Current().Position}}
	return &String{
		Value:    string(data),
		Location: loc,
	}, loc, nil
}

func (p *Parser) parseComment(pos lexer.Position) (*Comment, core.Location, error) {
	start := p.Current()

	data, found := p.DelimString("-->", false)
	if !found {
		return nil, core.Location{}, p.NewError(p.Current().Position, p.Current().Position, "<!--", locale.ErrNotFoundEndTag)
	}
	end := p.Current()
	p.Next(3) // 跳过 --> 三个字符

	loc := core.Location{
		URI:   p.Location.URI,
		Range: core.Range{Start: pos.Position, End: p.Current().Position},
	}
	return &Comment{
		Location: loc,
		Value: String{
			Location: core.Location{URI: p.Location.URI, Range: core.Range{Start: start.Position, End: end.Position}},
			Value:    string(data),
		},
	}, loc, nil
}

func (p *Parser) parseStartElement(pos lexer.Position) (*StartElement, core.Location, error) {
	p.Spaces(0) // 跳过空格

	start := p.Current()
	name, found := p.DelimFunc(func(r rune) bool { return unicode.IsSpace(r) || r == '/' || r == '>' }, false)
	if !found || len(name) == 0 {
		return nil, core.Location{}, p.NewError(p.Current().Position, p.Current().Position, "", locale.ErrInvalidXML)
	}

	elem := &StartElement{
		Location: core.Location{URI: p.Location.URI},
		Name:     parseName(name, p.Location.URI, start.Position, p.Current().Position),
	}

	attrs, err := p.parseAttributes()
	if err != nil {
		return nil, core.Location{}, err
	}
	elem.Attributes = attrs

	p.Spaces(0)
	if p.Match("/>") {
		elem.Range = core.Range{Start: pos.Position, End: p.Current().Position}
		elem.SelfClose = true
		return elem, elem.Location, nil
	}
	if p.Match(">") {
		elem.Range = core.Range{Start: pos.Position, End: p.Current().Position}
		return elem, elem.Location, nil
	}

	return nil, core.Location{}, p.NewError(p.Current().Position, p.Current().Position, string(name), locale.ErrNotFoundEndTag)
}

func parseName(name []byte, uri core.URI, start, end core.Position) Name {
	index := bytes.IndexByte(name, ':')
	if index < 0 {
		return Name{
			Location: core.Location{URI: uri, Range: core.Range{Start: start, End: end}},
			Local: String{
				Location: core.Location{URI: uri, Range: core.Range{Start: start, End: end}},
				Value:    string(name),
			},
		}
	}

	prefix := name[:index]
	character := start.Character + utf8.RuneCount(prefix)
	return Name{
		Location: core.Location{URI: uri, Range: core.Range{Start: start, End: end}},
		Prefix: String{
			Location: core.Location{URI: uri, Range: core.Range{Start: start, End: core.Position{Line: start.Line, Character: character}}},
			Value:    string(prefix),
		},
		Local: String{
			Location: core.Location{URI: uri, Range: core.Range{Start: core.Position{Line: start.Line, Character: character + 1}, End: end}},
			Value:    string(name[index+1:]),
		},
	}
}

// pos 表示当前元素的起始位置，包含了 < 元素
func (p *Parser) parseEndElement(pos lexer.Position) (*EndElement, core.Location, error) {
	start := p.Current() // 名称开始的定位，传递过来的 pos 表示的 < 起始位置

	name, found := p.Delim('>', false)
	if !found || len(name) == 0 {
		return nil, core.Location{}, p.NewError(p.Current().Position, p.Current().Position, "", locale.ErrInvalidXML)
	}
	end := p.Current()
	p.Next(1) // 去掉 > 符号

	loc := core.Location{
		URI:   p.Location.URI,
		Range: core.Range{Start: pos.Position, End: p.Current().Position},
	}
	return &EndElement{
		Location: loc,
		Name:     parseName(name, p.Location.URI, start.Position, end.Position),
	}, loc, nil
}

// pos 表示当前元素的起始位置，包含了 < 元素
func (p *Parser) parseCData(pos lexer.Position) (*CData, core.Location, error) {
	start := p.Current()
	var value []byte

	for {
		v, found := p.DelimString(cdataEnd, false)
		if !found {
			return nil, core.Location{}, p.NewError(pos.Position, p.Current().Position, cdataStart, locale.ErrNotFoundEndTag)
		}
		value = append(value, v...)

		curr := p.Current()
		p.Move(curr.SubRune(']').SubRune(']')) // 回滚两个字符，用于匹配转义内容
		if p.Match(cdataEscape) {
			value = append(value, cdataEscape[2:]...)
			continue
		}

		p.Move(curr)
		break
	}

	end := p.Current()
	p.Next(len(cdataEnd)) // 将 ]]> 从流中去掉

	loc := core.Location{
		URI:   p.Location.URI,
		Range: core.Range{Start: pos.Position, End: p.Current().Position},
	}
	return &CData{
		BaseTag: BaseTag{
			Base: Base{
				Location: loc,
			},
			StartTag: Name{
				Local: String{
					Location: core.Location{
						URI: p.Location.URI,
						Range: core.Range{
							Start: pos.Position,
							End:   core.Position{Line: pos.Line, Character: pos.Character + len(cdataStart)},
						},
					},
					Value: cdataStart,
				},
				Location: core.Location{
					URI: p.Location.URI,
					Range: core.Range{
						Start: pos.Position,
						End:   core.Position{Line: pos.Line, Character: pos.Character + len(cdataStart)},
					},
				},
			},
			EndTag: Name{
				Local: String{
					Location: core.Location{
						URI: p.Location.URI,
						Range: core.Range{
							Start: end.Position,
							End:   core.Position{Line: end.Line, Character: end.Character + len(cdataEnd)},
						},
					},
					Value: cdataEnd,
				},
				Location: core.Location{
					URI: p.Location.URI,
					Range: core.Range{
						Start: end.Position,
						End:   core.Position{Line: end.Line, Character: end.Character + len(cdataEnd)},
					},
				},
			},
		},
		Value: String{
			Location: core.Location{
				URI:   p.Location.URI,
				Range: core.Range{Start: start.Position, End: end.Position},
			},
			Value: string(value),
		},
	}, loc, nil
}

// pos 表示当前元素的起始位置，包含了 < 元素
func (p *Parser) parseInstruction(pos lexer.Position) (*Instruction, core.Location, error) {
	name, nameRange := p.getName()
	if len(name) == 0 {
		return nil, core.Location{}, p.NewError(p.Current().Position, p.Current().Position, "", locale.ErrInvalidXML)
	}
	elem := &Instruction{
		Location: core.Location{URI: p.Location.URI},
		Name: String{
			Location: core.Location{
				URI:   p.Location.URI,
				Range: nameRange,
			},
			Value: string(name),
		},
	}

	attrs, err := p.parseAttributes()
	if err != nil {
		return nil, core.Location{}, err
	}
	elem.Attributes = attrs

	p.Spaces(0)
	if p.Match("?>") {
		elem.Range = core.Range{Start: pos.Position, End: p.Current().Position}
		return elem, elem.Location, nil
	}

	return nil, core.Location{}, p.NewError(p.Current().Position, p.Current().Position, "<?", locale.ErrNotFoundEndTag)
}

func (p *Parser) parseAttributes() (attrs []*Attribute, err error) {
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
	p.Spaces(0)

	return attrs, nil
}

func (p *Parser) parseAttribute() (*Attribute, error) {
	p.Spaces(0) // 忽略空格
	pos, start := p.Current(), p.Current()

	name, nameRange := p.getName()
	if len(name) == 0 {
		return nil, nil
	}
	attr := &Attribute{
		Location: core.Location{URI: p.Location.URI},
		Name:     parseName(name, p.Location.URI, nameRange.Start, nameRange.End),
	}

	p.Spaces(0)
	if !p.Match("=") {
		return nil, p.NewError(p.Current().Position, p.Current().Position, "", locale.ErrInvalidXML)
	}

	p.Spaces(0)
	if !p.Match("\"") {
		return nil, p.NewError(p.Current().Position, p.Current().Position, "", locale.ErrInvalidXML)
	}

	pos = p.Current()
	value, found := p.Delim('"', true)
	if !found || len(value) == 0 {
		return nil, p.NewError(p.Current().Position, p.Current().Position, "", locale.ErrInvalidXML)
	}
	end := p.Current().SubRune('"') // 不包含 " 符号
	attr.Value = String{
		Location: core.Location{
			URI:   p.Location.URI,
			Range: core.Range{Start: pos.Position, End: end.Position},
		},
		Value: string(p.Bytes(pos.Offset, end.Offset)),
	}

	attr.Range = core.Range{Start: start.Position, End: p.Current().Position}

	return attr, nil
}

func (p *Parser) getName() ([]byte, core.Range) {
	start := p.Current()

	for {
		if p.AtEOF() {
			break
		}

		rs := p.Next(1)
		if len(rs) != 1 {
			continue
		}

		b := rs[0]
		if b == '"' || b == '=' || b == '<' || b == '>' || b == '?' || b == '/' || unicode.IsSpace(rune(b)) {
			p.Rollback()
			break
		}
	}

	end := p.Current()
	return p.Bytes(start.Offset, end.Offset), core.Range{Start: start.Position, End: end.Position}
}

// 找到与 start 相对应的结束符号位置
//
// 如果找不到对应的结束符号，则会向 p.h 输出一条错误信息，然后将定位至原始位置，并不返回错误。
// 这样可以保证最大限度地解析 xml 内容，不会因为一些非致命的错误而中断整个解析。
func (p *Parser) endElement(start *StartElement) error {
	if start.SelfClose {
		return nil
	}

	curr := p.Current()

	level := 0
	for {
		t, _, err := p.Token()
		if errors.Is(err, io.EOF) {
			p.Error(start.NewError(locale.ErrNotFoundEndTag).WithField(start.Name.String()))
			p.Move(curr)
			return nil
		} else if err != nil {
			return err
		}

		switch elem := t.(type) {
		case *StartElement:
			if elem.Name.Equal(start.Name) {
				level++
			}
		case *EndElement:
			if level == 0 && start.Match(elem) {
				return nil
			}
			level--
		}
	}
}

// NewError 生成 *core.Error 对象
//
// 其中的 URI 来自于 p.Location.URI
func (p *Parser) NewError(start, end core.Position, field string, key message.Reference, v ...interface{}) *core.Error {
	return core.NewError(key, v...).
		WithField(field).
		WithLocation(core.Location{
			URI:   p.Location.URI,
			Range: core.Range{Start: start, End: end},
		})
}
