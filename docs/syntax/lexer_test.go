// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package syntax

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/input"
)

func newLexerString(data string) *Lexer {
	return NewLexer(input.Block{Data: []byte(data)})
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

	tag, eof := l.Tag()
	a.NotNil(tag).False(eof)
	a.Equal(tag.Line, 0).
		Equal(string(tag.Data), `get /path desc
markdown desc line1
markdown desc line2
`).Equal(tag.Name, "@api")

	tag, eof = l.Tag()
	a.NotNil(tag).False(eof)
	a.Equal(tag.Line, 3).
		Equal(string(tag.Data), "xxx\n").
		Equal(tag.Name, "@apigroup")

	tag, eof = l.Tag()
	a.NotNil(tag).True(eof)
	a.Equal(tag.Line, 4).
		Equal(string(tag.Data), "t1,t2").
		Equal(tag.Name, "@apitags")

	tag, eof = l.Tag()
	a.Nil(tag).True(eof)
}
