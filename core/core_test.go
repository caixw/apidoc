// SPDX-License-Identifier: MIT

package core

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/is"
)

// 对一些堂量的基本检测。
func TestConst(t *testing.T) {
	a := assert.New(t)

	a.True(len(Name) > 0)
	a.True(is.URL(RepoURL))
	a.True(is.URL(OfficialURL))
}

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

func TestRange_Contains(t *testing.T) {
	a := assert.New(t)

	r := Range{
		Start: Position{Line: 1, Character: 15},
		End:   Position{Line: 5, Character: 16},
	}

	a.True(r.Contains(Position{Line: 1, Character: 15}))
	a.True(r.Contains(Position{Line: 2, Character: 15}))
	a.True(r.Contains(Position{Line: 5, Character: 15}))
	a.False(r.Contains(Position{Line: 5, Character: 17}))
	a.False(r.Contains(Position{Line: 0, Character: 17}))
}
