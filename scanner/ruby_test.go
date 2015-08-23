// Copyright 2015 by caixw, All rights reserved.
// Use of this source rubyCode is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

import (
	"testing"

	"github.com/caixw/apidoc/core"
	"github.com/issue9/assert"
)

var _ core.ScanFunc = Ruby

//////////////////////////////////// 测试单个注释块

var (
	rubyCode1 = `
=begin
line1
line2
line3
=end`

	rubyComment1 = []byte(`
line1
line2
line3
`)

	rubyCode2 = `
int x = 5;
# line1
# line2
# line3
`

	rubyComment2 = []byte(` line1
 line2
 line3
`)
)

func TestRuby__SingleBlock(t *testing.T) {
	a := assert.New(t)

	fn := func(rubyCode string, rubyComment []byte) {
		block, pos := Ruby([]byte(rubyCode))
		a.Equal(block, rubyComment).Equal(pos, len(rubyCode))
	}

	fn(rubyCode1, rubyComment1)
	fn(rubyCode2, rubyComment2)
}

//////////////////////////////////// 测试多个注释块

var (
	rubyMultBlock1 = `
int x = 5
=begin
 comment1
 comment1
=end

=begin
comment2
comment2
=end
`

	rubyComments1 = [][]byte{
		[]byte(`
 comment1
 comment1
`),
		[]byte(`
comment2
comment2
`),
	}

	rubyMultBlock2 = `
int x=5
# comment1
# comment1
# 

# comment2
# comment2
`

	rubyComments2 = [][]byte{
		[]byte(` comment1
 comment1
 
`),
		[]byte(` comment2
 comment2
`),
	}
)

func TestRuby__MultBlock(t *testing.T) {
	a := assert.New(t)

	fn := func(rubyCode string, rubyComments [][]byte) {
		rubyCodebs := []byte(rubyCode)
		for _, c := range rubyComments {
			block, pos := Ruby(rubyCodebs)
			a.Equal(block, c)
			rubyCodebs = rubyCodebs[pos:]
		}
	}

	fn(rubyMultBlock1, rubyComments1)
	fn(rubyMultBlock2, rubyComments2)
}
