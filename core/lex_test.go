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
	a.Equal(" line2 ", l.nextLine()) // 空格会被过滤
	l.backup()
	a.Equal(" line2 ", l.nextLine())

	a.Equal("", l.nextLine()) // 没有更多内容了
}

func TestLex_nextWord(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]byte("word1 word2\nword3"), 100, "file.go")
	a.NotNil(l)

	l.next()               // w
	l.next()               // o
	w, eol := l.nextWord() // rd1
	a.False(eol).Equal(w, "rd1")
	l.backup() // 多次调用backup，只启作用一次。
	l.backup()
	a.Equal(l.next(), 'r')

	w, eol = l.nextWord()
	a.False(eol).Equal(w, "d1")

	// 对空格的操作，不返回任何值。
	w, eol = l.nextWord()
	w, eol = l.nextWord()
	w, eol = l.nextWord()
	a.False(eol).Equal(w, "")

	// 第2个单词
	l.skipSpace()
	w, eol = l.nextWord() // word2
	a.True(eol).Equal(w, "word2")

	// eol，不会查找下一行的内容
	w, eol = l.nextWord()
	a.True(eol).Equal(w, "")
	// eol，不会查找下一行的内容
	w, eol = l.nextWord()
	a.True(eol).Equal(w, "")

	// 跳到下一行,eol
	l.next()
	w, eol = l.nextWord()
	a.True(eol).Equal(w, "word3")

	// eol,没有再多的内容了
	w, eol = l.nextWord()
	a.True(eol).Equal(w, "")
}

func TestLex_match(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]byte("line1\n line2 \n"), 100, "file.go")
	a.NotNil(l)

	a.True(l.match("line"))
	a.Equal('1', l.next())

	l.next() // \n
	l.next() // 空格
	l.next() // l

	a.False(l.match("2222")) // 不匹配，不会移动位置
	a.True(l.match("ine2"))  // 正确匹配
	l.backup()
	l.backup()
	a.Equal('i', l.next())

	// 超过剩余字符的长度。
	a.False(l.match("ne2\n\n"))
}

func TestLex_skipSpace(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]byte("  ln  1  \n 2 \n"), 100, "file.go")
	a.NotNil(l)

	l.skipSpace()
	a.Equal('l', l.next())
	l.next() // n

	l.skipSpace()
	l.backup() // lexer.backup对lexer.skipSpace()不启作用
	a.Equal('1', l.next())

	l.skipSpace() // 不能跳过\n
	l.skipSpace() // 不能跳过\n
	a.Equal('\n', l.next())

	l.skipSpace()
	a.Equal('2', l.next())

	l.skipSpace()
	a.Equal('\n', l.next())
	l.next()
	l.next()
	a.Equal(eof, l.next())

	// 文件结尾
	l.skipSpace()
	a.Equal(eof, l.next())
}

func TestLex_scanApiExample(t *testing.T) {
	a := assert.New(t)

	// 正常测试
	code := ` xml
<root>
    <data>123</data>
</root>`
	matchCode := `
<root>
    <data>123</data>
</root>`
	l := newLexer([]byte(code), 100, "file.go")
	a.NotNil(l)
	e, err := l.scanApiExample()
	a.NotError(err).
		Equal(e.Type, "xml").
		Equal(e.Code, matchCode)

	code = ` xml <root>
    <data>123</data>
</root>
 @apiURL abc/test`
	matchCode = `<root>
    <data>123</data>
</root>`
	l = newLexer([]byte(code), 100, "file.go")
	a.NotNil(l)
	e, err = l.scanApiExample()
	a.NotError(err).
		Equal(e.Type, "xml").
		Equal(len(e.Code), len(matchCode)).
		Equal(e.Code, matchCode)
}
