// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// lexer 提供基本的代码解析功能。
package lexer

import (
	"strings"
	"unicode"
)

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

// 当前位置在源代码中的行号
func (l *Lexer) lineNumber() int {
	count := 0
	for i := 0; i < l.pos; i++ {
		if l.data[i] == '\n' {
			count++
		}
	}
	return count
}

// 返回一个语法错误的error接口。
func (l *Lexer) SyntaxError(msg string) *SyntaxError {
	return &SyntaxError{
		Line:    l.lineNumber(),
		Message: msg,
	}
}

// 读取从当前位置到 delimiter 之间的所有字符,会去掉尾部空格
func (l *Lexer) Read(delimiter string) string {
	rs := []rune{}

	for {
		if l.pos >= len(l.data) || l.Match(delimiter) { // EOF或是到了下个标签处
			if delimiter != "\n" {
				l.Backup() // 若是eof，backup不会发生任何操作
			}
			break
		}
		rs = append(rs, l.data[l.pos])
		l.pos++
	} // end for

	return strings.TrimSpace(string(rs))
}

// 读取从当前位置到delimiter之间的所有内容，并按空格分成n个数组。
func (l *Lexer) ReadN(n int, delimiter string) ([]string, error) {
	ret := make([]string, 0, n)
	size := 0
	rs := []rune{}

	for {
		if l.pos >= len(l.data) || l.Match(delimiter) { // EOF或是到了下个标签处
			if delimiter != "\n" {
				l.Backup() // 若是eof，backup不会发生任何操作
			}

			if len(rs) > 0 {
				// 最后一条数据，去掉尾部空格
				ret = append(ret, strings.TrimRightFunc(string(rs), unicode.IsSpace))
			}
			break
		}

		r := l.data[l.pos]
		l.pos++
		if unicode.IsSpace(r) {
			if len(rs) == 0 { // 多个连续空格
				continue
			}
			if size < n-1 {
				ret = append(ret, string(rs))
				rs = rs[:0]
				size++
				continue
			}
		}

		rs = append(rs, r)
	} // end for

	if len(ret) < n {
		return nil, l.SyntaxError("未指定足够的参数")
	}
	return ret, nil
}

// 判断接下去的几个字符连接起来是否正好为word，若不匹配，则不移动指针。
// 可通过Lexer.backup来撤消最后一次调用。
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

// 撤消match函数的最后次调用。指针指向执行这些函数之前的位置。
func (l *Lexer) Backup() {
	l.pos -= l.width
	l.width = 0
}
