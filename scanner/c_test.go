// Copyright 2015 by caixw, All rights reserved.
// Use of this source cCode is governed by a MIT
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
	cCode1 = `
int x = 5;
/* line1
line2
line3*/`

	cComment1 = []byte(` line1
line2
line3`)

	cCode2 = `
int x = 5;
// line1
// line2
// line3
`

	cComment2 = []byte(` line1
 line2
 line3
`)
)

func TestC__SingleBlock(t *testing.T) {
	a := assert.New(t)

	fn := func(cCode string, cComment []byte) {
		block, pos := C([]byte(cCode))
		a.Equal(block, cComment).Equal(pos, len(cCode))
	}

	fn(cCode1, cComment1)
	fn(cCode2, cComment2)
}

//////////////////////////////////// 测试多个注释块

var (
	cMultBlock1 = `
int x = 5
/*
 comment1
 comment1
 */

 /*comment2
 comment2
 */
`

	cComments1 = [][]byte{
		[]byte(`
 comment1
 comment1
 `),
		[]byte(`comment2
 comment2
 `),
	}

	cMultBlock2 = `
int x=5
// comment1
// comment1
// 

// comment2
// comment2
`

	cComments2 = [][]byte{
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

	fn := func(cCode string, cComments [][]byte) {
		cCodebs := []byte(cCode)
		for _, c := range cComments {
			block, pos := C(cCodebs)
			a.Equal(block, c)
			cCodebs = cCodebs[pos:]
		}
	}

	fn(cMultBlock1, cComments1)
	fn(cMultBlock2, cComments2)
}
