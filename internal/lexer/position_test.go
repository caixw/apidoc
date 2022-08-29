// SPDX-License-Identifier: MIT

package lexer

import (
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/apidoc/v7/core"
)

func TestPosition_Equal(t *testing.T) {
	a := assert.New(t, false)

	v1 := Position{}
	a.True(v1.Equal(Position{}))

	v1.Offset = 5
	a.False(v1.Equal(Position{}))
	a.True(v1.Equal(Position{Offset: 5, Position: core.Position{Line: 5}}))
}

func TestPosition_AddRune_SubRune(t *testing.T) {
	a := assert.New(t, false)

	v1 := Position{}

	v1 = v1.AddRune(' ')
	a.Equal(v1.Line, 0).Equal(v1.Offset, 1).Equal(v1.Character, 1)

	v1 = v1.AddRune('中')
	a.Equal(v1.Line, 0).Equal(v1.Offset, 4).Equal(v1.Character, 2)

	v1 = v1.AddRune('\n')
	a.Equal(v1.Line, 1).Equal(v1.Offset, 5).Equal(v1.Character, 0)
	v1 = v1.AddRune('\n')
	a.Equal(v1.Line, 2).Equal(v1.Offset, 6).Equal(v1.Character, 0)

	v1 = v1.AddRune('a')
	a.Equal(v1.Line, 2).Equal(v1.Offset, 7).Equal(v1.Character, 1)

	v1 = v1.SubRune('b')
	a.Equal(v1.Line, 2).Equal(v1.Offset, 6).Equal(v1.Character, 0)

	v1 = v1.SubRune('\n')
	a.Equal(v1.Line, 1).Equal(v1.Offset, 5).Equal(v1.Character, 0)
}
