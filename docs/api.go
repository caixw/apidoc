// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package docs

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/caixw/apidoc/docs/syntax"
	"github.com/caixw/apidoc/locale"
)

// @api 的格式如下：
//
// @api GET /users/{id}/logs 获取用户信息
// @group g1
// @tags t1,t2
// @deprecated desc
// @query page int default desc
// @query size int default desc
// @query state array.string [normal,lock] 状态码
// @param id int desc
// @param id int desc
// @header name desc
// @header name desc
//
// @request object * 通用的请求主体
// @param count int optional desc
// @param list array must desc
// @param list.id int optional desc
// @param list.name int must desc
// @param list.groups array.string optional.xxxx desc markdown enum:
//  * xx: xxxxx
//  * xx: xxxxx
// @example application/json default
// {
//  count: 5,
//  list: [
//    {id:1, name: 'name1', 'groups': [1,2]},
//    {id:2, name: 'name2', 'groups': [1,2]}
//  ]
// }
//
// @request object application/xml 特定的请求主体
//
// @response 200 array.object 通用的返回内容定义
// @apiheader string xxx
// @param id int desc
// @param name string desc
// @param group object desc
// @param group.id int desc
//
// @response 404 object application/json
// @apiheader string xxx
// @param code int desc
// @param message string desc
// @param detail array.object desc
// @param detail.id string desc
// @param detail.message string desc

func (docs *Docs) parseAPI(l *syntax.Lexer) error {
	api := &API{}

	for tag, eof := l.Tag(); !eof; tag, eof = l.Tag() {
		switch string(bytes.ToLower(tag.Name)) {
		case "@api":
			if err := api.parseAPI(l, tag); err != nil {
				return err
			}
		case "@apirequest":
			if err := api.parseRequest(l, tag); err != nil {
				return err
			}
		case "@apiresponse":
			if err := api.parseResponse(l, tag); err != nil {
				return err
			}
		default:
			return tag.Error(locale.ErrInvalidTag, string(tag.Name))
		}
	}

	doc := docs.getDoc(api.group)
	doc.locker.Lock()
	defer doc.locker.Unlock()
	doc.Apis = append(doc.Apis, api)

	return nil
}

// 分析 @api 以及子标签
func (api *API) parseAPI(l *syntax.Lexer, tag *syntax.Tag) error {
	if api.Method != "" || api.Path != "" || api.Summary != "" {
		return tag.Error(locale.ErrDuplicateTag, "@api")
	}
	data := syntax.Split(tag.Data, 3)
	if len(data) != 3 {
		return tag.Error(locale.ErrTagArgNotEnough, "@api")
	}

	api.Method = strings.ToUpper(string(data[0])) // TODO 验证请求方法
	api.Path = string(data[1])
	api.Summary = string(data[2])

	for tag, eof := l.Tag(); !eof; tag, eof = l.Tag() {
		switch string(bytes.ToLower(tag.Name)) {
		case "@apigroup":
			if api.group != "" {
				return tag.Error(locale.ErrDuplicateTag, "@apiGroup")
			}
			api.group = string(tag.Data)
		case "@apitags":
			if len(api.Tags) > 0 {
				return tag.Error(locale.ErrDuplicateTag, "@apiTags")
			}

			data := tag.Data
			start := 0
			for {
				index := bytes.IndexByte(tag.Data, ',')

				if index <= 0 {
					api.Tags = append(api.Tags, string(data[start:]))
					break
				}

				api.Tags = append(api.Tags, string(data[start:index]))
				data = tag.Data[index+1:]
			}
		case "@apideprecated":
			api.Deprecated = string(tag.Data)
		case "@apiquery":
			if api.Params == nil {
				api.Params = make([]*Param, 0, 10)
			}

			params := syntax.Split(tag.Data, 4)
			if len(params) != 4 {
				return tag.Error(locale.ErrTagArgNotEnough, "@apiQuery")
			}

			api.Queries = append(api.Queries, &Param{
				Name:     string(params[0]),
				Summary:  string(params[3]),
				Optional: true, // TODO
				Type: &Schema{
					Type:    string(params[1]),
					Default: string(params[2]),
				},
			})
		case "@apiparam":
			if api.Params == nil {
				api.Params = make([]*Param, 0, 10)
			}

			params := syntax.Split(tag.Data, 4)
			if len(params) != 4 {
				return tag.Error(locale.ErrTagArgNotEnough, "@apiParam")
			}

			api.Params = append(api.Params, &Param{
				Name:     string(params[0]),
				Summary:  string(params[3]),
				Optional: true,
				Type: &Schema{
					Type:    string(params[1]),
					Default: string(params[2]),
				},
			})
		case "@apiheader":
			if api.Params == nil {
				api.Params = make([]*Param, 0, 10)
			}

			params := syntax.Split(tag.Data, 4)
			if len(params) != 4 {
				return tag.Error(locale.ErrTagArgNotEnough, "@apiHeader")
			}

			api.headers = append(api.headers, &Header{
				Name:     string(params[0]),
				Summary:  string(params[3]),
				Optional: true, // TODO
			})
		default:
			l.Backup(tag)
			return nil
		}
	}
	return nil
}

