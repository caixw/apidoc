// SPDX-License-Identifier: MIT

package core

import (
	"testing"

	"github.com/issue9/assert"
)

func TestPosition_Equal(t *testing.T) {
	a := assert.New(t)

	p1 := Position{}
	a.True(p1.Equal(Position{}))
	a.False(p1.Equal(Position{Line: -1}))
}

func TestRange_IsEmpty(t *testing.T) {
	a := assert.New(t)

	r := Range{
		Start: Position{},
		End:   Position{},
	}
	a.True(r.IsEmpty())

	r.Start.Line = 11
	r.End.Line = 11
	a.True(r.IsEmpty())

	r.Start.Character = 55
	a.False(r.IsEmpty())
	r.End.Character = 55
	a.True(r.IsEmpty())
}
