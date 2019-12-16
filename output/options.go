// SPDX-License-Identifier: MIT

package output

import (
	"encoding/xml"
	"strings"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/openapi"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
)

// 几种输出的类型
const (
	ApidocXML   = "apidoc+xml"
	OpenapiYAML = "openapi+yaml"
	OpenapiJSON = "openapi+json"
)

var stylesheetURL string

type marshaler func(*doc.Doc) ([]byte, error)

func init() {
	stylesheetURL = vars.OfficialURL + "/" + vars.DocVersion() + "/apidoc.xsl"
}

// Options 指定了渲染输出的相关设置项。
type Options struct {
	// 导出的文件类型格式，默认为 apidoc 的 XML 文件。
	Type string `yaml:"type,omitempty"`

	// 文档的保存路径
	Path string `yaml:"path,omitempty"`

	// 只输出该标签的文档，若为空，则表示所有。
	Tags []string `yaml:"tags,omitempty"`

	// xslt 文件地址
	//
	// 默认值为 https://apidoc.tools/docs/ 下当前版本的 apidoc.xsl，比如：
	//  https://apidoc.tools/docs/v5/apidoc.xsl
	Style string `yaml:"style,omitempty"`

	procInst []string  // 保存所有 xml 的指令内容，包括编码信息
	marshal  marshaler // Type 对应的转换函数
	xml      bool      // 是否为 xml 内容
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

// buf 表示是否数据只存在于内存，如果为 true，则不需要检测 Path 是否正常。
func (o *Options) sanitize(buf bool) *message.SyntaxError {
	if o == nil {
		return message.NewLocaleError("", "", 0, locale.ErrRequired)
	}

	if o.Type == "" {
		o.Type = ApidocXML
	}

	switch o.Type {
	case ApidocXML:
		o.marshal = apidocMarshaler
	case OpenapiJSON:
		o.marshal = openapi.JSON
	case OpenapiYAML:
		o.marshal = openapi.YAML
	default:
		return message.NewLocaleError("", "type", 0, locale.ErrInvalidValue)
	}

	o.xml = strings.HasSuffix(o.Type, "+xml")

	if o.Path == "" && !buf {
		return message.NewLocaleError("", "path", 0, locale.ErrRequired)
	}

	if o.xml {
		if o.Style == "" {
			o.Style = stylesheetURL
		}

		o.procInst = []string{
			xml.Header,
			`<?xml-stylesheet type="text/xsl" href="` + o.Style + `"?>`,
		}
	}

	return nil
}

func apidocMarshaler(d *doc.Doc) ([]byte, error) {
	return xml.MarshalIndent(d, "", "\t")
}
