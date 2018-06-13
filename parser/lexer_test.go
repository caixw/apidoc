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
