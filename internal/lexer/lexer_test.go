// SPDX-License-Identifier: MIT

package lexer

import (
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/caixw/apidoc/v7/core"
)

func TestNewLexer(t *testing.T) {
	a := assert.New(t, false)

	l, err := New(core.Block{Data: []byte("// doc")})
	a.NotError(err).NotNil(l)

	// 中间三个字节表示一个汉字
	l, err = New(core.Block{Data: []byte{0b01001001, 0b11100111, 0b10011111, 0b10100101, 0b01110100}})
	a.NotError(err).NotNil(l)

	// 中间三个字节表示一个汉字，编码错误
	l, err = New(core.Block{Data: []byte{0b01001001, 0b01100111, 0b10011111, 0b10100101, 0b01110100}})
	a.Error(err).Nil(l)
}

func TestLexer_Position(t *testing.T) {
	a := assert.New(t, false)

	loc := core.Location{
		Range: core.Range{
			Start: core.Position{Line: 11, Character: 2},
		},
	}
	l, err := New(core.Block{Data: []byte("l0\nl1\nl2\nl3\n"), Location: loc})
	a.NotError(err).NotNil(l)

	l.Next(3)
	a.Equal(l.Current(), Position{
		Position: core.Position{Line: 12, Character: 0},
		Offset:   3,
	})

	l.Next(4)
	a.Equal(l.Current(), Position{
		Position: core.Position{Line: 13, Character: 1},
		Offset:   7,
	})

	l.Next(3)
	l.Next(3)
	a.Equal(l.Current(), Position{
		Position: core.Position{Line: 15, Character: 0},
		Offset:   12,
	})

	l, err = New(core.Block{Data: []byte("12中文ab"), Location: loc})
	a.NotError(err).NotNil(l)
	l.Next(2)
	a.Equal(l.Current(), Position{
		Position: core.Position{Line: 11, Character: 4},
		Offset:   2,
	})

	l.Next(2)
	a.Equal(l.Current(), Position{
		Position: core.Position{Line: 11, Character: 6},
		Offset:   8,
	})

	l.Next(10)
	a.True(l.AtEOF())

	p := Position{
		Position: core.Position{Line: 0, Character: 2},
		Offset:   2,
	}
	l.Move(p)
	a.False(l.AtEOF()).
		Equal(-1, l.prev.Offset).
		Equal(l.current, p)

	l.Move(Position{Offset: 11111111})
	a.True(l.AtEOF())

	// 不能为负
	a.Panic(func() {
		l.Move(Position{Offset: -1})
	})
}

func TestLexer_Match(t *testing.T) {
	a := assert.New(t, false)

	l, err := New(core.Block{Data: []byte("ab中\ncd")})
	a.NotError(err).NotNil(l)

	a.False(l.Match("")).Equal(0, l.current.Offset)
	a.False(l.Match("b")).Equal(0, l.current.Offset)
	a.True(l.Match("ab")).Equal(2, l.current.Offset)
	a.False(l.Match("ab")).Equal(2, l.current.Offset)
	a.True(l.Match("中")).Equal(l.current, Position{
		Position: core.Position{Line: 0, Character: 3},
		Offset:   5,
	})

	l.Rollback()
	l.Rollback()
	l.Rollback()
	a.True(l.Match("中")).Equal(l.current, Position{
		Position: core.Position{Line: 0, Character: 3},
		Offset:   5,
	})

	l.Next(len(l.Data))
	a.False(l.Match("ab")).Equal(l.current, Position{
		Position: core.Position{Line: 1, Character: 2},
		Offset:   8,
	})
}

func TestLexer_Spaces(t *testing.T) {
	a := assert.New(t, false)
	l, err := New(core.Block{Data: []byte("    0 \n  1 ")})
	a.NotError(err).NotNil(l)

	l.Spaces('\n')
	a.Equal(string(l.Next(1)), "0")

	// 无法跳过换行符
	l.Next(1)
	l.Spaces('\n')
	l.Spaces('\n')
	l.Spaces('\n')
	l.Spaces('\n')
	a.Empty(l.Spaces('\n'))
	a.Equal(string(l.Next(1)), "\n")

	a.Equal(l.Spaces('\t'), "  ")
	l.Rollback()
	a.Equal(2, len(l.Spaces(0)))

	l.Rollback()
	a.Equal(0, len(l.Spaces(' ')))
	a.Equal(2, len(l.Spaces(0)))

	a.Equal(string(l.Next(1)), "1")

	l.Next(1)
	l.Spaces('\n')
	l.Spaces('\n')
	a.Equal(l.current.Offset, len(l.Data))
}