func (api *API) parseRequest(l *syntax.Lexer, tag *syntax.Tag) error {
	data := syntax.Split(tag.Data, 3)
	if len(data) != 2 {
		return tag.Error(locale.ErrInvalidFormat, "@apiRequest")
	}

	if api.Requests == nil {
		api.Requests = make([]*Request, 0, 3)
	}

	req := &Request{
		Mimetype: string(data[1]),
	}
	api.Requests = append(api.Requests, req)

	schema := &Schema{}
	if serr := buildSchema(schema, nil, data[1], nil, nil); serr != nil {
		serr.File = tag.File
		serr.Line = tag.Line
		return serr
	}
	req.Type = schema

	for tag, eof := l.Tag(); !eof; tag, eof = l.Tag() {
		switch string(bytes.ToLower(tag.Name)) {
		case "@apiexample":
			if req.Examples == nil {
				req.Examples = make([]*Example, 0, 3)
			}
			data := syntax.Split(tag.Data, 3)
			req.Examples = append(req.Examples, &Example{
				Mimetype: string(data[0]),
				Summary:  string(data[1]),
				Value:    string(data[2]),
			})
		case "@apiparam":
			params := syntax.Split(tag.Data, 4)
			if len(params) != 4 {
				return tag.Error(locale.ErrTagArgNotEnough, "@apiParam")
			}

			if err := buildSchema(schema, data[0], data[1], data[2], data[3]); err != nil {
				err.File = tag.File
				err.Line = tag.Line
				return err
			}
		default:
			l.Backup(tag)
			return nil
		}
	}

	return nil
}

func (api *API) parseResponse(l *syntax.Lexer, tag *syntax.Tag) error {
	data := syntax.Split(tag.Data, 3)
	if len(data) != 3 {
		return tag.Error(locale.ErrInvalidFormat, "@apiResponse")
	}

	if api.Responses == nil {
		api.Responses = make([]*Response, 10)
	}

	schema := &Schema{}
	if serr := buildSchema(schema, nil, data[1], nil, nil); serr != nil {
		serr.File = tag.File
		serr.Line = tag.Line
		return serr
	}
	resp := &Response{}
	resp.Mimetype = string(data[1])
	api.Responses = append(api.Responses, resp)

	for tag, eof := l.Tag(); !eof; tag, eof = l.Tag() {
		switch string(bytes.ToLower(tag.Name)) {
		case "@apiexample":
			if resp.Examples == nil {
				resp.Examples = make([]*Example, 0, 3)
			}
			data := syntax.Split(tag.Data, 3)
			resp.Examples = append(resp.Examples, &Example{
				Mimetype: string(data[0]),
				Summary:  string(data[1]),
				Value:    string(data[2]),
			})

		case "@apiheader":
			data := syntax.Split(tag.Data, 3)
			if len(data) != 2 {
				return tag.Error(locale.ErrInvalidFormat, "@apiHeader")
			}

			if resp.Headers == nil {
				resp.Headers = make([]*Header, 0, 3)
			}
			optional, err := strconv.ParseBool(string(data[2]))
			if err != nil {
				return &syntax.Error{
					File:        tag.File,
					Line:        tag.Line,
					MessageKey:  locale.ErrInvalidFormat,
					MessageArgs: []interface{}{"@apiHeader"},
				}
			}
			header := &Header{
				Name:     string(data[0]),
				Summary:  string(data[1]),
				Optional: optional,
			}
			resp.Headers = append(resp.Headers, header)
		case "@apiparam":
			data := syntax.Split(tag.Data, 4)
			if len(data) != 4 {
				return tag.Error(locale.ErrInvalidFormat, "@apiParam")
			}

			if serr := buildSchema(schema, data[0], data[1], data[2], data[3]); serr != nil {
				serr.File = tag.File
				serr.Line = tag.Line
				return serr
			}
		default:
			l.Backup(tag)
			return nil
		}
	}

	return nil
}
