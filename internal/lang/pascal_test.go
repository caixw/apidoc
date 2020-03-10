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

	l := &Lexer{data: []byte(`"123""123"`)}
	a.True(b.BeginFunc(l))
	raw, data, ok := b.EndFunc(l)
	a.True(ok).
		Equal(len(data), 0). // 不返回内容
		Equal(len(raw), 0).  // 不返回内容
		True(l.AtEOF())      // 到达末尾

	l = &Lexer{data: []byte(`"123"""123"`)}
	a.True(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(len(data), 0).                      // 不返回内容
		Equal(len(raw), 0).                       // 不返回内容
		Equal(string(l.data[l.offset:]), "123\"") // 未到达末尾
}