func TestLexer_Delim(t *testing.T) {
	a := assert.New(t, false)

	l, err := New(core.Block{Data: []byte("123")})
	a.NotError(err).NotNil(l)
	val, found := l.Delim('\n', true)
	a.False(found).Nil(val)
	val, found = l.Delim(0, true)
	a.False(found).Nil(val)

	l, err = New(core.Block{Data: []byte("123\n")})
	a.NotError(err).NotNil(l)
	val, found = l.Delim('\n', true)
	a.True(found).Equal(string(val), "123\n").
		Equal(l.current.Offset, 4)

	l = &Lexer{Block: core.Block{Data: []byte("123\n")}, current: Position{Offset: 1}}
	val, found = l.Delim('\n', true)
	a.True(found).Equal(string(val), "23\n").
		Equal(l.current.Offset, 4)
	l.Rollback()
	a.Equal(l.current.Offset, 1)
	val, found = l.Delim('\n', true)
	a.True(found).Equal(string(val), "23\n").
		Equal(l.current.Offset, 4)
}

func TestLexer_DelimFunc(t *testing.T) {
	a := assert.New(t, false)

	l, err := New(core.Block{Data: []byte("123456789\n123456789\n123")})
	a.NotError(err).NotNil(l)

	// 不存在
	val, found := l.DelimFunc(func(r rune) bool { return r == 0 }, true)
	a.False(found).Nil(val)
	val, found = l.DelimFunc(func(r rune) bool { return r == '\r' }, true)
	a.False(found).Nil(val)
	a.Panic(func() {
		l.DelimFunc(nil, false)
	})

	bs, found := l.DelimFunc(func(r rune) bool { return r == '9' }, true)
	a.True(found).Equal(string(bs), "123456789")

	bs, found = l.DelimFunc(func(r rune) bool { return r == '9' }, false)
	a.True(found).Equal(string(bs), "\n12345678")

	bs, found = l.DelimFunc(func(r rune) bool { return r == '9' }, false)
	a.True(found).Empty(bs)

	bs, found = l.DelimFunc(func(r rune) bool { return r == '9' }, true)
	a.True(found).Equal(string(bs), "9")

	bs, found = l.DelimFunc(func(r rune) bool { return r == '\n' }, false)
	a.True(found).Empty(bs)
	a.Equal(l.Current(), Position{Offset: 19, Position: core.Position{Line: 1, Character: 0}})

	bs, found = l.DelimFunc(func(r rune) bool { return r == '\n' }, true)
	a.True(found).Equal(string(bs), "\n")
	a.Equal(l.Current(), Position{Offset: 20, Position: core.Position{Line: 2, Character: 0}})

	bs, found = l.DelimFunc(func(r rune) bool { return r == '3' }, false)
	a.True(found).Equal(string(bs), "12")
	a.False(l.AtEOF())
	bs, found = l.DelimFunc(func(r rune) bool { return r == '3' }, false)
	a.True(found).Empty(bs)
	a.False(l.AtEOF())
	bs, found = l.DelimFunc(func(r rune) bool { return r == '3' }, true)
	a.True(found).Equal(string(bs), "3")
	a.True(l.AtEOF())
}

func TestLexer_All(t *testing.T) {
	a := assert.New(t, false)

	l, err := New(core.Block{Data: []byte("123")})
	a.NotError(err).NotNil(l)
	a.Equal(l.All(), "123").True(l.AtEOF())

	l, err = New(core.Block{Data: []byte("123\n456")})
	a.NotError(err).NotNil(l)
	l.Next(1)
	a.Equal(l.All(), "23\n456").True(l.AtEOF())

	l.Rollback()
	a.Equal(l.Current().Offset, 1)
	a.Equal(l.All(), "23\n456").True(l.AtEOF())
}

func TestLexer_Bytes(t *testing.T) {
	a := assert.New(t, false)

	l, err := New(core.Block{Data: []byte("123")})
	a.NotError(err).NotNil(l)
	a.Equal(l.Bytes(1, 2), "2").False(l.AtEOF())
}

func TestLexer_DelimString(t *testing.T) {
	a := assert.New(t, false)

	l, err := New(core.Block{Data: []byte("1234567")})
	a.NotError(err).NotNil(l)
	a.Panic(func() { l.DelimString("", true) })

	val, found := l.DelimString("45", true)
	a.True(found).Equal(val, "12345").
		Equal(l.Current().Offset, 5)

	l.Rollback()
	val, found = l.DelimString("45", false)
	a.True(found).Equal(val, "123").
		Equal(l.Current().Offset, 3)

	val, found = l.DelimString("7", true)
	a.True(found).Equal(val, "4567").
		Equal(l.Current().Offset, 7)
	a.True(l.AtEOF(), l.current.Offset, l.lastIndex)

	l, err = New(core.Block{Data: []byte("123444567\n8910")})
	a.NotError(err).NotNil(l)

	val, found = l.DelimString("45", true)
	a.True(found).Equal(val, "1234445").
		Equal(l.Current().Offset, 7)
	l.Rollback()
	a.Equal(0, l.Current().Offset)
	val, found = l.DelimString("45", true)
	a.True(found).Equal(val, "1234445").
		Equal(l.Current().Offset, 7)

	val, found = l.DelimString("891", true)
	a.True(found).Equal(val, "67\n891").
		Equal(l.Current().Offset, 13)

	val, found = l.DelimString("891", true)
	a.False(found).Nil(val).Equal(l.Current().Offset, 13)

	val, found = l.DelimString("891", true)
	a.False(found).Nil(val).Equal(l.Current().Offset, 13)
}
