// SPDX-License-Identifier: MIT

package jsonrpc

import (
	"testing"

	"github.com/issue9/assert"
)

func TestNewHandler(t *testing.T) {
	a := assert.New(t)

	type f func()

	a.Panic(func() {
		newHandler(5)
	})

	// 参数数量不正确
	a.Panic(func() {
		newHandler(func(*int) error { return nil })
	})

	// 参数数量不正确
	a.Panic(func() {
		newHandler(func(*f, *int) error { return nil })
	})

	// 参数数量不正确
	a.Panic(func() {
		newHandler(func(*int, *f) error { return nil })
	})

	// 返回值不正确
	a.Panic(func() {
		newHandler(func(*int, *int) int { return 0 })
	})

	// 返回值不正确
	a.NotPanic(func() {
		newHandler(func(*int, *int) *Error { return nil })
	})

	// 正常签名
	a.NotPanic(func() {
		newHandler(func(*int, *int) error { return nil })
	})
}
