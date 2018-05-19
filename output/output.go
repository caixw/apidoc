// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package output 对解析后的数据进行渲染输出。
package output

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/issue9/utils"
	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/openapi"
)

// Options 指定了渲染输出的相关设置项。
type Options struct {
	Dir     string        `yaml:"dir"`              // 文档的保存目录
	Groups  []string      `yaml:"groups,omitempty"` // 仅输出这些组，为空表示输出所有
	Elapsed time.Duration `yaml:"-"`                // TODO 添加到 openapi 的扩展字段中
}

// Sanitize 对 Options 作一些初始化操作。
func (o *Options) Sanitize() *openapi.Error {
	if len(o.Dir) == 0 {
		return &openapi.Error{Field: "dir", Message: locale.Sprintf(locale.ErrRequired)}
	}

	return nil
}

// 指定的组是否需要输出
func (o *Options) groupIsEnable(group string) bool {
	if len(o.Groups) == 0 {
		return true
	}

	for _, g := range o.Groups {
		if g == group {
			return true
		}
	}

	return false
}

// Render 渲染 docs 的内容，具体的渲染参数由 o 指定。
func Render(docs map[string]*openapi.OpenAPI, o *Options) error {
	// 文档目录下的文件名可能改变，先清除目录下的所有文件。
	if err := os.RemoveAll(o.Dir); err != nil {
		return err
	}

	if !utils.FileExists(o.Dir) {
		if err := os.MkdirAll(o.Dir, os.ModePerm); err != nil {
			return err
		}
	}

	// 未指定输出组，则输出所有内容。
	if len(o.Groups) == 0 {
		for name, doc := range docs {
			render(filepath.Join(o.Dir, name), doc)
		}

		return nil
	}

	for _, group := range o.Groups {
		doc, found := docs[group]
		if !found {
			return &openapi.Error{Field: "group", Message: locale.Sprintf(locale.ErrGroupNotExists)}
		}

		if err := render(filepath.Join(o.Dir, group), doc); err != nil {
			return err
		}
	}
	return nil
}

func render(path string, doc *openapi.OpenAPI) error {
	data, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, os.ModePerm)
}
