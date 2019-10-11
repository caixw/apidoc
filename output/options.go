// SPDX-License-Identifier: MIT

package output

import (
	"encoding/xml"

	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
)

var stylesheetURL string

func init() {
	stylesheetURL = vars.OfficialURL + "/" + vars.DocVersion() + "/apidoc.xsl"
}

// Options 指定了渲染输出的相关设置项。
type Options struct {
	// 文档的保存路径
	Path string `yaml:"path,omitempty"`

	// 只输出该标签的文档，若为空，则表示所有。
	Tags []string `yaml:"tags,omitempty"`

	// xslt 文件地址
	//
	// 默认值为 https://apidoc.tools/docs/ 下当前版本的 apidoc.xsl
	Style string `yaml:"style,omitempty"`

	procInst []string
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

func (o *Options) sanitize() *message.SyntaxError {
	if o == nil {
		return message.NewLocaleError("", "", 0, locale.ErrRequired)
	}

	if o.Path == "" {
		return message.NewLocaleError("", "path", 0, locale.ErrRequired)
	}

	if o.Style == "" {
		o.Style = stylesheetURL
	}

	o.procInst = []string{
		xml.Header,
		`<?xml-stylesheet type="text/xsl" href="` + o.Style + `"?>`,
	}

	return nil
}
