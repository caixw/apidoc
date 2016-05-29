// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"errors"
	"os"
)

type Options struct {
	AppVersion string `json:"-"`       // apidoc程序的版本号
	Elapsed    int64  `json:"-"`       // 编译用时，单位毫秒
	Version    string `json:"version"` // 文档的版本号
	Dir        string `json:"dir"`     // 文档的保存目录
	Title      string `json:"title"`   // 文档的标题
	BaseURL    string `json:"baseURL"` // api文档中url的前缀
	// Language string // 产生的ui界面语言
	//Type string   `json:"type"` // 输出的语言格式
	//Groups     []string `json:"groups"`     // 需要打印的分组内容。
	//Timezone   string   `json:"timezone"`   // 时区
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

func checkOptions(o *Options) error {
	if len(o.Dir) == 0 {
		return errors.New("未指定 Dir")
	}
	o.Dir += string(os.PathSeparator)

	if len(o.Title) == 0 {
		o.Title = "APIDOC"
	}

	return nil
}
