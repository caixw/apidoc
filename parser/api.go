// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package parser

import (
	"bytes"
	"strings"

	"github.com/issue9/version"

	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/openapi"
)

// @api 的格式如下：
//
// @api GET /users/{id}/logs 获取用户信息
// @group g1
// @tags t1,t2
// @version 1.0
// @deprecated desc
// @query page int default desc
// @query size int default desc
// @query state array.string [normal,lock] 状态码
// @param id int desc
// @param id int desc
// @header name desc
// @header name desc
//
// @request application/json object
// @param count int optional desc
// @param list array must desc
// @param list.id int optional desc
// @param list.name int must desc
// @param list.groups array.string.enum optional desc markdown enum{normal:xx,locked:xx}
// @example
// {
//  count: 5,
//  list: [
//    {id:1, name: 'name1', 'groups': [1,2]},
//    {id:2, name: 'name2', 'groups': [1,2]}
//  ]
// }
//
// @request application/yaml object
//
// @response 200 application/json array
// @apiheader string xxx
// @param id int desc
// @param name string desc
// @param group object desc
// @param group.id int desc
//
// @response 404 application/json object
// @apiheader string xxx
// @param code int desc
// @param message string desc
// @param detail array.object desc
// @param detail.id string desc
// @param detail.message string desc

type api struct {
	method      string
	path        string
	summary     string
	description string
	group       string
	tags        []string
	version     string
	deprecated  bool
	params      []*openapi.Parameter // 包含 query 和 param

	request   *openapi.RequestBody
	responses map[string]*openapi.Response
}

func (p *parser) parseAPI(l *lexer) error {
	obj := &api{}
	for tag, eof := l.tag(); !eof; tag, eof = l.tag() {
		switch string(bytes.ToLower(tag.name)) {
		case "@api":
			if obj.method != "" || obj.path != "" || obj.summary != "" {
				return tag.syntaxError(locale.ErrDuplicateTag, "@api")
			}
			data := split(tag.data, 3)
			if len(data) != 3 {
				return tag.syntaxError(locale.ErrTagArgNotEnough, "@api")
			}

			obj.method = strings.ToUpper(string(data[0])) // TODO 验证请求方法
			obj.path = string(data[1])
			obj.summary = string(data[2])

			if err := obj.parseAPI(l); err != nil {
				return err
			}
		case "@apirequest":
			data := split(tag.data, 2)
			if len(data) != 2 {
				return tag.syntaxError(locale.ErrInvalidFormat, "@apiRequest")
			}

			if err := obj.parseRequest(l, tag); err != nil {
				return err
			}
		case "@apiresponse":
			// TODO
		default:
			return tag.syntaxError(locale.ErrInvalidTag, string(tag.name))
		}
	}
	return nil
}

func (obj *api) parseRequest(l *lexer, tag *tag) error {
	data := split(tag.data, 2)
	if len(data) != 2 {
		return tag.syntaxError(locale.ErrInvalidFormat, "@apiRequest")
	}

	if obj.request == nil {
		obj.request = &openapi.RequestBody{
			Content: make(map[string]*openapi.MediaType, 3),
		}
	}

	var typ, subtype string
	index := bytes.IndexByte(data[1], ',')
	if index > 0 {
		t1 := data[1][:index]
		subtype = string(data[1][index+1:])
		typ = string(t1)

		if typ != "array" {
			return tag.syntaxError(locale.ErrInvalidFormat, "@apiRequest")
		}
	}

	mimetype := string(data[0])
	obj.request.Content[mimetype] = &openapi.MediaType{
		Schema: &openapi.Schema{
			Type: typ,
		},
	}

	schema := obj.request.Content[mimetype].Schema

	if subtype != "" {
		schema.Items.Type = subtype
	}

	for tag, eof := l.tag(); !eof; tag, eof = l.tag() {
		switch string(bytes.ToLower(tag.name)) {
		case "@apiexample":
			obj.request.Content[mimetype].Example = openapi.ExampleValue(string(tag.data))
		case "@apiparam":
			// TODO
		default:
			l.backup(tag)
			return nil
		}
	}

	return nil
}

func (obj *api) parseAPI(l *lexer) error {
	for tag, eof := l.tag(); !eof; tag, eof = l.tag() {
		switch string(bytes.ToLower(tag.name)) {
		case "@apigroup":
			if obj.group != "" {
				return tag.syntaxError(locale.ErrDuplicateTag, "@apiGroup")
			}
			obj.group = string(tag.data)
		case "@apitags":
			if len(obj.tags) > 0 {
				return tag.syntaxError(locale.ErrDuplicateTag, "@apiTags")
			}

			data := tag.data
			start := 0
			for {
				index := bytes.IndexByte(tag.data, ',')

				if index <= 0 {
					obj.tags = append(obj.tags, string(data[start:]))
					break
				}

				obj.tags = append(obj.tags, string(data[start:index]))
				data = tag.data[index+1:]
			}
		case "@apiversion":
			if obj.version != "" {
				return tag.syntaxError(locale.ErrDuplicateTag, "@apiVersion")
			}
			obj.version = string(tag.data)

			if !version.SemVerValid(obj.version) {
				return tag.syntaxError(locale.ErrInvalidFormat, "@apiVersion")
			}
		case "@apideprecated":
			// TODO 输出警告信息
			obj.deprecated = true
		case "@apiquery":
			if obj.params == nil {
				obj.params = make([]*openapi.Parameter, 0, 10)
			}
			// TODO 复杂类型的检测 @apiquery state string.array.enum default desc
			params := split(tag.data, 4)
			if len(params) != 4 {
				return tag.syntaxError(locale.ErrTagArgNotEnough, "@apiQuery")
			}

			obj.params = append(obj.params, &openapi.Parameter{
				Name:            string(params[0]),
				IN:              openapi.ParameterINQuery,
				Description:     openapi.Description(params[3]),
				Required:        false,
				AllowEmptyValue: true,
				Schema: &openapi.Schema{
					Type:    string(params[1]), // TODO 检测类型是否符合 openapi 要求
					Default: string(params[2]),
				},
			})
		case "@apiparam":
			if obj.params == nil {
				obj.params = make([]*openapi.Parameter, 0, 10)
			}

			params := split(tag.data, 4)
			if len(params) != 4 {
				return tag.syntaxError(locale.ErrTagArgNotEnough, "@apiParam")
			}

			obj.params = append(obj.params, &openapi.Parameter{
				Name:        string(params[0]),
				IN:          openapi.ParameterINPath,
				Description: openapi.Description(params[3]),
				Required:    true,
				Schema: &openapi.Schema{
					Type:    string(params[1]), // TODO 检测类型是否符合 openapi 要求
					Default: string(params[2]),
				},
			})
		case "@apiheader":
			if obj.params == nil {
				obj.params = make([]*openapi.Parameter, 0, 10)
			}

			params := split(tag.data, 4)
			if len(params) != 4 {
				return tag.syntaxError(locale.ErrTagArgNotEnough, "@apiHeader")
			}

			obj.params = append(obj.params, &openapi.Parameter{
				Name:            string(params[0]),
				IN:              openapi.ParameterINHeader,
				Description:     openapi.Description(params[3]),
				Required:        false,
				AllowEmptyValue: true,
			})
		default:
			l.backup(tag)
			return nil
		}
	}
	return nil
}
