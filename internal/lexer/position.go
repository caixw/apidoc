// SPDX-License-Identifier: MIT

package lexer

import "github.com/caixw/apidoc/v6/core"

// Position 描述 Lexer 中的定位信息
type Position struct {
	core.Position

	// 表示的是字节的偏移量，
	// 而 Position.Character 表示的是当前行`字符`的偏移量
	Offset int
}

// Equal 判断与 v 是否相等
func (p Position) Equal(v Position) bool {
	return p.Offset == v.Offset
}
