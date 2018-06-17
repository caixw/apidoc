// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package parser

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/input"
)

func newLexerString(data string) *lexer {
	return newLexer(input.Block{Data: []byte(data)})
}

func newTagString(data string) *tag {
	return &tag{
		data: []byte(data),
	}
}

func TestLexer_tag(t *testing.T) {
	a := assert.New(t)
	l := newLexerString(`@api get /path desc
markdown desc line1
markdown desc line2
   @apigroup xxx
 @apitags t1,t2`)

	tag, eof := l.tag()
	a.NotNil(tag).False(eof)
	a.Equal(tag.ln, 0).
		Equal(string(tag.data), `get /path desc
markdown desc line1
markdown desc line2
`).Equal(tag.name, "@api")

	tag, eof = l.tag()
	a.NotNil(tag).False(eof)
	a.Equal(tag.ln, 3).
		Equal(string(tag.data), "xxx\n").
		Equal(tag.name, "@apigroup")

	tag, eof = l.tag()
	a.NotNil(tag).True(eof)
	a.Equal(tag.ln, 4).
		Equal(string(tag.data), "t1,t2").
		Equal(tag.name, "@apitags")

	tag, eof = l.tag()
	a.Nil(tag).True(eof)
}

func TestSplit(t *testing.T) {
	a := assert.New(t)

	tag := []byte("@tag s1\ts2  \n  s3")

	bs := split(tag, 1)
	a.Equal(bs, [][]byte{[]byte("@tag s1\ts2  \n  s3")})

	bs = split(tag, 2)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1\ts2  \n  s3")})

	bs = split(tag, 3)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2  \n  s3")})

	bs = split(tag, 4)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2"), []byte("s3")})

	// 不够
	bs = split(tag, 5)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2"), []byte("s3")})
}
