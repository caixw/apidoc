// SPDX-License-Identifier: MIT

package output

import (
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

// Options 指定了渲染输出的相关设置项。
type Options struct {
	// 文档的保存路径
	Path string `yaml:"path,omitempty"`

	// 只输出该标签的文档，若为空，则表示所有。
	Tags []string `yaml:"tags,omitempty"`
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

	return nil
}
