// SPDX-License-Identifier: MIT

// Package makeutil 为 internal/docs 中一系列生成工具提供的公共内容
package makeutil

// Header 表示输出文件的文件头内容
const Header = "该文件由工具自动生成，请勿手动修改！"

// PanicError 如果 err 不为 nil，则 panic
func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}
