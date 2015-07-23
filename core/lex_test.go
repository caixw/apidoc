// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLex_lineNumber(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]byte("\n\n"), 100, "file.go")
	a.NotNil(l)

	a.Equal(100, l.lineNumber())

	l.next()
	a.Equal(101, l.lineNumber())

	l.next()
	a.Equal(102, l.lineNumber())

	l.backup()
	a.Equal(101, l.lineNumber())
}

func TestLex_next(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]byte("ab\ncd\n"), 100, "file.go")
	a.NotNil(l)

	a.Equal('a', l.next())
	a.Equal('b', l.next())
	a.Equal('\n', l.next())
	a.Equal('c', l.next())

	// 退回一个字符
	l.backup()
	a.Equal('c', l.next())

	// 退回多个字符
	l.backup()
	l.backup()
	l.backup()
	a.Equal('c', l.next())

	a.Equal('d', l.next())
	a.Equal('\n', l.next())
	a.Equal(eof, l.next()) // 文件结束
	a.Equal(eof, l.next())
}

func TestLex_nextLine(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]byte("line1\n line2 \n"), 100, "file.go")
	a.NotNil(l)

	a.Equal("line1", l.nextLine())
	l.backup()
	l.backup()
	a.Equal('l', l.next())

	a.Equal("ine1", l.nextLine())
	a.Equal("line2", l.nextLine()) // 空格会被过滤
	l.backup()
	a.Equal("line2", l.nextLine())

	a.Equal("", l.nextLine()) // 没有更多内容了
}
