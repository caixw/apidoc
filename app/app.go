// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// app 提供了一些公共的函数、结构体及设置项。
package app

// 一些代码级别的配置项。
const (
	// 版本号
	//
	// 版本号按照 http://semver.org/lang/zh-CN/ 中的规则，分成以下四个部分：
	// 主版本号.次版本号.修订号.修订日期
	Version = "2.2.59.160604"

	// 程序的正式名称
	Name = "apidoc"

	// 源码仓库地址
	RepoURL = "https://github.com/caixw/apidoc"

	// 官网
	OfficialURL = "https://caixw.github.io/apidoc"

	// 配置文件名称。
	ConfigFilename = ".apidoc.json"
)
