// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// 用于描述 block.Type 的值。
const (
	blockTypeNone     int8 = iota
	blockTypeString        // 字符串，将被忽略。
	blockTypeSComment      // 单行注释
	blockTypeMComment      // 多行注释
)

// block 定义了与语言相关的三种类型的代码块：单行注释，多行注释，字符串。
type block struct {
	BeginFunc func(l *lexer) bool
	EndFunc   func(l *lexer) ([]rune, bool)

	Type   int8   // 代码块的类型，可以是字符串，单行注释或是多行注释。
	Begin  string // 块的起始字符串
	End    string // 块的结束字符串，单行注释不用定义此值
	Escape string // 当 Type 为 blockTypeString 时，此值表示转义字符，Type 为其它值时，此值无意义；
}

func (b *block) match(l *lexer) bool {
	if b.BeginFunc != nil {
		return b.BeginFunc(l)
	}

	return l.match(b.Begin)
}

// 返回从当前位置到定义结束的所有字符
// 返回值 bool 提示是否正常找到结束标记
func (b *block) end(l *lexer) ([]rune, bool) {
	switch {
	case b.EndFunc != nil:
		return b.EndFunc(l)
	case b.Type == blockTypeString:
		return b.endString(l)
	case b.Type == blockTypeMComment:
		return b.endMComments(l)
	case b.Type == blockTypeSComment:
		return b.endSComments(l)
	}

	return nil, false
}

// 从 l 的当前位置开始往后查找，直到找到 b 中定义的 end 字符串，
// 将 l 中的指针移到该位置。
// 正常找到结束符的返回 true，否则返回 false。
func (b *block) endString(l *lexer) ([]rune, bool) {
LOOP:
	for {
		switch {
		case l.atEOF():
			break LOOP
		case (len(b.Escape) > 0) && l.match(b.Escape):
			l.next()
		case l.match(b.End):
			return nil, true
		default:
			l.next()
		}
	} // end for
	return nil, false
}

// 从 l 的当前位置往后开始查找连续的相同类型单行代码块。
func (b *block) endSComments(l *lexer) ([]rune, bool) {
	// 跳过除换行符以外的所有空白字符。
	skipSpace := func() {
		for {
			r, w := utf8.DecodeRune(l.data[l.pos:])
			if !unicode.IsSpace(r) || r == '\n' {
				break
			}
			l.pos += w
		}
	} // end skipSpace

	ret := make([]rune, 0, 1000)
	for {
		for { // 读取一行的内容到 ret 变量中
			r := l.next()
			ret = append(ret, r)

			if l.atEOF() || r == '\n' {
				break
			}
		}

		skipSpace()            // 去掉新行的前导空格，若是存在的话。
		if !l.match(b.Begin) { // 不是接连着的注释块了，结束当前的匹配
			break
		}
	}

	if len(ret) > 0 { // 最后一个换行符返还给 lexer
		l.pos--
	}

	return ret, true
}

// 从 l 的当前位置一直到定义的 b.End 之间的所有字符。
// 会对每一行应用 filterSymbols 规则。
func (b *block) endMComments(l *lexer) ([]rune, bool) {
	lines := make([][]rune, 0, 20)
	line := make([]rune, 0, 100)

LOOP:
	for {
		switch {
		case l.atEOF():
			return nil, false
		case l.match(b.End):
			lines = append(lines, b.filterSymbols(line))
			break LOOP
		default:
			r := l.next()
			line = append(line, r)
			if r == '\n' {
				lines = append(lines, b.filterSymbols(line))
				line = make([]rune, 0, 100)
			}
		}
	}

	ret := make([]rune, 0, 1000)
	for _, v := range lines {
		ret = append(ret, v...)
	}
	return ret, true
}

// 行首若出现`空白字符+symbol+空白字符`的组合，则去掉这些字符。
// symbol 为 b.Begin 中的任意字符。
func (b *block) filterSymbols(line []rune) []rune {
	for k, v := range line {
		if unicode.IsSpace(v) { // 跳过行首的空格
			continue
		}

		// 不存在指定的符号，直接返回原数据
		if strings.IndexRune(b.Begin, v) < 0 {
			return line
		}

		// 若下个字符正好是是空格
		if len(line) > k+1 && unicode.IsSpace(line[k+1]) {
			return line[k+2:]
		} else {
			return line
		}
	}

	return line
}
