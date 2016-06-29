// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/caixw/apidoc/app"
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
	Type   int8   // 代码块的类型，可以是字符串，单行注释或是多行注释
	Begin  string // 块的起始字符串
	End    string // 块的结束字符串，单行注释不用定义此值
	Escape string // 当 Type 为 blockTypeString 时，此值表示转义字符，Type 为其它值时，此值无意义；
}

type lexer struct {
	data    []byte
	pos     int
	isAtEOF bool
}

// 是否已经在文件末尾。
func (l *lexer) atEOF() bool {
	return l.isAtEOF || l.pos >= len(l.data)
}

// 获取当前的字符，并将指针指向下一个字符。
func (l *lexer) next() rune {
	r, w := utf8.DecodeRune(l.data[l.pos:])
	l.pos += w

	if r == utf8.RuneError && w == 0 { // EOF
		l.isAtEOF = true
	}

	return r
}

// 接下来的 n 个字符是否匹配指定的字符串，
// 若匹配，则将指定移向该字符串这后，否则不作任何操作。
func (l *lexer) match(word string) bool {
	if l.atEOF() || (l.pos+len(word) > len(l.data)) { // 剩余字符没有word长，直接返回false
		return false
	}

	width := 0
	for _, r := range word {
		rr, w := utf8.DecodeRune(l.data[l.pos:])
		if rr != r {
			l.pos -= width
			return false
		}

		l.pos += w
		width += w
	}

	return true
}

var newLine = []byte("\n")

func (l *lexer) lineNumber() int {
	return bytes.Count(l.data[:l.pos], newLine)
}

// 构建一个语法错误的信息。
func (l *lexer) syntaxError(msg string) *app.SyntaxError {
	return &app.SyntaxError{
		Line:    l.lineNumber(),
		Message: msg,
	}
}

// 从当前位置往后查找，直到找到第一个与 blocks 中某个相匹配的，并返回该 block 。
func (l *lexer) block(blocks []*block) *block {
	for {
		if l.atEOF() {
			return nil
		}

		for _, block := range blocks {
			if l.match(block.Begin) {
				return block
			}
		}

		l.next()
	}
}

// 返回从当前位置到定义结束的所有字符
func (b *block) end(l *lexer) ([]rune, *app.SyntaxError) {
	var rs []rune
	var err *app.SyntaxError

	switch b.Type {
	case blockTypeString:
		err = b.endString(l)
	case blockTypeMComment:
		rs, err = b.endMComments(l)
	case blockTypeSComment:
		rs, err = b.endSComments(l)
	}
	return rs, err
}

// 从 l 的当前位置开始往后查找，直到找到 b 中定义的 end 字符串，
// 将 l 中的指针移到该位置。
// 正常找到结束符的返回 true，否则返回 false。
func (b *block) endString(l *lexer) *app.SyntaxError {
LOOP:
	for {
		switch {
		case l.atEOF():
			break LOOP
		case (len(b.Escape) > 0) && l.match(b.Escape):
			l.next()
		case l.match(b.End):
			return nil
		default:
			l.next()
		}
	} // end for
	return l.syntaxError("未找到字符串的结束标记")
}

// 从 l 的当前位置往后开始查找连续的相同类型单行代码块。
func (b *block) endSComments(l *lexer) ([]rune, *app.SyntaxError) {
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

	return ret, nil
}

// 从 l 的当前位置一直到定义的 b.End 之间的所有字符。
// 会对每一行应用 filterSymbols 规则。
func (b *block) endMComments(l *lexer) ([]rune, *app.SyntaxError) {
	lines := make([][]rune, 0, 20)
	line := make([]rune, 0, 100)

LOOP:
	for {
		switch {
		case l.atEOF():
			return nil, l.syntaxError("未找到注释的结束标记:" + b.End)
		case l.match(b.End):
			lines = append(lines, filterSymbols(line))
			break LOOP
		default:
			r := l.next()
			line = append(line, r)
			if r == '\n' {
				lines = append(lines, filterSymbols(line))
				line = make([]rune, 0, 100)
			}
		}
	}

	ret := make([]rune, 0, 1000)
	for _, v := range lines {
		ret = append(ret, v...)
	}
	return ret, nil
}

// 行首若出现`空白字符+symbol+空白字符`的组合，则去掉这些字符。
// symbol 为全局变量 app.Symbols 中定义的任意字符。
func filterSymbols(line []rune) []rune {
	for k, v := range line {
		if unicode.IsSpace(v) { // 跳过行首的空格
			continue
		}

		// 不存在指定的符号，直接返回原数据
		if strings.IndexRune(app.Symbols, v) < 0 {
			return line
		}

		// 若下个字符正好是是空格
		if len(line) > k+1 && unicode.IsSpace(line[k+1]) {
			return line[k+2:]
		}
	}

	return line
}
