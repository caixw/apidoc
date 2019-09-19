// SPDX-License-Identifier: MIT

package output

import (
	"encoding/xml"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/openapi"
	"github.com/caixw/apidoc/v5/message"
)

// 文档类型定义
const (
	ApidocXML   = "apidoc+xml"
	OpenapiJSON = "openapi+json"
	OpenapiYAML = "openapi+yaml"
	RAMLYAML    = "raml+yaml"
)

type marshaler func(v *doc.Doc) ([]byte, error)

// Options 指定了渲染输出的相关设置项。
type Options struct {
	// 文档的保存路径，建议使用绝对路径。
	Path string `yaml:"path,omitempty"`

	// 输出类型
	Type string `yaml:"type,omitempty"`

	// 只输出该标签的文档，若为空，则表示所有。
	Tags []string `yaml:"tags,omitempty"`

	marshal marshaler
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

// Sanitize 检测内容是否合法
func (o *Options) Sanitize() *message.SyntaxError {
	if o == nil {
		return message.NewError("", "", 0, locale.ErrRequired)
	}

	if o.Path == "" {
		return message.NewError("", "path", 0, locale.ErrRequired)
	}

	if o.Type == "" {
		o.Type = ApidocXML
	}

	switch o.Type {
	case ApidocXML:
		o.marshal = xmlMarshal
	case OpenapiJSON:
		o.marshal = openapi.JSON
	case OpenapiYAML:
		o.marshal = openapi.YAML
	case RAMLYAML:
		// TODO
	default:
		return message.NewError("", "type", 0, locale.ErrInvalidValue)
	}

	return nil
}

func xmlMarshal(v *doc.Doc) ([]byte, error) {
	return xml.Marshal(v)
}
