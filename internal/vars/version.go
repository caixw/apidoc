// SPDX-License-Identifier: MIT

package vars

import "strings"

// 主版本号，实际版本号可能还会加上构建日期，
// 可通过 Version() 函数获取完整的版本号。
const version = "5.0.0"

var (
	fullVersion = version
	docVersion  string
	buildDate   string
	commitHash  string
)

func init() {
	dotIndex := strings.IndexByte(version, '.')
	if dotIndex <= 0 {
		panic("无效的版本号 mainVersion")
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
// 取程序版本号的主版本号部分，返回格式为 v + 主版本号，比如 v5
func DocVersion() string {
	return docVersion
}
