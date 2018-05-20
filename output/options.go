// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"encoding/json"
	"os"
	"time"

	"github.com/issue9/utils"
	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/config/conferr"
	"github.com/caixw/apidoc/locale"
)

const (
	typeJSON = "json"
	typeYAML = "yaml"
)

// Options 指定了渲染输出的相关设置项。
type Options struct {
	// 文档的保存目录
	Dir string `yaml:"dir"`

	// 是否对 Dir 作清理操作
	//
	// 如果为 true，则每次都会清空 Dir 目录下的所有内容；
	// 否则为覆盖同名文件的操作，默认为 false。
	Clean bool `yaml:"clean,omitempty"`

	// 仅输出这些组，为空表示输出所有
	//
	// 若指定的组名实际上不存在，则不会有任何影响。
	Groups []string `yaml:"groups,omitempty"`

	// 输出类型
	//
	// 与 openapi 支持的格式相同，目前可以是 json 和 yaml
	// 默认值为 yaml
	Type string `yaml:"type,omitempty"`

	marshal func(v interface{}) ([]byte, error) // 根据 type 决定转换的函数
	Elapsed time.Duration                       `yaml:"-"` // TODO 添加到 openapi 的扩展字段中
}

// Sanitize 对 Options 作一些初始化操作。
func (o *Options) Sanitize() *conferr.Error {
	if len(o.Dir) == 0 {
		return conferr.New("dir", locale.Sprintf(locale.ErrRequired))
	}

	if o.Clean {
		if err := os.RemoveAll(o.Dir); err != nil {
			return conferr.New("dir", err.Error())
		}
	}

	if !utils.FileExists(o.Dir) {
		if err := os.MkdirAll(o.Dir, os.ModePerm); err != nil {
			return conferr.New("dir", err.Error())
		}
	}

	if o.Type == "" {
		o.Type = typeYAML
	}

	switch o.Type {
	case typeJSON:
		o.marshal = json.Marshal
	case typeYAML:
		o.marshal = yaml.Marshal
	default:
		return conferr.New("type", locale.Sprintf(locale.ErrInvalidValue))
	}

	return nil
}

// 指定的组名是否包含在输出列表里。
func (o *Options) contains(group string) bool {
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
