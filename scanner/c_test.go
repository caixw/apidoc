// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

import (
	"testing"

	"github.com/caixw/apidoc/core"
	"github.com/issue9/assert"
)

var _ core.ScanFunc = C

//////////////////////////////////// 测试单个注释块

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

func TestC__SingleBlock(t *testing.T) {
	a := assert.New(t)

	fn := func(code string, comment []byte) {
		block, pos := C([]byte(code))
		a.Equal(block, comment).Equal(pos, len(code))
	}

	fn(code1, comment1)
	fn(code2, comment2)
}

//////////////////////////////////// 测试多个注释块

var (
	mb1 = `
int x = 5
/*
 comment1
 comment1
 */

 /*comment2
 comment2
 */
`

	comments1 = [][]byte{
		[]byte(`
 comment1
 comment1
 `),
		[]byte(`comment2
 comment2
 `),
	}

	mb2 = `
int x=5
// comment1
// comment1
// 

// comment2
// comment2
`

	comments2 = [][]byte{
		[]byte(` comment1
 comment1
 
`),
		[]byte(` comment2
 comment2
`),
	}
)

func TestC__MultBlock(t *testing.T) {
	a := assert.New(t)

	fn := func(code string, comments [][]byte) {
		codebs := []byte(code)
		for _, c := range comments {
			block, pos := C(codebs)
			a.Equal(block, c)
			codebs = codebs[pos:]
		}
		//block, pos := C([]byte(code))
		//a.Equal(block, comment).Equal(pos, len(code))
	}

	fn(mb1, comments1)
	fn(mb2, comments2)
}
