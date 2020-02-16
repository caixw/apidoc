// SPDX-License-Identifier: MIT

package doc

import (
	xmessage "golang.org/x/text/message"

	"github.com/caixw/apidoc/v6/message"
)

// Block 表示原始的注释代码块
type Block struct {
	File string
	Line int
	Data []byte
}

// NewLocaleError 生成基于 Block 定位的错误信息
func (b *Block) NewLocaleError(field string, key xmessage.Reference, v ...interface{}) error {
	return message.NewLocaleError(b.File, field, b.Line, key, v...)
}
