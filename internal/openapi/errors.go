// SPDX-License-Identifier: MIT

package openapi

	import "github.com/caixw/apidoc/v5/message"

// Sanitizer 数据验证接口
type Sanitizer interface {
	Sanitize() *message.SyntaxError
}
