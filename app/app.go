// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// app 提供了一些公共的函数、结构体及代码级别的设置项。
package app

import (
	"os"

	"github.com/issue9/term/colors"
)

// 一些代码级别的配置项。
// 可运行 go test 来检测常量是否符合规范。
const (
	// 版本号，规则参照 http://semver.org
	Version = "3.0.1-alpha+20160723"

	// 程序的正式名称
	Name = "apidoc"

	// 源码仓库地址
	RepoURL = "https://github.com/caixw/apidoc"

	// 官网
	OfficialURL = "http://apidoc.tools"

	// 配置文件名称。
	ConfigFilename = ".apidoc.json"

	// 默认的文档标题
	DefaultTitle = "APIDOC"

	// 默认的分组名称，在不指定分组名称的时候，
	// 系统会给其加到此分组中，同时也是默认的索引文件名。
	DefaultGroupName = "index"

	// 输出的 profile 文件的名称
	Profile = "apidoc.prof"

	// 默认的语言，目前仅能保证简体中文是最新的。
	DefaultLocale = "cmn-Hans"
)

// Message 向终端输出不同颜色的提示信息
//
// color 是输出的字体颜色，仅对 prefix
// 参数起作用，其它字符串依然使用系统默认的颜色。
func Message(out *os.File, color colors.Color, prefix string, v ...interface{}) {
	colors.Fprint(out, color, colors.Default, prefix)
	colors.Fprintln(out, colors.Default, colors.Default, v...)
}

// Warn 输出警告性的信息
func Warn(v ...interface{}) {
	Message(os.Stderr, colors.Cyan, "[WARN] ", v...)
}

// Error 输出错误的信息
func Error(v ...interface{}) {
	Message(os.Stderr, colors.Red, "[ERROR] ", v...)
}

// Info 输出提示信息
func Info(v ...interface{}) {
	Message(os.Stdout, colors.Green, "[INFO] ", v...)
}
