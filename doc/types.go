// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"strings"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/doc/lexer"
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
	Name     string      `yaml:"name" json:"name"`                             // 参数名称
	Type     string      `yaml:"type" json:"type"`                             // 类型
	Summary  string      `yaml:"summary" json:"summary"`                       // 参数介绍
	Optional bool        `yaml:"optional,omitempty" json:"optional,omitempty"` // 是否可以为空
	Default  interface{} `yaml:"default,omitempty" json:"default,omitempty"`   // 默认值
}

// Example 示例
type Example struct {
	Mimetype string `yaml:"mimetype" json:"mimetype"`
	Summary  string `yaml:"summary,omitempty" json:"summary,omitempty"`
	Value    string `yaml:"value" json:"value"` // 示例内容
}

// 解析示例代码，格式如下：
//  @apiExample application/json
//  {
//      "id": 1,
//      "name": "name",
//  }
func (body *Body) parseExample(tag *lexer.Tag) error {
	lines := tag.Lines(2)
	if len(lines) != 2 {
		return tag.ErrInvalidFormat()
	}

	words := lexer.SplitWords(lines[0], 2)

	if body.Examples == nil {
		body.Examples = make([]*Example, 0, 3)
	}

	example := &Example{
		Mimetype: string(words[0]),
		Value:    string(lines[1]),
	}
	if len(words) == 2 { // 如果存在简介
		example.Summary = string(words[1])
	}

	body.Examples = append(body.Examples, example)

	return nil
}

func (body *Body) parseHeader(tag *lexer.Tag) error {
	data := tag.Words(3)
	if len(data) != 3 {
		return tag.ErrInvalidFormat()
	}

	if body.Headers == nil {
		body.Headers = make([]*Header, 0, 3)
	}

	body.Headers = append(body.Headers, &Header{
		Name:     string(data[0]),
		Summary:  string(data[2]),
		Optional: isOptional(data[1]),
	})

	return nil
}

func newResponse(l *lexer.Lexer, tag *lexer.Tag) (*Response, error) {
	data := tag.Words(3)
	if len(data) != 3 {
		return nil, tag.ErrInvalidFormat()
	}

	schema := &Schema{}
	if err := buildSchema(tag, schema, nil, data[1], nil, data[2]); err != nil {
		return nil, err
	}
	resp := &Response{
		Body: Body{
			Mimetype: string(data[1]),
		},
	}

	for tag, eof := l.Tag(); !eof; tag, eof = l.Tag() {
		switch strings.ToLower(tag.Name) {
		case "@apiexample":
			if err := resp.parseExample(tag); err != nil {
				return nil, err
			}
		case "@apiheader":
			if err := resp.parseHeader(tag); err != nil {
				return nil, err
			}
		case "@apiparam":
			data := tag.Words(4)
			if len(data) != 4 {
				return nil, tag.ErrInvalidFormat()
			}

			if err := buildSchema(tag, schema, data[0], data[1], data[2], data[3]); err != nil {
				return nil, err
			}
		default:
			l.Backup(tag)
			return resp, nil
		}
	}

	return resp, nil
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

	opt, def, err := parseOptional(string(data[1]), data[2])
	if err != nil {
		return nil, err
	}

	return &Param{
		Name:     string(data[0]),
		Summary:  string(data[3]),
		Type:     string(data[1]),
		Default:  def,
		Optional: opt,
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
