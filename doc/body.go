// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/caixw/apidoc/doc/lexer"
	"github.com/caixw/apidoc/doc/schema"
)

// Request 表示用户请求所表示的数据。
type Request = Body

// Response 表示一次请求或是返回的数据。
type Response struct {
	Body
	Status int `yaml:"status" json:"status"`
}

// API 和 Doc 都有这个属性，且都需要 parseResponse 方法。
// 抽象为一个嵌套对象使用。
type responses struct {
	Responses []*Response `yaml:"responses,omitempty" json:"responses,omitempty"`
}

// Body 表示请求和返回的共有内容
type Body struct {
	Mimetype string         `yaml:"mimetype,omitempty" json:"mimetype,omitempty"`
	Headers  []*Header      `yaml:"headers,omitempty" json:"headers,omitempty"`
	Type     *schema.Schema `yaml:"type" json:"type"`
	Examples []*Example     `yaml:"examples,omitempty" json:"examples,omitempty"`
}

// Header 报头
type Header struct {
	Name     string `yaml:"name" json:"name"`                             // 参数名称
	Summary  string `yaml:"summary" json:"summary"`                       // 参数介绍
	Optional bool   `yaml:"optional,omitempty" json:"optional,omitempty"` // 是否可以为空
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

var requiredBytes = []byte("required")

func isOptional(data []byte) bool {
	return !bytes.Equal(bytes.ToLower(data), requiredBytes)
}

// 解析 @apiHeader 标签，格式如下：
//  @apiheader content-type required desc
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

// 解析 @apiparam 标签，格式如下：
//  @apiparam group object reqiured desc
func (body *Body) parseParam(tag *lexer.Tag) error {
	data := tag.Words(4)
	if len(data) != 4 {
		return tag.ErrInvalidFormat()
	}

	return body.Type.Build(tag, data[0], data[1], data[2], data[3])
}

func (resps *responses) parseResponse(l *lexer.Lexer, tag *lexer.Tag) error {
	if resps.Responses == nil {
		resps.Responses = make([]*Response, 0, 3)
	}

	resp, err := newResponse(l, tag)
	if err != nil {
		return err
	}
	resps.Responses = append(resps.Responses, resp)

	return nil
}

// 解析 @apiResponse 及子标签，格式如下：
//  @apiresponse 200 array.object * 通用的返回内容定义
//  @apiheader content-type required desc
//  @apiparam id int reqiured desc
//  @apiparam name string reqiured desc
//  @apiparam group object reqiured desc
//  @apiparam group.id int reqiured desc
func newResponse(l *lexer.Lexer, tag *lexer.Tag) (*Response, error) {
	data := tag.Words(4)
	if len(data) != 4 {
		return nil, tag.ErrInvalidFormat()
	}

	status, err := strconv.Atoi(string(data[0]))
	if err != nil {
		return nil, tag.ErrInvalidFormat()
	}

	s := &schema.Schema{}
	if err := s.Build(tag, nil, data[1], nil, data[3]); err != nil {
		return nil, err
	}
	resp := &Response{
		Status: status,
		Body: Body{
			Mimetype: string(data[2]),
			Type:     s,
		},
	}

LOOP:
	for tag := l.Tag(); tag != nil; tag = l.Tag() {
		fn := resp.parseExample
		switch strings.ToLower(tag.Name) {
		case "@apiexample":
			fn = resp.parseExample
		case "@apiheader":
			fn = resp.parseHeader
		case "@apiparam":
			fn = resp.parseParam
		default:
			l.Backup(tag)
			break LOOP
		}

		if err := fn(tag); err != nil {
			return nil, err
		}
	}

	return resp, nil
}
