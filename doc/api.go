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
	apiCallback
	Path        string    `yaml:"path" json:"path"`
	Description Markdown  `yaml:"description,omitempty" json:"description,omitempty"`
	Tags        []string  `yaml:"tags,omitempty" json:"tags,omitempty"`
	Deprecated  string    `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	Servers     []string  `yaml:"servers" json:"servers"`
	Callback    *Callback `yaml:"callback,omitempty" json:"callback,omitempty"`

	// 记录起始位置，方便错误定位
	file string
	line int
}

// Callback 回调内容
//
//  @apiCallback GET
//  @apiquery ...
//  @apiparam ...
type Callback struct {
	responses
	apiCallback
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
		case "@apicallback":
			parse = api.parseCallback
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

// 解析 @apiCallback 标签
func (api *API) parseCallback(l *lexer, tag *lexerTag) {
	if api.Callback != nil {
		tag.err(locale.ErrDuplicateTag)
		return
	}

	if len(tag.Data) == 0 {
		tag.err(locale.ErrRequired)
		return
	}

	data := tag.words(2)

	c := &Callback{
		apiCallback: apiCallback{
			Method: string(data[0]),
		},
	}

	if len(data) == 2 {
		c.Summary = string(data[1])
	}

LOOP:
	for tag := l.tag(); tag != nil; tag = l.tag() {
		parse := api.parseAPI
		switch strings.ToLower(tag.Name) {
		case "@apiquery":
			parse = c.parseQuery
		case "@apiparam":
			parse = c.parseParam
		case "@apirequest":
			parse = c.parseRequest
		case "@apiresponse":
			parse = c.parseResponse
		default:
			l.backup(tag)
			break LOOP
		}

		parse(l, tag)
	}

	api.Callback = c
}

// 解析 @api 标签，格式如下：
//  @api GET /path summary
func (api *API) parseapi(l *lexer, tag *lexerTag) {
	if api.Method != "" || api.Path != "" || api.Summary != "" {
		tag.err(locale.ErrDuplicateTag)
		return
	}

	lines := tag.lines(2)

	if len(lines) == 2 {
		api.Description = Markdown(lines[1])
	}

	data := splitWords(lines[0], 3)
	if len(data) != 3 {
		tag.err(locale.ErrInvalidFormat)
		return
	}

	api.Method = strings.ToUpper(string(data[0]))
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
