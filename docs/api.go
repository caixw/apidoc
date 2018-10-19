// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package docs

import (
	"bytes"
	"strings"

	"github.com/caixw/apidoc/docs/syntax"
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
//
// @request object * 通用的请求主体
// @header name desc optional
// @header name desc optional
// @param count int optional desc
// @param list array must desc
// @param list.id int optional desc
// @param list.name int must desc
// @param list.groups array.string optional.xxxx desc markdown enum:
//  * xx: xxxxx
//  * xx: xxxxx
// @example application/json summary
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
// @response 200 array.object * 通用的返回内容定义
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

	group string
}

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
			return tag.ErrInvalidTag()
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
		return tag.ErrDuplicateTag()
	}
	data := tag.Split(3)
	if len(data) != 3 {
		return tag.ErrInvalidFormat()
	}

	api.Method = strings.ToUpper(string(data[0])) // TODO 验证请求方法
	api.Path = string(data[1])
	api.Summary = string(data[2])

	for tag, eof := l.Tag(); !eof; tag, eof = l.Tag() {
		switch string(bytes.ToLower(tag.Name)) {
		case "@apigroup":
			if api.group != "" {
				return tag.ErrDuplicateTag()
			}
			api.group = string(tag.Data)
		case "@apitags":
			if len(api.Tags) > 0 {
				return tag.ErrDuplicateTag()
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

			params := tag.Split(4)
			if len(params) != 4 {
				return tag.ErrInvalidFormat()
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

			params := tag.Split(4)
			if len(params) != 4 {
				return tag.ErrInvalidFormat()
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
		default:
			l.Backup(tag)
			return nil
		}
	}
	return nil
}

func (api *API) parseRequest(l *syntax.Lexer, tag *syntax.Tag) error {
	data := tag.Split(3)
	if len(data) != 3 {
		return tag.ErrInvalidFormat()
	}

	if api.Requests == nil {
		api.Requests = make([]*Request, 0, 3)
	}

	req := &Request{
		Mimetype: string(data[1]),
		Type:     &Schema{},
	}
	api.Requests = append(api.Requests, req)

	if err := buildSchema(tag, req.Type, nil, data[1], nil, data[2]); err != nil {
		return err
	}

	for tag, eof := l.Tag(); !eof; tag, eof = l.Tag() {
		switch string(bytes.ToLower(tag.Name)) {
		case "@apiexample":
			if err := req.parseExample(tag); err != nil {
				return err
			}
		case "@apiheader":
			if err := req.parseHeader(tag); err != nil {
				return err
			}
		case "@apiparam":
			params := tag.Split(4)
			if len(params) != 4 {
				return tag.ErrInvalidFormat()
			}

			if err := buildSchema(tag, req.Type, data[0], data[1], data[2], data[3]); err != nil {
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
	data := tag.Split(3)
	if len(data) != 3 {
		return tag.ErrInvalidFormat()
	}

	if api.Responses == nil {
		api.Responses = make([]*Response, 10)
	}

	schema := &Schema{}
	if err := buildSchema(tag, schema, nil, data[1], nil, data[2]); err != nil {
		return err
	}
	resp := &Response{
		Body: Body{
			Mimetype: string(data[1]),
		},
	}
	api.Responses = append(api.Responses, resp)

	for tag, eof := l.Tag(); !eof; tag, eof = l.Tag() {
		switch string(bytes.ToLower(tag.Name)) {
		case "@apiexample":
			if err := resp.parseExample(tag); err != nil {
				return err
			}
		case "@apiheader":
			if err := resp.parseHeader(tag); err != nil {
				return err
			}
		case "@apiparam":
			data := tag.Split(4)
			if len(data) != 4 {
				return tag.ErrInvalidFormat()
			}

			if err := buildSchema(tag, schema, data[0], data[1], data[2], data[3]); err != nil {
				return err
			}
		default:
			l.Backup(tag)
			return nil
		}
	}

	return nil
}
