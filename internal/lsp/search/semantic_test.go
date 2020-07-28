// SPDX-License-Identifier: MIT

package search

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
)

func TestTokenBuilder_append(t *testing.T) {
	a := assert.New(t)

	b := &tokenBuilder{}
	b.append(core.Range{
		Start: core.Position{Line: 1, Character: 11},
		End:   core.Position{Line: 1, Character: 12},
	}, 1)
	a.Equal(b.tokens[0], []int{1, 11, 1, 1, 0})

	// 空值，不会添加内容
	b.append(core.Range{}, 1)
	a.Equal(1, len(b.tokens))

	// 长度为 0
	a.Panic(func() {
		b.append(core.Range{
			Start: core.Position{Line: 1, Character: 11},
			End:   core.Position{Line: 1, Character: 11},
		}, 1)
	})

	// 长度为负数
	a.Panic(func() {
		b.append(core.Range{
			Start: core.Position{Line: 1, Character: 11},
			End:   core.Position{Line: 1, Character: 10},
		}, 1)
	})
}

func TestTokenBuilder_build(t *testing.T) {
	a := assert.New(t)

	b := &tokenBuilder{
		tokens: [][]int{
			{1, 2, 5, 0, 0},
			{1, 12, 5, 0, 0},
			{2, 2, 7, 0, 0},
			{2, 2, 5, 0, 0},
			{2, 5, 5, 0, 0},
			{11, 1, 5, 0, 0},
		},
	}
	a.Equal(b.build(), []int{
		1, 2, 5, 0, 0,
		0, 10, 5, 0, 0,
		1, 2, 7, 0, 0,
		0, 0, 5, 0, 0,
		0, 3, 5, 0, 0,
		9, 1, 5, 0, 0,
	})
}

func TestTokenBuilder_sort(t *testing.T) {
	a := assert.New(t)

	b := &tokenBuilder{
		tokens: [][]int{
			{11, 1, 5, 0, 0},
			{1, 2, 5, 0, 0},
			{1, 12, 5, 0, 0},
			{2, 5, 5, 0, 0},
			{2, 2, 7, 0, 0},
			{2, 2, 5, 0, 0},
		},
	}

	b.sort()
	a.Equal(b.tokens, [][]int{
		{1, 2, 5, 0, 0},
		{1, 12, 5, 0, 0},
		{2, 2, 7, 0, 0},
		{2, 2, 5, 0, 0},
		{2, 5, 5, 0, 0},
		{11, 1, 5, 0, 0},
	})
}
