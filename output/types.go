// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

type Options struct {
	AppVersion string // apidoc程序的版本号
	Version    string // 文档的版本号
	DocDir     string // 文档的保存目录
	Title      string // 文档的标题
	Elapsed    int64  // 编译用时，单位毫秒
	// Language string // 产生的ui界面语言
}

// 用于页首和页脚的附加信息
type info struct {
	Groups     map[string]string // 分组名称与文件的对照表
	CurrGroup  string            // 当前所在的分组页，若为空，表示在列表页
	Date       string
	Version    string
	AppVersion string
	Title      string
	Elapsed    string
}
