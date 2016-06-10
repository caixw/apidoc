// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// output 对解后的数据进行渲染输出。
package output

import (
	"fmt"
	"os"
	"time"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/doc"
)

// 支持的渲染方式
var renderTypes = []string{
	"html",
}

// 渲染输出的相关设置项。
type Options struct {
	Dir     string        `json:"dir"`  // 文档的保存目录
	Type    string        `json:"type"` // 渲染方式，默认为 html
	Elapsed time.Duration `json:"-"`    // 编译用时

	// Language string // 产生的ui界面语言
	//Groups     []string `json:"groups"`     // 需要打印的分组内容。
	//Timezone   string   `json:"timezone"`   // 时区
}

// 对 Options 作一些初始化操作。
func (o *Options) Init() *app.OptionsError {
	if len(o.Dir) == 0 {
		return &app.OptionsError{Field: "Dir", Message: "不能为空"}
	}
	o.Dir += string(os.PathSeparator)

	if !isSuppertedType(o.Type) {
		return &app.OptionsError{Field: "Type", Message: "不支持该类型"}
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
