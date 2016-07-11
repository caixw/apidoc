// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"testing"
	"unicode/utf8"

	"github.com/issue9/assert"
)

func TestLexer_lineNumber(t *testing.T) {
	a := assert.New(t)

	l := &lexer{data: []byte("l0\nl1\nl2\nl3\n")}
	l.pos = 3
	a.Equal(l.lineNumber(), 1)

	l.pos += 3
	a.Equal(l.lineNumber(), 2)

	l.pos += 3
	l.pos += 3
	a.Equal(l.lineNumber(), 4)
}

func TestLexer_next(t *testing.T) {
	a := assert.New(t)

	l := &lexer{
		data: []byte("ab\ncd"),
	}

	a.Equal('a', l.next())
	a.Equal('b', l.next())
	a.Equal('\n', l.next())
	a.Equal('c', l.next())
	a.False(l.atEOF())
	a.Equal('d', l.next())
	a.True(l.atEOF())
	a.Equal(utf8.RuneError, l.next())
	a.Equal(utf8.RuneError, l.next())
	a.True(l.atEOF())
}

func TestLexer_match(t *testing.T) {
	a := assert.New(t)

	l := &lexer{
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

func TestLexer_block(t *testing.T) {
	a := assert.New(t)

	l := &lexer{
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
	}

	blocks := []blocker{
		&block{Type: blockTypeSComment, Begin: "//"},
		&block{Type: blockTypeMComment, Begin: "/*", End: "*/"},
		&block{Type: blockTypeMComment, Begin: "\n=pod", End: "\n=cut"},
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: "\\"},
	}

	b := l.block(blocks) // scomment1
	a.Equal(b.(*block).Type, blockTypeSComment)
	rs, err := b.EndFunc(l)
	a.NotError(err).Equal(string(rs), " scomment1\n scomment2\n")

	b = l.block(blocks) // string1
	a.Equal(b.(*block).Type, blockTypeString)
	_, err = b.EndFunc(l)
	a.NotError(err)

	b = l.block(blocks) // string2
	a.Equal(b.(*block).Type, blockTypeString)
	_, err = b.EndFunc(l)
	a.NotError(err)

	b = l.block(blocks)
	a.Equal(b.(*block).Type, blockTypeMComment) // mcomment1
	rs, err = b.EndFunc(l)
	a.NotError(err).Equal(string(rs), "\nmcomment1\nmcomment2\n")

	/* 测试一段单行注释后紧跟 \n=pod 形式的多行注释，是否会出错 */

	b = l.block(blocks) // scomment3,scomment4
	a.Equal(b.(*block).Type, blockTypeSComment)
	rs, err = b.EndFunc(l)
	a.NotError(err).Equal(string(rs), " scomment3\n scomment4\n")

	b = l.block(blocks) // mcomment3,mcomment4
	a.Equal(b.(*block).Type, blockTypeMComment)
	rs, err = b.EndFunc(l)
	a.NotError(err).Equal(string(rs), "\n mcomment3\n mcomment4")
}
