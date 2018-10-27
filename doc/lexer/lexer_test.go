// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lexer

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/input"
)

func newLexerString(data string) *Lexer {
	return New(input.Block{Data: []byte(data)})
}

func newTagString(data string) *Tag {
	return &Tag{
		Data: []byte(data),
	}
}

func TestLexer_Tag(t *testing.T) {
	a := assert.New(t)
	l := newLexerString(`@api get /path desc
markdown desc line1
markdown desc line2
   @apigroup xxx
 @apitags t1,t2`)

	// @api
	tag := l.Tag()
	a.NotNil(tag)
	a.Equal(tag.Line, 0).
		Equal(string(tag.Data), `get /path desc
markdown desc line1
markdown desc line2`).Equal(tag.Name, "@api")

	// @apigroup
	tag = l.Tag()
	a.NotNil(tag)
	a.Equal(tag.Line, 3).
		Equal(string(tag.Data), "xxx").
		Equal(tag.Name, "@apigroup")

	// @apitags
	tag = l.Tag()
	a.NotNil(tag)
	a.Equal(tag.Line, 4).
		Equal(string(tag.Data), "t1,t2").
		Equal(tag.Name, "@apitags")

	// 没有标签了
	tag = l.Tag()
	a.Nil(tag)

	// 没有标签了，多次调用，结果是一样的
	tag = l.Tag()
	a.Nil(tag)
}

func TestSplitWords(t *testing.T) {
	a := assert.New(t)

	tag := []byte("@tag s1\ts2  \n  s3")

	bs := SplitWords(tag, 1)
	a.Equal(bs, [][]byte{[]byte("@tag s1\ts2  \n  s3")})

	bs = SplitWords(tag, 2)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1\ts2  \n  s3")})

	bs = SplitWords(tag, 3)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2  \n  s3")})

	bs = SplitWords(tag, 4)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2"), []byte("s3")})

	// 不够
	bs = SplitWords(tag, 5)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2"), []byte("s3")})
}

func TestSplitLines(t *testing.T) {
	a := assert.New(t)

	tag := []byte("@tag s1\ts2  \n  s3")

	bs := splitLines(tag, 1)
	a.Equal(bs, [][]byte{[]byte("@tag s1\ts2  \n  s3")})

	bs = splitLines(tag, 2)
	a.Equal(bs, [][]byte{[]byte("@tag s1\ts2  "), []byte("  s3")})

	// 不够
	bs = splitLines(tag, 3)
	a.Equal(bs, [][]byte{[]byte("@tag s1\ts2  "), []byte("  s3")})
}
