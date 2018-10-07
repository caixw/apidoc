// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package output 对解析后的数据进行渲染输出。
package output

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/caixw/apidoc/docs"
)

// Render 渲染 docs 的内容，具体的渲染参数由 o 指定。
func Render(docs *docs.Docs, o *Options) error {
	for name, doc := range docs.Docs {
		if !o.contains(name) {
			continue
		}

		data, err := o.marshal(doc)
		if err != nil {
			return err
		}

		path := filepath.Join(o.Dir, name)
		if err = ioutil.WriteFile(path, data, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
