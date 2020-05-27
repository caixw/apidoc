// SPDX-License-Identifier: MIT

package mock

import (
	"math/rand"
	"strconv"

	"github.com/issue9/rands"

	"github.com/caixw/apidoc/v7/internal/ast"
)

var randOptions = &struct {
	maxSliceSize  int
	maxNumber     int
	maxStringSize int
	minStringSize int
	StringData    []byte
	indent        string

	// 当前文件提供了一些生成随机测试数据的函数

	// 测试数据为了方便验证正确性，生成的值是固定的，
	// 而普通的 mock 数据值是随机的。通过此值判断生成哪种数据。
	//
	// 测试环境下，生成的数据，数值固定为 1024，字符串固定为 “1024”
	// 枚举值，则永远取第一个元素作为值。
	test bool
}{
	maxSliceSize:  100,
	maxNumber:     10000,
	maxStringSize: 100,
	minStringSize: 5,
	StringData:    rands.AlphaNumber,
	indent:        "    ",
}

func generateBool() bool {
	if randOptions.test {
		return true
	}
	return (rand.Int() % 2) == 0
}

func isEnum(p *ast.Param) bool {
	return len(p.Enums) > 0
}

func generateNumber(p *ast.Param) int64 {
	if isEnum(p) {
		index := 0
		if !randOptions.test {
			index = rand.Intn(len(p.Enums))
		}
		v, err := strconv.ParseInt(p.Enums[index].Value.V(), 10, 32)
		if err != nil { // 这属于文档定义错误，直接 panic
			panic(err)
		}
		return v
	}

	if randOptions.test {
		return 1024
	}
	return rand.Int63n(int64(randOptions.maxNumber))
}

func generateString(p *ast.Param) string {
	if isEnum(p) {
		index := 0
		if !randOptions.test {
			index = rand.Intn(len(p.Enums))
		}
		return p.Enums[index].Value.V()
	}

	if randOptions.test {
		return "1024"
	}
	return rands.String(randOptions.minStringSize, randOptions.maxStringSize, randOptions.StringData)
}

// 生成随机的数组长度
func generateSliceSize() int {
	if randOptions.test {
		return 5
	}
	return rand.Intn(randOptions.maxSliceSize)
}
