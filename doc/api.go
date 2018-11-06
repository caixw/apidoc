// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"bytes"
	"sort"
	"strings"

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
	Name     string  `yaml:"name" json:"name"`                             // 参数名称
	Type     *Schema `yaml:"type" json:"type"`                             // 类型
	Summary  string  `yaml:"summary" json:"summary"`                       // 参数介绍
	Optional bool    `yaml:"optional,omitempty" json:"optional,omitempty"` // 是否为可选参数
}

func (doc *Doc) parseAPI(l *lexer) {
	api := &API{}

	for tag := l.tag(); tag != nil; tag = l.tag() {
		parse := api.parseAPI
		switch strings.ToLower(tag.Name) {
		case "@api":
			parse = api.parseAPI
		case "@apirequest":
			parse = api.parseRequest
		case "@apiresponse":
			parse = api.parseResponse
		default:
			tag.warn(locale.ErrInvalidTag)
			return
		}

		parse(l, tag)
	}

	doc.append(l, api)
}

func (doc *Doc) append(l *lexer, api *API) {
	doc.locker.Lock()
	defer doc.locker.Unlock()

	for _, item := range doc.Apis {
		if item.Method == api.Method && item.Path == api.Path {
			l.err(api.err("@api", locale.ErrDuplicateRoute))
		}
	}

	doc.Apis = append(doc.Apis, api)
}

type apiParser func(*API, *lexer, *lexerTag)

var apiParsers = map[string]apiParser{
	"@api":           (*API).parseapi,
	"@apiservers":    (*API).parseServers,
	"@apitags":       (*API).parseTags,
	"@apideprecated": (*API).parseDeprecated,
	"@apiquery":      (*API).parseQuery,
	"@apiparam":      (*API).parseParam,
}

// 检测内容是否正确
func (api *API) check(h *errors.Handler) {
	names, err := api.getPathParams()
	if err != nil {
		h.SyntaxWarn(err)
	}

	if len(names) != len(api.Params) {
		h.SyntaxWarn(api.err("@api", locale.ErrPathNotMatchParams))
	}

	for _, p := range api.Params {
		if !names[p.Name] {
			h.SyntaxWarn(api.err("@api", locale.ErrPathNotMatchParams))
		}
	}
}

func (api *API) getPathParams() (map[string]bool, *errors.Error) {
	names := make(map[string]bool, len(api.Params))

	state := '}'
	index := 0
	for i, b := range api.Path {
		switch b {
		case '{':
			if state != '}' {
				return nil, api.err("@api", locale.ErrPathInvalid)
			}

			state = '{'
			index = i
		case '}':
			if state != '{' {
				return nil, api.err("@api", locale.ErrPathInvalid)
			}
			names[api.Path[index+1:i]] = true
			state = '}'
		}
	} // end for

	// 缺少 } 结束符号
	if state == '{' {
		return nil, api.err("@api", locale.ErrPathInvalid)
	}

	return names, nil
}

// 分析 @api 以及子标签
func (api *API) parseAPI(l *lexer, tag *lexerTag) {
	l.backup(tag) // 进来时，第一个肯定是 @api 标签，退回该标签，统一让 for 处理。

	for tag := l.tag(); tag != nil; tag = l.tag() {
		fn, found := apiParsers[strings.ToLower(tag.Name)]
		if !found { // 未找到标签，返回给上一级
			l.backup(tag)
			return
		}

		fn(api, l, tag)
	}
}

// 解析 @api 标签，格式如下：
//  @api GET /path summary
func (api *API) parseapi(l *lexer, tag *lexerTag) {
	if api.Method != "" || api.Path != "" || api.Summary != "" {
		tag.err(locale.ErrDuplicateTag)
		return
	}

	data := tag.words(3)
	if len(data) != 3 {
		tag.err(locale.ErrInvalidFormat)
		return
	}

	api.Method = strings.ToUpper(string(data[0])) // TODO 验证请求方法
	api.Path = string(data[1])
	api.Summary = string(data[2])
	api.file = tag.File
	api.line = tag.Line
}

