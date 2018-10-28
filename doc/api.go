// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"bytes"
	"strings"

	"github.com/caixw/apidoc/doc/lexer"
	"github.com/caixw/apidoc/doc/schema"
)

// API 表示单个 API 文档
type API struct {
	Method      string      `yaml:"method" json:"method"`
	Path        string      `yaml:"path" json:"path"`
	Summary     string      `yaml:"summary" json:"summary"`
	Description Markdown    `yaml:"description,omitempty" json:"description,omitempty"`
	Tags        []string    `yaml:"tags,omitempty" json:"tags,omitempty"`
	Queries     []*Param    `yaml:"queries,omitempty" json:"queries,omitempty"` // 查询参数
	Params      []*Param    `yaml:"params,omitempty" json:"params,omitempty"`   // URL 参数
	Requests    []*Request  `yaml:"requests,omitempty" json:"requests,omitempty"`
	Responses   []*Response `yaml:"responses" json:"responses"`
	Deprecated  string      `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	Server      string      `yaml:"server" json:"server"`
}

// Param 简单参数的描述，比如查询参数等
type Param struct {
	Name     string         `yaml:"name" json:"name"`                             // 参数名称
	Type     *schema.Schema `yaml:"type" json:"type"`                             // 类型
	Summary  string         `yaml:"summary" json:"summary"`                       // 参数介绍
	Optional bool           `yaml:"optional,omitempty" json:"optional,omitempty"` // 是否为可选参数
}

func (doc *Doc) parseAPI(l *lexer.Lexer) error {
	api := &API{}
	var parse func(*lexer.Lexer, *lexer.Tag) error

	for tag := l.Tag(); tag != nil; tag = l.Tag() {
		switch strings.ToLower(tag.Name) {
		case "@api":
			parse = api.parseAPI
		case "@apirequest":
			parse = api.parseRequest
		case "@apiresponse":
			parse = api.parseResponse
		default:
			return tag.ErrInvalidTag()
		}

		if err := parse(l, tag); err != nil {
			return err
		}
	}

	doc.append(api)

	return nil
}

func (doc *Doc) append(api *API) {
	doc.locker.Lock()
	doc.Apis = append(doc.Apis, api)
	doc.locker.Unlock()
}

type apiParser func(*API, *lexer.Lexer, *lexer.Tag) error

var apiParsers = map[string]apiParser{
	"@api":           (*API).parseapi,
	"@apiserver":     (*API).parseServer,
	"@apitags":       (*API).parseTags,
	"@apideprecated": (*API).parseDeprecated,
	"@apiquery":      (*API).parseQuery,
	"@apiparam":      (*API).parseParam,
}

// 分析 @api 以及子标签
func (api *API) parseAPI(l *lexer.Lexer, tag *lexer.Tag) error {
	l.Backup(tag) // 进来时，第一个肯定是 @api 标签，退回该标签，统一让 for 处理。

	for tag := l.Tag(); tag != nil; tag = l.Tag() {
		fn, found := apiParsers[strings.ToLower(tag.Name)]
		if !found {
			l.Backup(tag)
			return nil
		}

		if err := fn(api, l, tag); err != nil {
			return err
		}
	}

	return nil
}

func (api *API) parseapi(l *lexer.Lexer, tag *lexer.Tag) error {
	if api.Method != "" || api.Path != "" || api.Summary != "" {
		return tag.ErrDuplicateTag()
	}
	data := tag.Words(3)
	if len(data) != 3 {
		return tag.ErrInvalidFormat()
	}

	api.Method = strings.ToUpper(string(data[0])) // TODO 验证请求方法
	api.Path = string(data[1])
	api.Summary = string(data[2])

	return nil
}

func (api *API) parseServer(l *lexer.Lexer, tag *lexer.Tag) error {
	if api.Server != "" {
		return tag.ErrDuplicateTag()
	}
	api.Server = string(tag.Data)
	return nil
}

var separatorTag = []byte{','}

func (api *API) parseTags(l *lexer.Lexer, tag *lexer.Tag) error {
	if len(api.Tags) > 0 {
		return tag.ErrDuplicateTag()
	}

	tags := bytes.Split(tag.Data, separatorTag)
	api.Tags = make([]string, 0, len(tags))
	for _, tag := range tags {
		api.Tags = append(api.Tags, string(tag))
	}

	return nil
}

func (api *API) parseDeprecated(l *lexer.Lexer, tag *lexer.Tag) error {
	api.Deprecated = string(tag.Data)
	return nil
}

func (api *API) parseQuery(l *lexer.Lexer, tag *lexer.Tag) error {
	if api.Params == nil {
		api.Params = make([]*Param, 0, 10)
	}

	p, err := newParam(tag)
	if err != nil {
		return err
	}
	api.Queries = append(api.Queries, p)
	return nil
}

func (api *API) parseParam(l *lexer.Lexer, tag *lexer.Tag) error {
	if api.Params == nil {
		api.Params = make([]*Param, 0, 10)
	}

	p, err := newParam(tag)
	if err != nil {
		return err
	}
	api.Params = append(api.Params, p)

	return nil
}

func (api *API) parseRequest(l *lexer.Lexer, tag *lexer.Tag) error {
	data := tag.Words(3)
	if len(data) != 3 {
		return tag.ErrInvalidFormat()
	}

	if api.Requests == nil {
		api.Requests = make([]*Request, 0, 3)
	}

	req := &Request{
		Mimetype: string(data[1]),
		Type:     &schema.Schema{},
	}
	api.Requests = append(api.Requests, req)

	if err := req.Type.Build(tag, nil, data[1], nil, data[2]); err != nil {
		return err
	}

	for tag := l.Tag(); tag != nil; tag = l.Tag() {
		switch strings.ToLower(tag.Name) {
		case "@apiexample":
			if err := req.parseExample(tag); err != nil {
				return err
			}
		case "@apiheader":
			if err := req.parseHeader(tag); err != nil {
				return err
			}
		case "@apiparam":
			params := tag.Words(4)
			if len(params) != 4 {
				return tag.ErrInvalidFormat()
			}

			if err := req.Type.Build(tag, params[0], params[1], params[2], params[3]); err != nil {
				return err
			}
		default:
			l.Backup(tag)
			return nil
		}
	}

	return nil
}

func (api *API) parseResponse(l *lexer.Lexer, tag *lexer.Tag) error {
	if api.Responses == nil {
		api.Responses = make([]*Response, 0, 10)
	}

	resp, err := newResponse(l, tag)
	if err != nil {
		return err
	}
	api.Responses = append(api.Responses, resp)

	return nil
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
