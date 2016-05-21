// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// lexer 提供基本的代码解析功能。
package lexer

import "unicode"

type Lexer struct {
	data  []rune
	pos   int // 当前指针位置
	width int // 最后移的字符数量
}

// 声明一个新的 Lexer 实例。
func New(data []rune) *Lexer {
	// TODO(caixw) Lexer 会大量产生，将其封装到 sync.Pool 是否对性能有一定提升。
	return &Lexer{
		data: data,
	}
}

// 当前位置在源代码中的行号，起始行为 0
func (l *Lexer) lineNumber() int {
	count := 0
	for i := 0; i < l.pos; i++ {
		if l.data[i] == '\n' {
			count++
		}
	}

	return count
}

// 构建一个语法错误的信息。
func (l *Lexer) SyntaxError(msg string) *SyntaxError {
	return &SyntaxError{
		Line:    l.lineNumber(),
		Message: msg,
	}
}

// 判断接下去的几个字符连接起来是否正好为 word。
// 若是，则移动指针到 word 之后，且返回 true；否则不移动指针，返回 false。
//
// NOTE: 可通过 Backup 来撤消最后一次 Match 调用。
func (l *Lexer) Match(word string) bool {
	if l.pos+len(word) >= len(l.data) { // 剩余字符没有word长，直接返回false
		return false
	}

	width := 0
	for _, r := range word {
		rr := l.data[l.pos]
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

// 撤消 Match 函数的最后次调用。指针指向执行这些函数之前的位置。
func (l *Lexer) Backup() {
	l.pos -= l.width
	l.width = 0
}

// 读取从当前位置到 delimiter 之间的所有字符，不包含前导空格，不包含 delimiter。
// delimiter 字符串返回未读字符串中。
func (l *Lexer) Read(delimiter string) []rune {
	l.SkipSpace()

	start := l.pos
	for {
		if l.pos >= len(l.data) || l.Match(delimiter) {
			l.Backup() // 只有 Match() 会触发 Backup()，EOF 不会发生任何事情
			break
		}
		l.pos++
	}
	return trimRight(l.data[start:l.pos])
}

// 往后读取，真到碰到第一个空字符或是结尾。返回字符串去掉首尾空字符。
func (l *Lexer) ReadWord() []rune {
	l.SkipSpace()

	start := l.pos
	for {
		if l.pos >= len(l.data) || unicode.IsSpace(l.data[l.pos]) {
			break
		}
		l.pos++
	}
	return l.data[start:l.pos]
}

// 往后读取一行内容，不包含首尾空格。
func (l *Lexer) ReadLine() []rune {
	l.SkipSpace()

	start := l.pos
	for {
		if l.pos >= len(l.data) || l.data[l.pos] == '\n' {
			break
		}
		l.pos++
	}
	return trimRight(l.data[start:l.pos])
}

// 跳过之后的空白字符。
func (l *Lexer) SkipSpace() {
	for {
		if l.pos >= len(l.data) || !unicode.IsSpace(l.data[l.pos]) {
			return
		}

		l.pos++
	}
}

// 指针移向下一个字符
func (l *Lexer) Next() {
	if !l.AtEOF() {
		l.pos++
	}
}

// 是否已经到结尾
func (l *Lexer) AtEOF() bool {
	return l.pos >= len(l.data)
}

// 去掉首尾的空格
func trimRight(data []rune) []rune {
	end := len(data) - 1
	for ; end >= 0; end-- {
		if !unicode.IsSpace(data[end]) {
			break
		}
	}

	return data[:end+1]
}
