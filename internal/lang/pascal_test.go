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

	l := &lexer{data: []byte(`"123""123"`)}
	a.True(b.BeginFunc(l))
	ret, ok := b.EndFunc(l)
	a.True(ok).
		Equal(len(ret), 0). // 不返回内容
		True(l.atEOF())     // 到达末尾

	l = &lexer{data: []byte(`"123"""123"`)}
	a.True(b.BeginFunc(l))
	ret, ok = b.EndFunc(l)
	a.True(ok).
		Equal(len(ret), 0).                    // 不返回内容
		Equal(string(l.data[l.pos:]), "123\"") // 未到达末尾
}
