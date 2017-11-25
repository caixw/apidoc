// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package app 提供了一些公共的函数、结构体及代码级别的设置项。
package app

// 一些代码级别的配置项。
// 可运行 go test 来检测常量是否符合规范。
const (
	// 主版本号，实际版本号可能还会加上构建日期，
	// 可通过 Version() 函数获取实际的版本号。
	mainVersion = "4.0.0"

	// 程序的正式名称
	Name = "apidoc"

	// 源码仓库地址
	RepoURL = "https://github.com/caixw/apidoc"

	// 官网
	OfficialURL = "http://apidoc.tools"

	// 配置文件名称。
	ConfigFilename = ".apidoc.yaml"

	// 默认的文档标题
	DefaultTitle = "APIDOC"

	// 默认的分组名称，在不指定分组名称的时候，
	// 系统会给其加到此分组中，同时也是默认的索引文件名。
	DefaultGroupName = "index"

	// 输出的 profile 文件的名称
	Profile = "apidoc.prof"

	// 默认的语言，目前仅能保证简体中文是最新的。
	DefaultLocale = "cmn-Hans"

	// 生成的 JSON 数据存放的目录
	JSONDataDirName = "data"
)

var (
	version    string
	buildDate  string
	commitHash string
)

func init() {
	if len(buildDate) == 0 {
		version = mainVersion
	} else {
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
