// Copyright 2016 t. caixw, All rights reserved.
// Use of this source code is governed t. a MIT
// license that can t. found in the LICENSE file.

package syntax

import (
	"unicode"

	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/types"
	"github.com/caixw/apidoc/vars"
)

// 简单的词法分析
type lexer struct {
	input *Input
	pos   int // 当前指针位置
	width int // 最后移的字符数量
}

// 表示一个标签的内容
type tag struct {
	lexer *lexer
	ln    int // 当前 tag 在文件中的起始行号
	data  []rune
	pos   int
}

// 声明一个新的 lexer 实例。
func newLexer(input *Input) *lexer {
	// TODO(caixw) lexer 会大量产生，将其封装到 sync.Pool 是否对性能有一定提升。
	return &lexer{
		input: input,
	}
}

// 当前位置在源代码文件中的行号
func (l *lexer) lineNumber() int {
	count := 0
	for i := 0; i < l.pos; i++ {
		if l.input.Data[i] == '\n' {
			count++
		}
	}

	return count + l.input.Line
}

// 输出一条错误信息
func (l *lexer) syntaxError(format string, v ...interface{}) {
	if l.input.Error == nil {
		return
	}

	err := &types.SyntaxError{
		File:    l.input.File,
		Line:    l.lineNumber(),
		Message: locale.Sprintf(format, v...),
	}

	l.input.Error.Println(err)
}

// 输出一条警告信息
func (l *lexer) syntaxWarn(format string, v ...interface{}) {
	if l.input.Warn == nil {
		return
	}

	err := &types.SyntaxError{
		File:    l.input.File,
		Line:    l.lineNumber(),
		Message: locale.Sprintf(format, v...),
	}

	l.input.Warn.Println(err)
}

// 判断接下去的几个字符连接起来是否正好为 word，且处在行首位置(word 之前不能有非空白字符)。
// 若是，则移动指针到 word 之后，且返回 true；否则不移动指针，返回 false。
//
// NOTE: 可通过 backup 来撤消最后一次 match 调用。
func (l *lexer) match(word string) bool {
	if l.atEOF() ||
		(l.pos+len(word) > len(l.input.Data)) || // 剩余字符没有 word 长
		!l.prefixIsSpace() {
		return false
	}

	width := 0
	for _, r := range word {
		rr := l.input.Data[l.pos]
		if unicode.ToLower(rr) != unicode.ToLower(r) {
			l.pos -= width
			return false
		}

		l.pos++
		width++
	}

	l.width = width
	return true
}

// 当前行在当前位置之前，只有空白字符。
func (l *lexer) prefixIsSpace() bool {
	if l.pos <= 0 {
		return true
	}

	pos := l.pos
	for {
		pos--
		r := l.input.Data[pos]

		if l.input.Data[pos] == '\n' || pos == 0 {
			return true
		}

		if !unicode.IsSpace(r) {
			return false
		}
	}
}

// 接下的单词是否和一个标签匹配。
func (l *lexer) matchTag(tagName string) bool {
	if !l.match(tagName) {
		return false
	}

	if !unicode.IsSpace(l.input.Data[l.pos]) {
		l.backup()
		return false
	}

	return true
}

// 撤消 match 函数的最后次调用。指针指向执行这些函数之前的位置。
func (l *lexer) backup() {
	l.pos -= l.width
	l.width = 0
}

// 将当前内容到下一个标签之前的内容，组成一个 tag 实例。
func (l *lexer) readTag() *tag {
	l.skipSpace()

	start := l.pos
	ln := l.lineNumber() // 记录行号
	l.width = 0          // 防止外层已经调用 match
	for {
		if l.atEOF() || l.match(vars.API) { // 直到碰到下个标签或是结束
			l.backup() // 退回标签本身的字符串
			break
		}
		l.pos++
	}

	return &tag{
		lexer: l,
		data:  trimRight(l.input.Data[start:l.pos]),
		ln:    ln,
	}
}

