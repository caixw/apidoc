// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/errors"
	"github.com/caixw/apidoc/internal/input"
)

func newLexerString(data string) *lexer {
	erro := log.New(ioutil.Discard, "[ERRO]", 0)
	warn := log.New(ioutil.Discard, "[WARN]", 0)
	h := errors.NewHandler(errors.NewLogHandlerFunc(erro, warn))
	return newLexer(input.Block{Data: []byte(data)}, h)
}

func newTagString(data string) *lexerTag {
	return newLexerString(data).tag()
}

func TestLexer_Tag(t *testing.T) {
	a := assert.New(t)
	l := newLexerString(`@api get /path desc
markdown desc line1
markdown desc line2
   @apigroup xxx
 @apitags t1,t2`)

	// @api
	tag := l.tag()
	a.NotNil(tag)
	a.Equal(tag.Line, 0).
		Equal(string(tag.Data), `get /path desc
markdown desc line1
markdown desc line2`). // 最后的换行符会被去掉
		Equal(tag.Name, "@api")

	// @apigroup
	tag = l.tag()
	a.NotNil(tag)
	a.Equal(tag.Line, 3).
		Equal(string(tag.Data), "xxx").
		Equal(tag.Name, "@apigroup")

	// @apitags
	tag = l.tag()
	a.NotNil(tag)
	a.Equal(tag.Line, 4).
		Equal(string(tag.Data), "t1,t2").
		Equal(tag.Name, "@apitags")

	// 没有标签了
	tag = l.tag()
	a.Nil(tag)

	// 没有标签了，多次调用，结果是一样的
	tag = l.tag()
	a.Nil(tag)

	l = newLexerString("@api")
	tag = l.tag()
	a.Equal(tag.Name, "@api").
		Empty(tag.Data)
}

func TestLexer_Backup(t *testing.T) {
	a := assert.New(t)
	l := &lexer{}

	tag := &lexerTag{Name: "@api"}
	l.backup(tag)
	a.Equal(l.tag(), tag)
	a.Nil(l.tag())
}

func TestSplitWords(t *testing.T) {
	a := assert.New(t)

	tag := []byte("@tag s1\ts2  \n  s3")

	bs := splitWords(tag, 1)
	a.Equal(bs, [][]byte{[]byte("@tag s1\ts2  \n  s3")})

	bs = splitWords(tag, 2)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1\ts2  \n  s3")})

	bs = splitWords(tag, 3)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2  \n  s3")})

	bs = splitWords(tag, 4)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2"), []byte("s3")})

	// 不够
	bs = splitWords(tag, 5)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2"), []byte("s3")})

	tag = []byte("@tag s1 s2  ")
	bs = splitWords(tag, 4)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2")})
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
