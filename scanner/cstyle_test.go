// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

import (
	"testing"

	"github.com/issue9/assert"
)

var _ scanFunc = cstyle

var code1 = `
int x = 5;
/* line1
line2
line3*/
`
var comment1 = []byte(` line1
line2
line3`)

var code2 = `
int x = 5;
// line1
// line2
// line3
`

var comment2 = []byte(` line1
 line2
 line3
`)

var code3 = `
  int x=5
  // line1
  // line2
  // line3
`
var comment3 = []byte(` line1
 line2
 line3
`)

func TestCStyle(t *testing.T) {
	a := assert.New(t)

	fn := func(code string, comment []byte) {
		s := &scanner{
			data: []byte(code),
		}
		block, ln, err := cstyle(s)
		a.NotError(err).NotNil(block)
		a.Equal(block, comment).Equal(ln, 2)
	}

	fn(code1, comment1)
	fn(code2, comment2)
	fn(code3, comment3)
}
