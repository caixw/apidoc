// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

import (
	"testing"

	"github.com/caixw/apidoc/core"
	"github.com/issue9/assert"
)

func TestScanner_next(t *testing.T) {
	a := assert.New(t)

	s := &scanner{
		data: []byte("ab\ncd"),
	}

	a.Equal('a', s.next())
	a.Equal('b', s.next())
	a.Equal('\n', s.next())
	a.Equal('c', s.next())
	a.Equal('d', s.next())
	a.Equal(eof, s.next())
	a.Equal(eof, s.next())
}

func TestScanner_match(t *testing.T) {
	a := assert.New(t)

	s := &scanner{
		data: []byte("ab\ncd"),
	}

	a.False(s.match("b")).Equal(0, s.pos)
	a.True(s.match("a")).Equal(1, s.pos)

	s.backup()
	s.backup()
	a.True(s.match("a")).Equal(1, s.pos)
	a.True(s.match("b")).Equal(2, s.pos)

	s.pos = len(s.data)
	a.False(s.match("ab"))
}

func TestScanner_skipSpace(t *testing.T) {
	a := assert.New(t)

	s := &scanner{
		data: []byte("  ab\n  cd"),
	}

	s.skipSpace()
	a.Equal(s.next(), 'a')

	s.next()
	s.skipSpace()
	a.Equal(s.next(), 'c')
}

func TestScanner_lineNumber(t *testing.T) {
	a := assert.New(t)

	s := &scanner{
		data: []byte("adf\n\nadf"),
		pos:  3,
	}
	a.Equal(1, s.lineNumber())

	s.pos = 4
	a.Equal(2, s.lineNumber())
}

func TestScanner_scanFile(t *testing.T) {
	a := assert.New(t)

	scanFile(cstyle, "./testcode/php1.php")

	a.Equal(len(docs), 2)
	a.Equal(docs[0].Group, "php1").
		Equal(docs[0].Method, "get").
		Equal(docs[0].URL, "/api/php1/get")
}

func TestScan(t *testing.T) {
	a := assert.New(t)

	docsMu.Lock()
	docs = []*core.Doc{}
	docsMu.Unlock()

	docs, err := Scan(&Options{SrcDir: "./testcode", Recursive: true, Type: "", Exts: nil})
	a.NotError(err).NotNil(docs)
	a.Equal(4, len(docs))

	for _, v := range docs {
		switch {
		case v.URL == "/api/php1/get":
			a.Equal(v.Method, "get")
		case v.URL == "/api/php2/post":
			a.Equal(v.Method, "post")
			a.Equal(v.Group, "php2")
		}
	}
}
