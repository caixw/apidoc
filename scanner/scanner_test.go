// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

import (
	"testing"

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
