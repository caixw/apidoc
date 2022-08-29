// SPDX-License-Identifier: MIT

package core

var (
	mainVersion = "7.2.3"
	metadata    string
	fullVersion = mainVersion
)

func init() {
	if metadata != "" {
		fullVersion += "+" + metadata
	}
}

// FullVersion 完整的版本号
//
// 会包含版本号、构建日期和最后的提交 ID，大致格式如下：
//
//	version+buildDate.commitHash
func FullVersion() string {
	return fullVersion
}

// Version 程序的版本号
//
// 遵守 https://semver.org/lang/zh-CN/ 规则。
// 程序不兼容或是文档格式不兼容时，需要提升主版本号。
func Version() string {
	return mainVersion
}
