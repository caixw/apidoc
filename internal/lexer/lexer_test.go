// SPDX-License-Identifier: MIT

package lexer

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/core"
)

func TestNewLexer(t *testing.T) {
	a := assert.New(t)

	l, err := New([]byte("// doc"))
	a.NotError(err).NotNil(l)

	// 中间三个字节表示一个汉字
	l, err = New([]byte{0b01001001, 0b11100111, 0b10011111, 0b10100101, 0b01110100})
	a.NotError(err).NotNil(l)

	// 中间三个字节表示一个汉字，编码错误
	l, err = New([]byte{0b01001001, 0b01100111, 0b10011111, 0b10100101, 0b01110100})
	a.Error(err).Nil(l)
}

func TestLexer_Position(t *testing.T) {
	a := assert.New(t)

	l, err := New([]byte("l0\nl1\nl2\nl3\n"))
	a.NotError(err).NotNil(l)

	l.Next(3)
	a.Equal(l.Position(), Position{
		Position: core.Position{Line: 1, Character: 0},
		Offset:   3,
	})

	l.Next(4)
	a.Equal(l.Position(), Position{
		Position: core.Position{Line: 2, Character: 1},
		Offset:   7,
	})

	l.Next(3)
	l.Next(3)
	a.Equal(l.Position(), Position{
		Position: core.Position{Line: 4, Character: 0},
		Offset:   12,
	})

	l, err = New([]byte("12中文ab"))
	a.NotError(err).NotNil(l)
	l.Next(2)
	a.Equal(l.Position(), Position{
		Position: core.Position{Line: 0, Character: 2},
		Offset:   2,
	})

	l.Next(2)
	a.Equal(l.Position(), Position{
		Position: core.Position{Line: 0, Character: 4},
		Offset:   8,
	})

	l.Next(10)
	a.True(l.AtEOF())

	p := Position{
		Position: core.Position{Line: 0, Character: 2},
		Offset:   2,
	}
	l.SetPosition(p)
	a.False(l.AtEOF()).
		Equal(0, l.prev.Offset).
		Equal(l.current, p)
}

func TestLexer_Match(t *testing.T) {
	a := assert.New(t)

	l := &Lexer{
		data: []byte("ab中\ncd"),
	}

	a.False(l.Match("")).Equal(0, l.current.Offset)
	a.False(l.Match("b")).Equal(0, l.current.Offset)
	a.True(l.Match("ab")).Equal(2, l.current.Offset)
	a.False(l.Match("ab")).Equal(2, l.current.Offset)
	a.True(l.Match("中")).Equal(l.current, Position{
		Position: core.Position{Line: 0, Character: 3},
		Offset:   5,
	})

	l.Rollback()
	a.True(l.Match("中")).Equal(l.current, Position{
		Position: core.Position{Line: 0, Character: 3},
		Offset:   5,
	})

	l.Next(len(l.data))
	a.False(l.Match("ab")).Equal(l.current, Position{
		Position: core.Position{Line: 1, Character: 2},
		Offset:   8,
	})
}

func TestLexer_Spaces(t *testing.T) {
	a := assert.New(t)
	l, err := New([]byte("    0 \n  1 "))
	a.NotError(err).NotNil(l)

	l.Spaces()
	a.Equal(string(l.Next(1)), "0")

	// 无法跳过换行符
	l.Next(1)
	l.Spaces()
	l.Spaces()
	l.Spaces()
	l.Spaces()
	a.Empty(l.Spaces())
	a.Equal(string(l.Next(1)), "\n")

	l.Next(1)
	a.Equal(1, len(l.Spaces()))
	l.Rollback()
	a.Equal(1, len(l.Spaces()))
	a.Equal(string(l.Next(1)), "1")

	l.Next(1)
	l.Spaces()
	l.Spaces()
	a.Equal(l.current.Offset, len(l.data))
}

