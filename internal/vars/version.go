// SPDX-License-Identifier: MIT

package vars

// 主版本号，实际版本号可能还会加上构建日期，
// 可通过 Version() 函数获取实际的版本号。
const mainVersion = "5.0.0"

var (
	version    = mainVersion
	buildDate  string
	commitHash string
)

func init() {
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
