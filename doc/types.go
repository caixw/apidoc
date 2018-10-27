// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"github.com/issue9/is"

	"github.com/caixw/apidoc/doc/lexer"
	"github.com/caixw/apidoc/doc/schema"
)

// Markdown 表示可以使用 markdown 文档
type Markdown string

// Tag 标签内容
type Tag struct {
	Name        string   `yaml:"name" json:"name"`                                   // 字面名称，需要唯一
	Description Markdown `yaml:"description,omitempty" json:"description,omitempty"` // 具体描述
}

// Server 服务信息
type Server struct {
	Name        string   `yaml:"name" json:"name"` // 字面名称，需要唯一
	URL         string   `yaml:"url" json:"url"`
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

// Param 简单参数的描述，比如查询参数等
type Param struct {
	Name     string         `yaml:"name" json:"name"`                             // 参数名称
	Type     *schema.Schema `yaml:"type" json:"type"`                             // 类型
	Summary  string         `yaml:"summary" json:"summary"`                       // 参数介绍
	Optional bool           `yaml:"optional,omitempty" json:"optional,omitempty"` // 是否为可选参数
}

// 解析链接元素，格式如下：
//  @tag text https://example.com
func newLink(tag *lexer.Tag) (*Link, error) {
	data := tag.Words(2)
	if len(data) != 2 {
		return nil, tag.ErrInvalidFormat()
	}

	if !is.URL(data[1]) {
		return nil, tag.ErrInvalidFormat()
	}

	return &Link{
		Text: string(data[0]),
		URL:  string(data[1]),
	}, nil
}

// 解析参数标签，格式如下：
// 用于路径参数和查义参数，request 和 request 中的不在此解析
//  @tag name type.subtype optional.defaultValue markdown desc
func newParam(tag *lexer.Tag) (*Param, error) {
	data := tag.Words(4)
	if len(data) != 4 {
		return nil, tag.ErrInvalidFormat()
	}

	s := &schema.Schema{}
	if err := s.Build(tag, nil, data[1], data[2], nil); err != nil {
		return nil, err
	}

	return &Param{
		Name:     string(data[0]),
		Summary:  string(data[3]),
		Type:     s,
		Optional: s.Default != nil,
	}, nil
}

// 解析联系人标签内容，格式可以是：
//  @apicontact name xx@example.com https://example.com
//  @apicontact name https://example.com xx@example.com
//  @apicontact name xx@example.com
//  @apicontact name https://example.com
func newContact(tag *lexer.Tag) (*Contact, error) {
	data := tag.Words(3)

	if len(data) < 2 {
		return nil, tag.ErrInvalidFormat()
	}

	contact := &Contact{Name: string(data[0])}
	v := string(data[1])
	switch checkContactType(v) {
	case 1:
		contact.URL = v
	case 2:
		contact.Email = v
	default:
		return nil, tag.ErrInvalidFormat()
	}

	if len(data) == 3 {
		v := string(data[2])
		switch checkContactType(v) {
		case 1:
			contact.URL = v
		case 2:
			contact.Email = v
		default:
			return nil, tag.ErrInvalidFormat()
		}
	}

	return contact, nil
}

func checkContactType(v string) int8 {
	switch {
	case is.Email(v): // Email 也属于一种 URL
		return 2
	case is.URL(v):
		return 1
	default:
		return 0
	}
}
