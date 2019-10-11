// SPDX-License-Identifier: MIT

// Package output 对解析后的数据进行渲染输出。
package output

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/caixw/apidoc/v5/doc"
)

// Render 渲染 doc 的内容
func Render(d *doc.Doc, opt *Options) error {
	if err := opt.sanitize(); err != nil {
		return err
	}

	filterDoc(d, opt)

	data, err := xml.MarshalIndent(d, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(opt.Path, data, os.ModePerm)
}

func filterDoc(d *doc.Doc, o *Options) {
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
