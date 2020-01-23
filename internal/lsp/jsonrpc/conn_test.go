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

	// 参数类型不正确
	a.Panic(func() {
		newHandler(func(bool, *f, *int) error { return nil })
	})

	// 参数类型不正确
	a.Panic(func() {
		newHandler(func(bool, *int, *f) error { return nil })
	})

	// 参数类型不正确
	a.Panic(func() {
		newHandler(func(*bool, *int, *int) error { return nil })
	})

	// 返回值不正确
	a.Panic(func() {
		newHandler(func(bool, *int, *int) int { return 0 })
	})

	// 返回值不正确
	a.NotPanic(func() {
		newHandler(func(bool, *int, *int) *Error { return nil })
	})

	// 正常签名
	a.NotPanic(func() {
		newHandler(func(bool, *int, *int) error { return nil })
	})
}
