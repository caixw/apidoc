// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package output 对解析后的数据进行渲染输出。
package output

import (
	"io/ioutil"
	"os"

	"github.com/caixw/apidoc/docs"
)

// Render 渲染 docs 的内容，具体的渲染参数由 o 指定。
func Render(ds *docs.Docs, o *Options) error {
	if len(o.Groups) == 0 {
		return render(ds, o)
	}

	ds1 := make([]*docs.Doc, 0, len(o.Groups))
	for _, doc := range ds.Docs {
		if o.contains(doc.Group) {
			ds1 = append(ds1, doc)
		}
	}
	ds.Docs = ds1

	return render(ds, o)
}

func render(docs *docs.Docs, o *Options) error {
	data, err := o.marshal(docs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(o.Path, data, os.ModePerm)
}
