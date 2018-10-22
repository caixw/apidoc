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

	// 两个性能文件的名称
	PprofCPU = "cpu"
	PprofMem = "mem"

	// 控制台的颜色
	InfoColor = colors.Green
	WarnColor = colors.Cyan
	ErroColor = colors.Red
)
