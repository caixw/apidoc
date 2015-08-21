// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

import (
	"testing"

	"github.com/caixw/apidoc/core"
	"github.com/issue9/assert"
)

var _ core.ScanFunc = CStyle

var (
	code1 = `
int x = 5;
/* line1
line2
line3*/`

	comment1 = []byte(` line1
line2
line3`)

	code2 = `
int x = 5;
// line1
// line2
// line3
`

	comment2 = []byte(` line1
 line2
 line3
`)
)

func TestCStyle(t *testing.T) {
	a := assert.New(t)

	fn := func(code string, comment []byte) {
		block, pos := CStyle([]byte(code))
		a.Equal(block, comment).Equal(pos, len(code))
	}

	fn(code1, comment1)
	fn(code2, comment2)
}
