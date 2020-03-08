// SPDX-License-Identifier: MIT

package build

import (
	"bytes"
	"encoding/xml"
	"strings"

	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/openapi"
	"github.com/caixw/apidoc/v6/internal/vars"
	"github.com/caixw/apidoc/v6/message"
	"github.com/caixw/apidoc/v6/spec"
)

// 几种输出的类型
const (
	ApidocXML   = "apidoc+xml"
	OpenapiYAML = "openapi+yaml"
	OpenapiJSON = "openapi+json"
)

const stylesheetURL = vars.OfficialURL + "/docs/" + spec.MajorVersion + "/apidoc.xsl"

type marshaler func(*spec.APIDoc) ([]byte, error)

// Output 指定了渲染输出的相关设置项。
type Output struct {
	// 导出的文件类型格式，默认为 apidoc 的 XML 文件。
	Type string `yaml:"type,omitempty"`

	// 文档的保存路径
	Path string `yaml:"path"`

	// 只输出该标签的文档，若为空，则表示所有。
	Tags []string `yaml:"tags,omitempty"`

	// xslt 文件地址
	//
	// 默认值为 https://apidoc.tools/docs/ 下当前版本的 apidoc.xsl，比如：
	//  https://apidoc.tools/docs/v6/apidoc.xsl
	Style string `yaml:"style,omitempty"`

	procInst []string  // 保存所有 xml 的指令内容，包括编码信息
	marshal  marshaler // Type 对应的转换函数
	xml      bool      // 是否为 xml 内容
}

func (o *Output) contains(tags ...string) bool {
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

// Sanitize 验证 options 的数据是否都合规
func (o *Output) Sanitize() error {
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

func apidocMarshaler(d *spec.APIDoc) ([]byte, error) {
	return xml.MarshalIndent(d, "", "\t")
}

func (o *Output) buffer(d *spec.APIDoc) (*bytes.Buffer, error) {
	filterDoc(d, o)

	buf := new(bytes.Buffer)

	if o.xml {
		for _, v := range o.procInst {
			if _, err := buf.WriteString(v); err != nil {
				return nil, err
			}

			if err := buf.WriteByte('\n'); err != nil {
				return nil, err
			}
		} // end range opt.procInst
	}

	data, err := o.marshal(d)
	if err != nil {
		return nil, err
	}
	if _, err = buf.Write(data); err != nil {
		return nil, err
	}

	return buf, nil
}

func filterDoc(d *spec.APIDoc, o *Output) {
	if len(o.Tags) == 0 {
		return
	}

	tags := make([]*spec.Tag, 0, len(o.Tags))
	for _, tag := range d.Tags {
		if o.contains(tag.Name) {
			tags = append(tags, tag)
		}
	}
	d.Tags = tags

	apis := make([]*spec.API, 0, len(d.Apis))
	for _, api := range d.Apis {
		if o.contains(api.Tags...) {
			apis = append(apis, api)
		}
	}
	d.Apis = apis
}
