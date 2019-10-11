// SPDX-License-Identifier: MIT

package vars

import "strings"

// 主版本号，实际版本号可能还会加上构建日期，
// 可通过 Version() 函数获取实际的版本号。
const mainVersion = "5.0.0"

var (
	version    = mainVersion
	docVersion string
	buildDate  string
	commitHash string
)

func init() {
	dotIndex := strings.IndexByte(mainVersion, '.')
	if dotIndex <= 0 {
		panic("无效的版本号 mainVersion")
	}
	docVersion = "v" + mainVersion[:dotIndex]

	if buildDate != "" {
		version = mainVersion + "+" + buildDate
	}
}

// Version 完整的版本号
func Version() string {
	return version
}

// CommitHash Git 上最后的提交记录 hash 值。
func CommitHash() string {
	return commitHash
}

// DocVersion 文档的版本号
//
// 取程序版本号的主版本号部分
func DocVersion() string {
	return docVersion
}
