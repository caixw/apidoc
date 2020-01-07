// SPDX-License-Identifier: MIT

// Package output 对解析后的数据进行渲染输出。
package output

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/caixw/apidoc/v6/doc"
)

// Render 渲染 doc 的内容
func Render(d *doc.Doc, opt *Options) error {
	if err := opt.sanitize(false); err != nil {
		return err
	}

	buf, err := buffer(d, opt)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(opt.Path, buf.Bytes(), os.ModePerm)
}

// Buffer 将内容导出到内存
func Buffer(d *doc.Doc, opt *Options) (*bytes.Buffer, error) {
	if err := opt.sanitize(true); err != nil {
		return nil, err
	}

	return buffer(d, opt)
}

func buffer(d *doc.Doc, opt *Options) (*bytes.Buffer, error) {
	filterDoc(d, opt)

	buf := new(bytes.Buffer)

	if opt.xml {
		for _, v := range opt.procInst {
			if _, err := buf.WriteString(v); err != nil {
				return nil, err
			}

			if err := buf.WriteByte('\n'); err != nil {
				return nil, err
			}
		} // end range opt.procInst
	}

	data, err := opt.marshal(d)
	if err != nil {
		return nil, err
	}
	if _, err = buf.Write(data); err != nil {
		return nil, err
	}

	return buf, nil
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
