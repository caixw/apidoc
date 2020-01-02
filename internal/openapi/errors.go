// SPDX-License-Identifier: MIT

package openapi

import "github.com/caixw/apidoc/v5/message"

// 数据验证接口
type sanitizer interface {
	Sanitize() *message.SyntaxError
}
