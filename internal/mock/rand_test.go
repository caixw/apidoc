// SPDX-License-Identifier: MIT

package mock

const indent = "    "

// 当前文件提供了一些生成随机测试数据的函数

// 测试数据为了方便验证正确性，生成的值是固定的，
// 而普通的 mock 数据值是随机的。通过此值判断生成哪种数据。
//
// 测试环境下，生成的数据，数值固定为 1024，字符串固定为 “1024”
// 枚举值，则永远取第一个元素作为值。
var testOptions = &GenOptions{
	Number:    func() interface{} { return 1024 },
	String:    func() string { return "1024" },
	Bool:      func() bool { return true },
	SliceSize: func() int { return 5 },
	Index:     func(max int) int { return 0 },
}
