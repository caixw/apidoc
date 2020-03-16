// SPDX-License-Identifier: MIT

package openapi

import "github.com/caixw/apidoc/v6/core"

// 数据验证接口
type sanitizer interface {
	Sanitize() *core.SyntaxError
}
