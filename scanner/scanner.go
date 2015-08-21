// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

import (
	"bytes"
	"sync"
	"unicode/utf8"

	"github.com/caixw/apidoc/core"
)

const eof = -1

var (
	docs   = []*core.Doc{}
	docsMu sync.Mutex
)

// 扫描scanner中的代码，提取最近的下一个代码块和其开始的行号。
// scanFunc必须是一个无状态的
type scanFunc func(*scanner) ([]rune, int, error)

type scanner struct {
	data  []byte
	pos   int
	width int
}

// 是否已经在文件末尾。
func (s *scanner) atEOF() bool {
	return s.pos >= len(s.data)
}

// 获取当前的字符，并将指针指向下一个字符。
func (s *scanner) next() rune {
	if s.atEOF() {
		return eof
	}

	r, w := utf8.DecodeRune(s.data[s.pos:])
	s.pos += w
	s.width = w
	return r
}

// 是否匹配指定的字符串，若匹配，则将指定移向该字符串这后，否则不作任何操作。
func (s *scanner) match(str string) bool {
	rs := []rune(str)
	if s.atEOF() {
		return false
	}

	width := 0
	for _, r := range rs {
		rr, w := utf8.DecodeRune(s.data[s.pos:])
		if rr != r {
			s.pos -= width
			return false
		}

		s.pos += w
		width += w
	}

	s.width = width
	return true
}

// 撤消s.next()/s.match()的最后一次操作。
func (s *scanner) backup() {
	s.pos -= s.width
	s.width = 0
}

// 当前所在的行号
func (s *scanner) lineNumber() int {
	// 当前行应该是\n数量加1
	return bytes.Count(s.data[:s.pos], []byte("\n")) + 1
}
