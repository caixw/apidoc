// SPDX-License-Identifier: MIT

package mock

import (
	"math/rand"

	"github.com/issue9/rands"
)

// 当前文件提供了一些生成随机测试数据的函数

// 测试数据为了方便验证正确性，生成的值是固定的，
// 而普通的 mock 数据值是随机的。通过此值判断生成哪种数据。
var test = false

func generateBool() bool {
	if test {
		return true
	}
	return (rand.Int() % 2) == 0
}

func generateNumber() int64 {
	if test {
		return 1024
	}
	return rand.Int63()
}

func generateString() string {
	if test {
		return "1024"
	}
	return rands.String(10, 32, rands.AlphaNumber)
}

func generateSliceSize() int {
	if test {
		return 5
	}
	return rand.Intn(100)
}
