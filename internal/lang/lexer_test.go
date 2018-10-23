// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lang

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLexer_lineNumber(t *testing.T) {
	a := assert.New(t)

	l := &Lexer{data: []byte("l0\nl1\nl2\nl3\n")}
	l.pos = 3
	a.Equal(l.lineNumber(), 1)

	l.pos += 3
	a.Equal(l.lineNumber(), 2)

	l.pos += 3
	l.pos += 3
	a.Equal(l.lineNumber(), 4)
}

func TestLexer_Match(t *testing.T) {
	a := assert.New(t)

	l := &Lexer{
		data: []byte("ab\ncd"),
	}

	a.False(l.Match("b")).Equal(0, l.pos)
	a.True(l.Match("ab")).Equal(2, l.pos)

	l.pos = len(l.data)
	a.False(l.Match("ab"))

	// 匹配结尾单词
	l.pos = 3 // c的位置
	a.True(l.Match("cd"))
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
		Blocks: blocks,
	}

	b := l.Block() // scomment1
	a.Equal(b.(*block).Type, blockTypeSComment)
	rs, err := b.EndFunc(l)
	a.NotError(err).Equal(rs, [][]byte{[]byte(" scomment1\n"), []byte(" scomment2\n")})

	b = l.Block() // string1
	a.Equal(b.(*block).Type, blockTypeString)
	_, err = b.EndFunc(l)
	a.NotError(err)

	b = l.Block() // string2
	a.Equal(b.(*block).Type, blockTypeString)
	_, err = b.EndFunc(l)
	a.NotError(err)

	b = l.Block()
	a.Equal(b.(*block).Type, blockTypeMComment) // mcomment1
	rs, err = b.EndFunc(l)
	a.NotError(err).Equal(rs, [][]byte{[]byte("\n"), []byte("mcomment1\n"), []byte("mcomment2\n")})

	/* 测试一段单行注释后紧跟 \n=pod 形式的多行注释，是否会出错 */

	b = l.Block() // scomment3,scomment4
	a.Equal(b.(*block).Type, blockTypeSComment)
	rs, err = b.EndFunc(l)
	a.NotError(err).Equal(rs, [][]byte{[]byte(" scomment3\n"), []byte(" scomment4\n")})

	b = l.Block() // mcomment3,mcomment4
	a.Equal(b.(*block).Type, blockTypeMComment)
	rs, err = b.EndFunc(l)
	a.NotError(err).Equal(rs, [][]byte{[]byte("\n"), []byte(" mcomment3\n"), []byte(" mcomment4")})
}

func TestLexer_SkipSpace(t *testing.T) {
	a := assert.New(t)
	l := &Lexer{data: []byte("  0 \n  1 ")}

	l.SkipSpace()
	a.Equal(l.data[l.pos], "0")

	// 无法跳过换行符
	l.pos++
	l.SkipSpace()
	l.SkipSpace()
	l.SkipSpace()
	l.SkipSpace()
	l.SkipSpace()
	a.Equal(l.data[l.pos], "\n")

	l.pos++
	l.SkipSpace()
	a.Equal(l.data[l.pos], "1")

	l.pos++
	l.SkipSpace()
	l.SkipSpace()
	a.Equal(l.pos, len(l.data))
}
