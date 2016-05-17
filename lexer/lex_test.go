// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lexer

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLexer_lineNumber(t *testing.T) {
	a := assert.New(t)
	l := New([]rune("\n\n"))
	a.NotNil(l)

	a.Equal(0, l.lineNumber())

	l.pos++
	a.Equal(1, l.lineNumber())

	l.pos++
	a.Equal(2, l.lineNumber())
}

func TestLexer_Read(t *testing.T) {
	a := assert.New(t)

	l := New([]rune("line1\n @api line2 \n"))
	a.NotNil(l)

	word := l.Read("@api")
	a.Equal(word, "line1")

	l.Match("@api")
	word = l.Read("@api")
	a.Equal(word, "line2")
}

func TestLexer_ReadN(t *testing.T) {
	a := assert.New(t)

	l := New([]rune("line1\n @api line2 \n"))
	a.NotNil(l)

	words, err := l.ReadN(1, "@api")
	a.NotError(err).Equal(words, []string{"line1"})

	l.Match("@api")
	words, err = l.ReadN(1, "@api") // 行尾并没有@api,匹配eof
	a.NotError(err).Equal(words, []string{"line2"})

	// 多词匹配
	l = New([]rune("word1 word2 word3 word4\n @api word5 word6 \n"))
	words, err = l.ReadN(2, "\n")
	a.NotError(err).Equal(words, []string{"word1", "word2 word3 word4"})

	l.Match("@api")
	words, err = l.ReadN(5, "\n")
	a.Error(err)

	l = New([]rune("word1 word2 word3 word4\n"))
	words, err = l.ReadN(1, "\n")
	a.NotError(err).Equal(words, []string{"word1 word2 word3 word4"})
}

func TestLexer_Match(t *testing.T) {
	a := assert.New(t)
	l := New([]rune("line1\n line2 \n"))
	a.NotNil(l)

	a.True(l.Match("line"))
	a.Equal('1', l.data[l.pos])
	l.pos++

	l.pos++ // \n
	l.pos++ // 空格
	l.pos++ // l

	a.False(l.Match("2222")) // 不匹配，不会移动位置
	a.True(l.Match("ine2"))  // 正确匹配
	l.Backup()
	l.Backup()
	a.Equal('i', l.data[l.pos])
	l.pos++

	// 超过剩余字符的长度。
	a.False(l.Match("ne2\n\n"))
}
