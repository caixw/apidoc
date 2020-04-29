// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert"
)

var _ Blocker = &pascalStringBlock{}

func TestPascalStringBlock(t *testing.T) {
	a := assert.New(t)

	b := newPascalStringBlock('"')
	a.NotNil(b)

	l, err := NewLexer([]byte(`"123""123"`), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	data, ok := b.EndFunc(l)
	a.True(ok).
		Equal(len(data), 0) // 不返回内容
	bs := l.Next(1)             // 继续向后推进，才会
	a.Empty(bs).True(l.AtEOF()) // 到达末尾

	l, err = NewLexer([]byte(`"123"""123"`), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(len(data), 0).            // 不返回内容
		Equal(string(l.All()), "123\"") // 未到达末尾
}
