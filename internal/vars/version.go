// SPDX-License-Identifier: MIT

package vars

// 程序的版本号
//
// 遵守 https://semver.org/lang/zh-CN/ 规则。
// 程序不兼容或是文档格式不兼容时，需要提升主版本号。
const version = "6.0.0"

var (
	fullVersion = version
	buildDate   string
	commitHash  string
)

func init() {
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
