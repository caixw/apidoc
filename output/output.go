// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package output 对解析后的数据进行渲染输出。
package output

import (
	"io/ioutil"
	"os"

	"github.com/caixw/apidoc/doc"
)

// Render 渲染 doc 的内容，具体的渲染参数由 o 指定。
func Render(doc *doc.Doc, o *Options) error {
	data, err := o.marshal(doc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(o.Path, data, os.ModePerm)
}
