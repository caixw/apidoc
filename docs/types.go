// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package docs

import (
	"sync"

	"github.com/caixw/apidoc/docs/lexer"
)

// Docs 文档集合
type Docs struct {
	Version string          // 当前的程序版本
	Docs    map[string]*Doc // 文档集，键名为分组名称
	locker  sync.Mutex
}

// Markdown 表示可以使用 markdown 文档
type Markdown string

// Tag 标签内容
type Tag struct {
	Name        string   `yaml:"name" json:"name"`                                   // 字面名称，需要唯一
	Description Markdown `yaml:"description,omitempty" json:"description,omitempty"` // 具体描述
}

// Contact 描述联系方式
type Contact struct {
	Name  string `yaml:"name" json:"name"`
	URL   string `yaml:"url" json:"url"`
	Email string `yaml:"email,omitempty" json:"email,omitempty"`
}

// Link 表示一个链接
type Link struct {
	Text string `yaml:"text" json:"text"`
	URL  string `yaml:"url" json:"url"`
}

// Request 表示用户请求所表示的数据。
type Request = Body

// Response 表示一次请求或是返回的数据。
type Response struct {
	Body
	Status string `yaml:"status" json:"status"`
}

// Body 表示请求和返回的共有内容
type Body struct {
	Mimetype string     `yaml:"mimetype,omitempty" json:"mimetype,omitempty"`
	Headers  []*Header  `yaml:"headers,omitempty" json:"headers,omitempty"`
	Type     *Schema    `yaml:"type" json:"type"`
	Examples []*Example `yaml:"examples,omitempty" json:"examples,omitempty"`
}

// Header 报头
type Header struct {
	Name     string `yaml:"name" json:"name"`                             // 参数名称
	Summary  string `yaml:"summary" json:"summary"`                       // 参数介绍
	Optional bool   `yaml:"optional,omitempty" json:"optional,omitempty"` // 是否可以为空
}

// Param 简单参数的描述，比如查询参数等
type Param struct {
	Name     string  `yaml:"name" json:"name"`                             // 参数名称
	Type     *Schema `yaml:"type" json:"type"`                             // 类型
	Summary  string  `yaml:"summary" json:"summary"`                       // 参数介绍
	Optional bool    `yaml:"optional,omitempty" json:"optional,omitempty"` // 是否可以为空
}

// Example 示例
type Example struct {
	Mimetype string `yaml:"mimetype" json:"mimetype"`
	Summary  string `yaml:"summary,omitempty" json:"summary,omitempty"`
	Value    string `yaml:"value" json:"value"` // 示例内容
}

func (body *Body) parseExample(tag *lexer.Tag) error {
	data := tag.Split(3)
	if len(data) != 3 {
		return tag.ErrInvalidFormat()
	}

	if body.Examples == nil {
		body.Examples = make([]*Example, 0, 3)
	}

	body.Examples = append(body.Examples, &Example{
		Mimetype: string(data[0]),
		Summary:  string(data[1]),
		Value:    string(data[2]),
	})

	return nil
}

func (body *Body) parseHeader(tag *lexer.Tag) error {
	data := tag.Split(3)
	if len(data) != 3 {
		return tag.ErrInvalidFormat()
	}

	if body.Headers == nil {
		body.Headers = make([]*Header, 0, 3)
	}

	body.Headers = append(body.Headers, &Header{
		Name:     string(data[0]),
		Summary:  string(data[1]),
		Optional: isRequired(string(data[2])),
	})

	return nil
}
