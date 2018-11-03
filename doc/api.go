// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"bytes"
	"sort"
	"strings"

	"github.com/caixw/apidoc/doc/lexer"
	"github.com/caixw/apidoc/doc/schema"
	"github.com/caixw/apidoc/errors"
	"github.com/caixw/apidoc/internal/locale"
)

// API 表示单个 API 文档
type API struct {
	responses
	Method      string     `yaml:"method" json:"method"`
	Path        string     `yaml:"path" json:"path"`
	Summary     string     `yaml:"summary" json:"summary"`
	Description Markdown   `yaml:"description,omitempty" json:"description,omitempty"`
	Tags        []string   `yaml:"tags,omitempty" json:"tags,omitempty"`
	Queries     []*Param   `yaml:"queries,omitempty" json:"queries,omitempty"` // 查询参数
	Params      []*Param   `yaml:"params,omitempty" json:"params,omitempty"`   // URL 参数
	Requests    []*Request `yaml:"requests,omitempty" json:"requests,omitempty"`
	Deprecated  string     `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	Servers     []string   `yaml:"servers" json:"servers"`

	// 记录起始位置，方便错误定位
	file string
	line int
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

	for tag := l.Tag(); tag != nil; tag = l.Tag() {
		parse := api.parseAPI
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

	doc.locker.Lock()
	doc.Apis = append(doc.Apis, api)
	doc.locker.Unlock()

	return nil
}

type apiParser func(*API, *lexer.Lexer, *lexer.Tag) error

var apiParsers = map[string]apiParser{
	"@api":           (*API).parseapi,
	"@apiservers":    (*API).parseServers,
	"@apitags":       (*API).parseTags,
	"@apideprecated": (*API).parseDeprecated,
	"@apiquery":      (*API).parseQuery,
	"@apiparam":      (*API).parseParam,
}

// 检测内容是否正确
func (api *API) check() error {
	names, err := api.getPathParams()
	if err != nil {
		return err
	}

	if len(names) != len(api.Params) {
		return api.errInvalidFormat("@api") // TODO 专门的错误提示，地址参数与实际参数对不上
	}

	for _, p := range api.Params {
		if !names[p.Name] {
			return api.errInvalidFormat("@api") // TODO 专门的错误提示，地址参数与实际参数对不上
		}
	}

	return nil
}

func (api *API) getPathParams() (map[string]bool, error) {
	names := make(map[string]bool, len(api.Params))

	state := '}'
	index := 0
	for i, b := range api.Path {
		switch b {
		case '{':
			if state != '}' {
				return nil, api.errInvalidFormat("@api")
			}

			state = '{'
			index = i
		case '}':
			if state != '{' {
				return nil, api.errInvalidFormat("@api")
			}
			names[api.Path[index+1:i]] = true
			state = '}'
		}
	} // end for

	// 缺少 } 结束符号
	if state == '{' {
		return nil, api.errInvalidFormat("@api")
	}

	return names, nil
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

// 解析 @api 标签，格式如下：
//  @api GET /path summary
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
	api.file = tag.File
	api.line = tag.Line

	return nil
}

func (api *API) errInvalidFormat(tag string) error {
	return &errors.Error{
		File:       api.file,
		Line:       api.line,
		Field:      tag,
		MessageKey: locale.ErrInvalidFormat,
	}
}

// 解析 @apiServers 标签，格式如下：
//  @apiServers s1,s2
func (api *API) parseServers(l *lexer.Lexer, tag *lexer.Tag) (err error) {
	if len(api.Servers) > 0 {
		return tag.ErrDuplicateTag()
	}

	if len(tag.Data) == 0 {
		return tag.ErrInvalidFormat()
	}

	if api.Servers, err = splitToArray(tag); err != nil {
		return err
	}

	return nil
}

// 解析 @apiTags 标签，格式如下：
//  @apiTags t1,t2
func (api *API) parseTags(l *lexer.Lexer, tag *lexer.Tag) (err error) {
	if len(api.Tags) > 0 {
		return tag.ErrDuplicateTag()
	}

	if len(tag.Data) == 0 {
		return tag.ErrInvalidFormat()
	}

	if api.Tags, err = splitToArray(tag); err != nil {
		return err
	}

	return nil
}

func splitToArray(tag *lexer.Tag) ([]string, error) {
	items := bytes.FieldsFunc(tag.Data, func(r rune) bool { return r == ',' })
	ret := make([]string, 0, len(items))

	for _, item := range items {
		ret = append(ret, string(bytes.TrimSpace(item)))
	}

	sort.Strings(ret)
	for i := 1; i < len(ret); i++ {
		if ret[i] == ret[i-1] {
			return nil, tag.ErrInvalidFormat() // 重复的名称
		}
	}

	return ret, nil
}

// 解析  @apiDeprecated 标签，格式如下：
//  @apiDeprecated description
func (api *API) parseDeprecated(l *lexer.Lexer, tag *lexer.Tag) error {
	if api.Deprecated != "" {
		return tag.ErrDuplicateTag()
	}

	if len(tag.Data) == 0 {
		return tag.ErrInvalidFormat()
	}

	api.Deprecated = string(tag.Data)
	return nil
}

// 解析 @apiQuery 标签，格式如下：
//  @apiQuery name type.subtype optional.defaultValue markdown desc
func (api *API) parseQuery(l *lexer.Lexer, tag *lexer.Tag) error {
	if api.Queries == nil {
		api.Queries = make([]*Param, 0, 10)
	}

	p, err := newParam(tag)
	if err != nil {
		return err
	}

	if api.queryExists(p.Name) {
		return tag.ErrDuplicateTag()
	}

	api.Queries = append(api.Queries, p)
	return nil
}

func (api *API) queryExists(name string) bool {
	for _, q := range api.Queries {
		if q.Name == name {
			return true
		}
	}

	return false
}

// 解析 @apiParam 标签，格式如下：
//  @apiParam name type.subtype optional.defaultValue markdown desc
func (api *API) parseParam(l *lexer.Lexer, tag *lexer.Tag) error {
	if api.Params == nil {
		api.Params = make([]*Param, 0, 3)
	}

	p, err := newParam(tag)
	if err != nil {
		return err
	}

	if api.paramExists(p.Name) {
		return tag.ErrDuplicateTag()
	}

	api.Params = append(api.Params, p)

	return nil
}

func (api *API) paramExists(name string) bool {
	for _, p := range api.Params {
		if p.Name == name {
			return true
		}
	}

	return false
}

// 解析 @apiRequest 及其子标签，格式如下：
//  @apirequest object * 通用的请求主体
//  @apiheader name optional desc
//  @apiheader name optional desc
//  @apiparam count int optional desc
//  @apiparam list array.string optional desc
//  @apiparam list.id int optional desc
//  @apiparam list.name int reqiured desc
//  @apiparam list.groups array.string optional.xxxx desc markdown enum:
//   * xx: xxxxx
//   * xx: xxxxx
//  @apiexample application/json summary
//  {
//   count: 5,
//   list: [
//     {id:1, name: 'name1', 'groups': [1,2]},
//     {id:2, name: 'name2', 'groups': [1,2]}
//   ]
//  }
func (api *API) parseRequest(l *lexer.Lexer, tag *lexer.Tag) error {
	data := tag.Words(3)
	if len(data) < 2 {
		return tag.ErrInvalidFormat()
	}

	if api.Requests == nil {
		api.Requests = make([]*Request, 0, 3)
	}

	var desc []byte
	if len(data) == 3 {
		desc = data[2]
	}

	req := &Request{
		Mimetype: string(data[1]),
		Type:     &schema.Schema{},
	}
	api.Requests = append(api.Requests, req)

	if err := req.Type.Build(tag, nil, data[0], nil, desc); err != nil {
		return err
	}

LOOP:
	for tag := l.Tag(); tag != nil; tag = l.Tag() {
		fn := req.parseExample
		switch strings.ToLower(tag.Name) {
		case "@apiexample":
			fn = req.parseExample
		case "@apiheader":
			fn = req.parseHeader
		case "@apiparam":
			fn = req.parseParam
		default:
			l.Backup(tag)
			break LOOP
		}

		if err := fn(tag); err != nil {
			return err
		}
	}

	return nil
}

// 检测内容是否都是有效的
/*
func (api *API) check() error {

}*/

// 解析参数标签，格式如下：
// 用于路径参数和查义参数，request 和 response 中的不在此解析
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