// 解析 @apiServers 标签，格式如下：
//  @apiServers s1,s2
func (api *API) parseServers(l *lexer, tag *lexerTag) {
	if len(api.Servers) > 0 {
		tag.err(locale.ErrDuplicateTag)
		return
	}

	if len(tag.Data) == 0 {
		tag.err(locale.ErrRequired)
		return
	}

	api.Servers = splitToArray(tag)
}

// 解析 @apiTags 标签，格式如下：
//  @apiTags t1,t2
func (api *API) parseTags(l *lexer, tag *lexerTag) {
	if len(api.Tags) > 0 {
		tag.err(locale.ErrDuplicateTag)
		return
	}

	if len(tag.Data) == 0 {
		tag.err(locale.ErrRequired)
		return
	}

	api.Tags = splitToArray(tag)
}

// 按 , 分隔内容。并去掉各个元素的首尾空格。
// 如果存在相同元素，则返回错误信息。
//
// 一般用于诸如 @apiTags 等标签的检测。
func splitToArray(tag *lexerTag) []string {
	items := bytes.FieldsFunc(tag.Data, func(r rune) bool { return r == ',' })
	ret := make([]string, 0, len(items))

	for _, item := range items {
		ret = append(ret, string(bytes.TrimSpace(item)))
	}

	sort.Strings(ret)
	for i := 1; i < len(ret); i++ {
		if ret[i] == ret[i-1] {
			tag.err(locale.ErrDuplicateValue)
			return nil
		}
	}

	return ret
}

// 解析  @apiDeprecated 标签，格式如下：
//  @apiDeprecated description
func (api *API) parseDeprecated(l *lexer, tag *lexerTag) {
	if api.Deprecated != "" {
		tag.err(locale.ErrDuplicateTag)
		return
	}

	if len(tag.Data) == 0 {
		tag.err(locale.ErrRequired)
		return
	}

	api.Deprecated = string(tag.Data)
}

// 解析 @apiQuery 标签，格式如下：
//  @apiQuery name type.subtype optional.defaultValue markdown desc
func (api *API) parseQuery(l *lexer, tag *lexerTag) {
	if api.Queries == nil {
		api.Queries = make([]*Param, 0, 10)
	}

	p, ok := newParam(tag)
	if !ok {
		return
	}

	if api.queryExists(p.Name) {
		tag.err(locale.ErrDuplicateValue)
		return
	}

	api.Queries = append(api.Queries, p)
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
func (api *API) parseParam(l *lexer, tag *lexerTag) {
	if api.Params == nil {
		api.Params = make([]*Param, 0, 3)
	}

	p, ok := newParam(tag)
	if !ok {
		return
	}

	if api.paramExists(p.Name) {
		tag.err(locale.ErrDuplicateValue)
		return
	}

	api.Params = append(api.Params, p)
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
func (api *API) parseRequest(l *lexer, tag *lexerTag) {
	data := tag.words(3)
	if len(data) < 2 {
		tag.err(locale.ErrInvalidFormat)
		return
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
		Type:     &Schema{},
	}
	api.Requests = append(api.Requests, req)

	if err := req.Type.build(nil, data[0], nil, desc); err != nil {
		tag.errWithError(err, locale.ErrInvalidFormat)
		return
	}

LOOP:
	for tag := l.tag(); tag != nil; tag = l.tag() {
		fn := req.parseExample
		switch strings.ToLower(tag.Name) {
		case "@apiexample":
			fn = req.parseExample
		case "@apiheader":
			fn = req.parseHeader
		case "@apiparam":
			fn = req.parseParam
		default:
			l.backup(tag)
			break LOOP
		}

		fn(tag)
	}
}

// 解析参数标签，格式如下：
// 用于路径参数和查义参数，request 和 response 中的不在此解析
//  @tag name type.subtype optional.defaultValue markdown desc
func newParam(tag *lexerTag) (p *Param, ok bool) {
	data := tag.words(4)
	if len(data) != 4 {
		tag.err(locale.ErrInvalidFormat)
		return nil, false
	}

	s := &Schema{}
	if err := s.build(nil, data[1], data[2], nil); err != nil {
		tag.errWithError(err, locale.ErrInvalidFormat)
		return nil, false
	}

	return &Param{
		Name:     string(data[0]),
		Summary:  string(data[3]),
		Type:     s,
		Optional: s.Default != nil,
	}, true
}
