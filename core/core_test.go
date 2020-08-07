// SPDX-License-Identifier: MIT

package core

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/is"
)

var (
	_ Ranger = Range{}
	_ Ranger = &Range{}
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

func TestRange_Equal(t *testing.T) {
	a := assert.New(t)

	r1 := Range{}
	a.True(r1.Equal(Range{}))
	a.False(r1.Equal(Range{End: Position{Line: 1}}))

	r1 = Range{End: Position{Line: 1}}
	a.True(r1.Equal(Range{End: Position{Line: 1}}))
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

func TestLocation_Equal(t *testing.T) {
	a := assert.New(t)

	l := Location{}
	a.True(l.Equal(Location{}))
	a.False(l.Equal(Location{URI: URI(".")}))

	l = Location{URI: URI("."), Range: Range{Start: Position{Line: 1}}}
	a.True(l.Equal(Location{URI: URI("."), Range: Range{Start: Position{Line: 1}}}))
	a.False(l.Equal(Location{}))
}
