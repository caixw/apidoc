// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLexer_LineNumber(t *testing.T) {
	a := assert.New(t)

	l := NewLexer([]byte("l0\nl1\nl2\nl3\n"), nil)
	l.pos = 3
	a.Equal(l.LineNumber(), 1)

	l.pos += 3
	a.Equal(l.LineNumber(), 2)

	l.pos += 3
	l.pos += 3
	a.Equal(l.LineNumber(), 4)
}

func TestLexer_match(t *testing.T) {
	a := assert.New(t)

	l := &Lexer{
		data: []byte("ab\ncd"),
	}

	a.False(l.match("b")).Equal(0, l.pos)
	a.True(l.match("ab")).Equal(2, l.pos)

	l.pos = len(l.data)
	a.False(l.match("ab"))

	// 匹配结尾单词
	l.pos = 3 // c的位置
	a.True(l.match("cd"))
}

func TestLexer_Block(t *testing.T) {
	a := assert.New(t)

	blocks := []Blocker{
		&block{Type: blockTypeSComment, Begin: "//"},
		&block{Type: blockTypeMComment, Begin: "/*", End: "*/"},
		&block{Type: blockTypeMComment, Begin: "\n=pod", End: "\n=cut"},
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: "\\"},
	}

	l := &Lexer{
		data: []byte(`// scomment1
// scomment2
func(){}
"/*string1"
"//string2"
/*
mcomment1
mcomment2
*/

// scomment3
// scomment4
=pod
 mcomment3
 mcomment4
=cut
`),
		blocks: blocks,
	}

	b := l.Block() // scomment1
	a.Equal(b.(*block).Type, blockTypeSComment)
	raw, data, err := b.EndFunc(l)
	a.NotError(err).
		Equal(string(raw), " scomment1\n scomment2\n").
		Equal(string(data), " scomment1\n scomment2\n")

	b = l.Block() // string1
	a.Equal(b.(*block).Type, blockTypeString)
	_, _, err = b.EndFunc(l)
	a.NotError(err)

	b = l.Block() // string2
	a.Equal(b.(*block).Type, blockTypeString)
	_, _, err = b.EndFunc(l)
	a.NotError(err)

	b = l.Block()
	a.Equal(b.(*block).Type, blockTypeMComment) // mcomment1
	raw, data, err = b.EndFunc(l)
	a.NotError(err).
		Equal(string(raw), "\nmcomment1\nmcomment2\n").
		Equal(string(data), "\nmcomment1\nmcomment2\n")

	/* 测试一段单行注释后紧跟 \n=pod 形式的多行注释，是否会出错 */

	b = l.Block() // scomment3,scomment4
	a.Equal(b.(*block).Type, blockTypeSComment)
	raw, data, err = b.EndFunc(l)
	a.NotError(err).
		Equal(string(raw), " scomment3\n scomment4\n").
		Equal(string(data), " scomment3\n scomment4\n")

	b = l.Block() // mcomment3,mcomment4
	a.Equal(b.(*block).Type, blockTypeMComment)
	raw, data, err = b.EndFunc(l)
	a.NotError(err).
		Equal(string(raw), "\n mcomment3\n mcomment4").
		Equal(string(data), "\n mcomment3\n mcomment4")
}

func TestLexer_skipSpace(t *testing.T) {
	a := assert.New(t)
	l := &Lexer{data: []byte("  0 \n  1 ")}

	l.skipSpace()
	a.Equal(l.data[l.pos], "0")

	// 无法跳过换行符
	l.pos++
	l.skipSpace()
	l.skipSpace()
	l.skipSpace()
	l.skipSpace()
	l.skipSpace()
	a.Equal(l.data[l.pos], "\n")

	l.pos++
	l.skipSpace()
	a.Equal(l.data[l.pos], "1")

	l.pos++
	l.skipSpace()
	l.skipSpace()
	a.Equal(l.pos, len(l.data))
}

func TestLexer_line(t *testing.T) {
	a := assert.New(t)

	l := NewLexer([]byte("123"), nil)
	a.Nil(l.line())

	l = NewLexer([]byte("123\n"), nil)
	a.Equal(string(l.line()), "123").
		Equal(l.pos, 3)

	l = &Lexer{data: []byte("123\n"), pos: 1}
	a.Equal(string(l.line()), "23").
		Equal(l.pos, 3)
}
