// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/internal/locale"
	"github.com/caixw/apidoc/internal/options"
	"github.com/caixw/apidoc/output/openapi"
)

type marshaler func(v *doc.Doc) ([]byte, error)

// 文档类型定义
const (
	ApidocJSON  = "apidoc+json"
	ApidocYAML  = "apidoc+yaml"
	OpenapiJSON = "openapi+json"
	OpenapiYAML = "openapi+yaml"
	RamlJSON    = "raml+json"
)

// Options 指定了渲染输出的相关设置项。
type Options struct {
	// 文档的保存路径，包含目录和文件名，若为空，则为当前目录下的
	Path string `yaml:"path,omitempty"`

	// 输出类型
	Type string `yaml:"type,omitempty"`

	// 只输出该标签的文档，若为空，则表示所有。
	Tags []string `yaml:"tags,omitempty"`

	marshal marshaler // 根据 type 决定转换的函数
}

func (o *Options) contains(tags ...string) bool {
	if len(o.Tags) == 0 {
		return true
	}

	for _, t := range o.Tags {
		for _, tag := range tags {
			if tag == t {
				return true
			}
		}
	}
	return false
}

// Sanitize 对 Options 作一些初始化操作。
func (o *Options) Sanitize() error {
	// TODO 改用默认值
	if o.Path == "" {
		return options.NewFieldError("path", locale.Sprintf(locale.ErrRequired))
	}

	if o.Type == "" {
		o.Type = ApidocJSON
	}

	switch o.Type {
	case ApidocJSON:
		o.marshal = apidocJSONMarshal
	case ApidocYAML:
		o.marshal = apidocYAMLMarshal
	case OpenapiJSON:
		o.marshal = openapi.JSON
	case OpenapiYAML:
		o.marshal = openapi.YAML
	case RamlJSON:
		// TODO
	default:
		return options.NewFieldError("type", locale.Sprintf(locale.ErrInvalidValue))
	}

	return nil
}

func apidocJSONMarshal(v *doc.Doc) ([]byte, error) {
	return json.Marshal(v)
}

func apidocYAMLMarshal(v *doc.Doc) ([]byte, error) {
	return yaml.Marshal(v)
}
