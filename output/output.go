// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"errors"
	"fmt"
	"os"

	"github.com/caixw/apidoc/doc"
)

// 支持的渲染方式
var renderTypes = []string{
	"html",
}

type Options struct {
	AppVersion string `json:"-"`       // apidoc 程序的版本号
	Elapsed    int64  `json:"-"`       // 编译用时，单位毫秒
	Version    string `json:"version"` // 文档的版本号
	Dir        string `json:"dir"`     // 文档的保存目录
	Title      string `json:"title"`   // 文档的标题
	BaseURL    string `json:"baseURL"` // api 文档中 url 的前缀
	Type       string `json:"type"`    // 渲染方式，默认为 html

	// Language string // 产生的ui界面语言
	//Groups     []string `json:"groups"`     // 需要打印的分组内容。
	//Timezone   string   `json:"timezone"`   // 时区
}

// 对 Options 作一些初始化操作。
func (o *Options) Init() error {
	if len(o.Dir) == 0 {
		return errors.New("未指定 Dir")
	}
	o.Dir += string(os.PathSeparator)

	if len(o.Title) == 0 {
		o.Title = "APIDOC"
	}

	if !isSuppertedType(o.Type) {
		return fmt.Errorf("不支持的渲染类型：[%v]", o.Type)
	}

	return nil
}

// 渲染 docs 的内容，具体的渲染参数由 o 指定。
func Render(docs *doc.Doc, o *Options) error {
	switch o.Type {
	case "html":
		return html(docs, o)
	default:
		return fmt.Errorf("不支持的渲染方式:[%v]", o.Type)
	}
}

func isSuppertedType(typ string) bool {
	for _, k := range renderTypes {
		if k == typ {
			return true
		}
	}

	return false
}
