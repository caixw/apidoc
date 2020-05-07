// SPDX-License-Identifier: MIT

package makeutil

// Header 表示输出文件的文件头内容
const Header = "该文件由工具自动生成，请勿手动修改！"

// PanicError 如果 err 不为 nil，则 panic
func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}
