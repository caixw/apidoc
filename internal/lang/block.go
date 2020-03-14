// SPDX-License-Identifier: MIT

package lang

import (
	"bytes"
	"strings"
	"unicode"
)

// Blocker 接口定义了解析代码块的所有操作
type Blocker interface {
	// 确定 l 的当前位置是否匹配 Blocker 的起始位置。
	BeginFunc(l *Lexer) bool

	// 确定 l 的当前位置是否匹配 Blocker 的结束位置
	//
	// ok 表示是否正确匹配；
	// data 表示匹配的内容，如果不使用返回的内容，可以返回空值。
	// 比如字符串，只需要返回 true，以确保找到了结束位置，但是 data 可以直接返回 nil。
	// raw 表示匹配情况下的原始内容，data 返回的可能是经过处理，而 raw 应该是未处理的。
	//
	// 如果在到达文件末尾都没有找到结束符，则应该返回 nil, false
	EndFunc(l *Lexer) (raw, data []byte, ok bool)
}

type (
	stringBlock struct {
		begin, end, escape string
	}

	singleComment struct {
		begin  string
		begins []byte
	}

	multipleComment struct {
		begin, end, escape string
		begins, ends       []byte
	}
)

func newString(begin, end, escape string) Blocker {
	return &stringBlock{
		begin:  begin,
		end:    end,
		escape: escape,
	}
}

func newSingleComment(begin string) Blocker {
	return &singleComment{
		begin:  begin,
		begins: []byte(begin),
	}
}

func newMultipleComment(begin, end, escape string) Blocker {
	return &multipleComment{
		begin:  begin,
		end:    end,
		escape: escape,

		begins: []byte(begin),
		ends:   []byte(end),
	}
}

// BeginFunc 实现 Blocker.BeginFunc
func (b *stringBlock) BeginFunc(l *Lexer) bool {
	return l.match(b.begin)
}

// 从 l 的当前位置开始往后查找，直到找到 b 中定义的 end 字符串，
// 将 l 中的指针移到该位置。
// 正常找到结束符的返回 true，否则返回 false。
//
// 第一个返回参数无用，仅是为了统一函数签名
func (b *stringBlock) EndFunc(l *Lexer) (raw, data []byte, ok bool) {
	for {
		switch {
		case l.AtEOF():
			return nil, nil, false
		case (len(b.escape) > 0) && l.match(b.escape):
			l.next(1)
		case l.match(b.end):
			return nil, nil, true
		default:
			l.next(1)
		}
	} // end for
}

// BeginFunc 实现 Blocker.BeginFunc
func (b *singleComment) BeginFunc(l *Lexer) bool {
	return l.match(b.begin)
}

// 从 l 的当前位置往后开始查找连续的相同类型单行代码块。
func (b *singleComment) EndFunc(l *Lexer) (raw, data []byte, ok bool) {
	data = make([]byte, 0, 120)
	raw = make([]byte, 0, 120)

	for {
		raw = append(raw, b.begins...)
		bs := l.delim('\n')
		if bs == nil { // 找不到换行符，直接填充到末尾
			all := l.all()
			data = append(data, all...)
			raw = append(raw, all...)
			break
		} else {
			data = append(data, bs...)
			raw = append(raw, bs...)
		}

		spaces := l.skipSpace()
		raw = append(raw, spaces...)
		if !l.match(b.begin) { // 不是接连着的注释块了，结束当前的匹配
			break
		}
	}

	return raw, data, true
}

// BeginFunc 实现 Blocker.BeginFunc
func (b *multipleComment) BeginFunc(l *Lexer) bool {
	return l.match(b.begin)
}

// 从 l 的当前位置一直到定义的 b.End 之间的所有字符。
// 会对每一行应用 filterSymbols 规则。
func (b *multipleComment) EndFunc(l *Lexer) (raw, data []byte, ok bool) {
	data = make([]byte, 0, 200)
	raw = append(make([]byte, 0, 200), b.begins...)
	line := make([]byte, 0, 100)

	for {
		switch {
		case l.AtEOF(): // 没有找到结束符号，直接到达文件末尾
			return nil, nil, false
		case l.match(b.end):
			raw = append(raw, b.ends...)
			if len(line) > 0 {
				data = append(data, line...)
			}
			return raw, data, true
		default:
			bs := l.next(1)
			raw = append(raw, bs...)
			line = append(line, bs...)
			if isNewline(bs) {
				data = append(data, filterSymbols(line, b.escape)...)
				line = line[:0]
			}
		}
	} // end for
}

// 行首若出现`空白字符+symbol+空白字符`的组合，则去掉 symbol 及之前的字符。
// symbol 为 charset 中的任意字符。
func filterSymbols(line []byte, charset string) []byte {
	if len(charset) == 0 {
		return line
	}

	// 过滤左侧的空格
	line = bytes.TrimLeftFunc(line, func(r rune) bool { return unicode.IsSpace(r) && r != '\n' })
	if len(line) == 0 {
		return line
	}

	// 过滤左侧的符号
	hasSymbol := false
	for index, v := range line {
		if strings.IndexByte(charset, v) >= 0 {
			hasSymbol = true
			continue
		}

		if hasSymbol && unicode.IsSpace(rune(v)) {
			return line[index:]
		}
		break
	}

	return line
}
