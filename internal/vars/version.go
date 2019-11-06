// SPDX-License-Identifier: MIT

package vars

import "strings"

// 程序的版本号
//
// 遵守 https://semver.org/lang/zh-CN/ 规则。
// 程序不兼容或是文档格式不兼容时，需要提升主版本号。
const version = "5.1.0"

var (
	fullVersion = version
	docVersion  string
	buildDate   string
	commitHash  string
)

func init() {
	dotIndex := strings.IndexByte(version, '.')
	if dotIndex <= 0 {
		panic("无效的版本号 version")
	}
	docVersion = "v" + version[:dotIndex]

	if buildDate != "" {
		fullVersion = version + "+" + buildDate
	}
}

// Version 完整的版本号
func Version() string {
	return fullVersion
}

// CommitHash Git 上最后的提交记录 hash 值。
func CommitHash() string {
	return commitHash
}

// DocVersion 文档的版本号
//
// 当文档格式不兼容时，此值也会发生变化。
func DocVersion() string {
	return docVersion
}