func TestLexer_Delim(t *testing.T) {
	a := assert.New(t)

	l, err := New([]byte("123"))
	a.NotError(err).NotNil(l)
	a.Nil(l.Delim('\n', true))
	a.Nil(l.Delim(0, true))

	l, err = New([]byte("123\n"))
	a.NotError(err).NotNil(l)
	a.Equal(string(l.Delim('\n', true)), "123\n").
		Equal(l.current.Offset, 4)

	l = &Lexer{data: []byte("123\n"), current: Position{Offset: 1}}
	a.Equal(string(l.Delim('\n', true)), "23\n").
		Equal(l.current.Offset, 4)
	l.Rollback()
	a.Equal(l.current.Offset, 1)
	a.Equal(string(l.Delim('\n', true)), "23\n").
		Equal(l.current.Offset, 4)
}

func TestLexer_DelimFunc(t *testing.T) {
	a := assert.New(t)

	l, err := New([]byte("123456789\n123456789\n123"))
	a.NotError(err).NotNil(l)
	a.Nil(l.DelimFunc(func(r rune) bool { return r == 0 }, true))    // 不存在
	a.Nil(l.DelimFunc(func(r rune) bool { return r == '\r' }, true)) // 不存在
	a.Nil(l.DelimFunc(func(r rune) bool { return r == '\r' }, true)) // 不存在

	bs := l.DelimFunc(func(r rune) bool { return r == '9' }, true)
	a.Equal(string(bs), "123456789")

	bs = l.DelimFunc(func(r rune) bool { return r == '9' }, false)
	a.Equal(string(bs), "\n12345678")

	bs = l.DelimFunc(func(r rune) bool { return r == '9' }, false)
	a.Empty(bs)

	bs = l.DelimFunc(func(r rune) bool { return r == '9' }, true)
	a.Equal(string(bs), "9")

	bs = l.DelimFunc(func(r rune) bool { return r == '\n' }, false)
	a.Empty(bs)
	a.Equal(l.Position(), Position{Offset: 19, Position: core.Position{Line: 1, Character: 0}})

	bs = l.DelimFunc(func(r rune) bool { return r == '\n' }, true)
	a.Equal(string(bs), "\n")
	a.Equal(l.Position(), Position{Offset: 20, Position: core.Position{Line: 2, Character: 0}})
}

func TestLexer_All(t *testing.T) {
	a := assert.New(t)

	l, err := New([]byte("123"))
	a.NotError(err).NotNil(l)
	a.Equal(l.All(), "123").True(l.AtEOF())

	l, err = New([]byte("123\n456"))
	a.NotError(err).NotNil(l)
	l.Next(1)
	a.Equal(l.All(), "23\n456").True(l.AtEOF())
}

func TestLexer_Bytes(t *testing.T) {
	a := assert.New(t)

	l, err := New([]byte("123"))
	a.NotError(err).NotNil(l)
	a.Equal(l.Bytes(1, 2), "2").False(l.AtEOF())
}

func TestLexer_DelimString(t *testing.T) {
	a := assert.New(t)

	l, err := New([]byte("1234567"))
	a.NotError(err).NotNil(l)
	a.Nil(l.DelimString(""))
	a.Equal(l.DelimString("45"), "12345").
		Equal(l.Position().Offset, 5)
	a.Equal(l.DelimString("7"), "67").
		Equal(l.Position().Offset, 7)
	l.Next(1)
	a.True(l.AtEOF())

	l, err = New([]byte("123444567\n8910"))
	a.NotError(err).NotNil(l)
	a.Equal(l.DelimString("45"), "1234445").
		Equal(l.Position().Offset, 7)
	a.Equal(l.DelimString("891"), "67\n891").
		Equal(l.Position().Offset, 13)
	a.Nil(l.DelimString("891")).
		Equal(l.Position().Offset, 13)
	a.Nil(l.DelimString("891")).
		Equal(l.Position().Offset, 13)
}
