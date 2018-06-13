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

	tag := l.tag()
	a.NotNil(tag)
	a.Equal(tag.ln, 0).
		Equal(string(tag.data), `@api get /path desc
markdown desc line1
markdown desc line2
`)

	tag = l.tag()
	a.NotNil(tag)
	a.Equal(tag.ln, 3).
		Equal(string(tag.data), "@apigroup xxx\n")

	tag = l.tag()
	a.NotNil(tag)
	a.Equal(tag.ln, 4).
		Equal(string(tag.data), "@apitags t1,t2")
}

func TestTag_split(t *testing.T) {
	a := assert.New(t)

	tag := newTagString("@tag s1\ts2  \n  s3")

	bs := tag.split(1)
	a.Equal(bs, [][]byte{[]byte("@tag s1\ts2  \n  s3")})

	bs = tag.split(2)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1\ts2  \n  s3")})

	bs = tag.split(3)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2  \n  s3")})

	bs = tag.split(4)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2"), []byte("s3")})

	// 不够
	bs = tag.split(5)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2"), []byte("s3")})
}
