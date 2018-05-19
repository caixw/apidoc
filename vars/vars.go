// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package vars 提供了一些公共的函数、结构体及代码级别的设置项。
package vars

import "github.com/issue9/term/colors"

// 一些公用的常量
const (
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

	// 默认的分组名称，在不指定分组名称的时候，系统会给其加到此分组中。
	DefaultGroupName = "openapi"

	// 默认的语言，目前仅能保证简体中文是最新的。
	// 需要保证存在于 locale.locales 中，否则运行时会报错。
	DefaultLocale = "zh-Hans"

	// 两个性能文件的名称
	PprofCPU = "cpu"
	PprofMem = "mem"

	// 控制台的颜色
	InfoColor = colors.Green
	WarnColor = colors.Cyan
	ErroColor = colors.Red
)
