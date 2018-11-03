// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package output 对解析后的数据进行渲染输出。
package output

import (
	"io/ioutil"
	"os"

	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/errors"
	"github.com/caixw/apidoc/internal/locale"
	opt "github.com/caixw/apidoc/options"
)

// Render 渲染 doc 的内容，具体的渲染参数由 o 指定。
func Render(d *doc.Doc, o *opt.Output) error {
	if o == nil {
		return &errors.Error{
			Field:      "output",
			MessageKey: locale.ErrRequired,
		}
	}

	opt, err := buildOptions(o)
	if err != nil {
		return err
	}

	filterDoc(d, opt)

	data, err := opt.marshal(d)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(opt.Path, data, os.ModePerm)
}

func filterDoc(d *doc.Doc, o *options) {
	if len(o.Tags) == 0 {
		return
	}

	tags := make([]*doc.Tag, 0, len(o.Tags))
	for _, tag := range d.Tags {
		if o.contains(tag.Name) {
			tags = append(tags, tag)
		}
	}
	d.Tags = tags

	apis := make([]*doc.API, 0, len(d.Apis))
	for _, api := range d.Apis {
		if o.contains(api.Tags...) {
			apis = append(apis, api)
		}
	}
	d.Apis = apis
}
