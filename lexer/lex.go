// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// lexer 提供基本的代码解析功能。
package lexer

type Lexer struct {
	data  []rune
	pos   int // 当前指针位置
	width int // 最后移的字符数量
}

// 声明一个新的 Lexer 实例。
func New(data []rune) *Lexer {
	// TODO(caixw) Lexer 会大量产生，将其封闭到 sync.Pool 是否对性能有一定提升。
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
// NOTE: 可通过 Backup 来撤消最后一次调用。
func (l *Lexer) Match(word string) bool {
	if l.pos+len(word) >= len(l.data) { // 剩余字符没有word长，直接返回false
		return false
	}

	width := 0
	for _, r := range word {
		rr := l.data[l.pos]
		if rr != r {
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

// 读取从当前位置到 delimiter 之间的所有字符。
//
// NOTE: 不包含 delimiter 字符串本身，该字符串会返回未读内容中。
func (l *Lexer) Read(delimiter string) []rune {
	rs := []rune{} // TODO 精简掉此内存分配

	for {
		if l.pos >= len(l.data) || l.Match(delimiter) {
			l.Backup() // 只有 Match() 会触发 Backup()，EOF 不会发生任何事情
			break
		}

		rs = append(rs, l.data[l.pos])
		l.pos++
	} // end for

	return rs
}

func (l *Lexer) ReadLine() []rune {
	return l.Read("\n")
}