// 读取从当前位置到结尾的所有内容
func (l *lexer) readEnd() string {
	if l.atEOF() {
		return ""
	}

	start := l.pos
	l.pos = len(l.input.Data)
	return string(l.input.Data[start:])
}

// 跳过之后的空白字符。
func (l *lexer) skipSpace() {
	for {
		if l.atEOF() || !unicode.IsSpace(l.input.Data[l.pos]) {
			return
		}

		l.pos++
	}
}

// 往后读取一行内容，不包含首尾空格。
func (l *lexer) readLine() string {
	l.skipSpace()

	start := l.pos
	for {
		if l.atEOF() || l.input.Data[l.pos] == '\n' {
			break
		}
		l.pos++
	}
	return string(trimRight(l.input.Data[start:l.pos]))
}

// 往后读取，真到碰到第一个空字符或是结尾。返回字符串去掉首尾空字符。
func (l *lexer) readWord() string {
	l.skipSpace()

	start := l.pos
	for {
		if l.atEOF() || unicode.IsSpace(l.input.Data[l.pos]) {
			break
		}
		l.pos++
	}
	return string(trimRight(l.input.Data[start:l.pos]))
}

// 是否已经到结尾
func (l *lexer) atEOF() bool {
	return l.pos >= len(l.input.Data)
}

// 往后读取，真到碰到第一个空字符或是结尾。返回字符串去掉首尾空字符。
func (t *tag) readWord() string {
	t.skipSpace()

	start := t.pos
	for {
		if t.atEOF() || unicode.IsSpace(t.data[t.pos]) {
			break
		}
		t.pos++
	}
	return string(trimRight(t.data[start:t.pos]))
}

// 当前位置在源代码文件中的行号
func (t *tag) lineNumber() int {
	count := t.ln
	for i := 0; i < t.pos; i++ {
		if t.data[i] == '\n' {
			count++
		}
	}

	return count
}

// 输出语法错误
func (t *tag) syntaxError(format string, v ...interface{}) {
	if t.lexer.input.Error == nil {
		return
	}

	err := &types.SyntaxError{
		File:    t.lexer.input.File,
		Line:    t.lineNumber(),
		Message: locale.Sprintf(format, v...),
	}
	t.lexer.input.Error.Println(err)
}

// 输出语法警告信息
func (t *tag) syntaxWarn(format string, v ...interface{}) {
	if t.lexer.input.Warn == nil {
		return
	}

	err := &types.SyntaxError{
		File:    t.lexer.input.File,
		Line:    t.lineNumber(),
		Message: locale.Sprintf(format, v...),
	}
	t.lexer.input.Warn.Println(err)
}

// 往后读取一行内容，不包含首尾空格。
func (t *tag) readLine() string {
	t.skipSpace()

	start := t.pos
	for {
		if t.atEOF() || t.data[t.pos] == '\n' {
			break
		}
		t.pos++
	}
	return string(trimRight(t.data[start:t.pos]))
}

// 读取从当前位置到结尾的所有内容，去掉首尾空格
func (t *tag) readEnd() string {
	if t.atEOF() {
		return ""
	}

	t.skipSpace()
	start := t.pos
	t.pos = len(t.data)
	return string(t.data[start:]) // 不用 trimRight，已经在初始化时去掉尾部的空格
}

// 是否在结尾处
func (t *tag) atEOF() bool {
	return t.pos >= len(t.data)
}

// 跳过之后的空白字符。
func (t *tag) skipSpace() {
	for {
		if t.atEOF() || !unicode.IsSpace(t.data[t.pos]) {
			return
		}

		t.pos++
	}
}

// 去掉尾部空格
func trimRight(data []rune) []rune {
	end := len(data) - 1
	for ; end >= 0; end-- {
		if !unicode.IsSpace(data[end]) {
			break
		}
	}

	return data[:end+1]
}
