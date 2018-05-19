// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package parser

import (
	"errors"
	"net/http"
	"strings"

	"github.com/caixw/apidoc/openapi"
)

type api struct {
	API         string               // @api 后面的内容，包含了 method, url 和 summary
	Group       string               `yaml:"group,omitempty"`
	Tags        []string             `yaml:"tags,omitempty"`
	Description openapi.Description  `yaml:"description,omitempty"`
	Deprecated  bool                 `yaml:"deprecated,omitempty"`
	OperationID string               `yaml:"operationId,omitempty" `
	Queries     map[string]string    `yaml:"queries,omitempty"`
	Params      map[string]string    `yaml:"params,omitempty"`
	Headers     map[string]string    `yaml:"header,omitempty"`
	Request     *request             `yaml:"request,omitempty"` // GET 此值可能为空
	Responses   map[string]*response `yaml:"responses"`
}

// 表示请求和返回的内容
type content struct {
	Description openapi.Description           `yaml:"description,omitempty"`
	Content     map[string]*openapi.MediaType `yaml:"content"`
}

type request content

type response struct {
	content
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
}

func (doc *doc) parseAPI(api *api) error {
	o, err := doc.getOperation(api)
	if err != nil {
		return err
	}

	o.Tags = api.Tags
	o.Description = api.Description
	o.Deprecated = api.Deprecated
	o.OperationID = api.OperationID

	if err := api.parseParameter(o); err != nil {
		return err
	}

	o.RequestBody = &openapi.RequestBody{
		Description: api.Request.Description,
		Content:     api.Request.Content,
	}

	o.Responses = make(map[string]*openapi.Response, len(api.Responses))
	for status, resp := range api.Responses {
		r := &openapi.Response{
			Description: resp.Description,
			Content:     resp.content.Content,
			Headers:     make(map[string]*openapi.Header, len(resp.Headers)),
		}
		for k, v := range resp.Headers {
			r.Headers[k] = &openapi.Header{Description: openapi.Description(v)}
		}

		o.Responses[status] = r
	}

	return nil
}

func (api *api) parseParameter(o *openapi.Operation) error {
	l := len(api.Queries) + len(api.Params) + len(api.Headers)
	o.Parameters = make([]*openapi.Parameter, 0, l)

	// queries
	for key, query := range api.Queries {
		queries := strings.SplitN(query, " ", 3)
		o.Parameters = append(o.Parameters, &openapi.Parameter{
			Name:        key,
			IN:          openapi.ParameterINQuery,
			Description: openapi.Description(queries[2]),
			Schema: &openapi.Schema{
				Type:    queries[0], // TODO 判断是否正确
				Default: queries[1], // TODO 判断是否可以和类型匹配
			},
		})
	}

	// params
	for key, param := range api.Params {
		params := strings.SplitN(param, " ", 2)
		o.Parameters = append(o.Parameters, &openapi.Parameter{
			Name:        key,
			IN:          openapi.ParameterINPath,
			Description: openapi.Description(params[1]),
			Schema: &openapi.Schema{
				Type: params[0], // TODO 判断是否正确
			},
		})
	}

	// headers
	for key, desc := range api.Headers {
		o.Parameters = append(o.Parameters, &openapi.Parameter{
			Name:        key,
			IN:          openapi.ParameterINHeader,
			Description: openapi.Description(desc),
			Schema: &openapi.Schema{
				Type: openapi.TypeString,
			},
		})
	}

	return nil
}

func (doc *doc) getOperation(api *api) (*openapi.Operation, error) {
	doc.locker.Lock()
	defer doc.locker.Unlock()

	if doc.OpenAPI.Paths == nil {
		doc.OpenAPI.Paths = make(map[string]*openapi.PathItem, 10)
	}

	strs := strings.SplitN(api.API, " ", 3)
	if len(strs) != 3 {
		return nil, errors.New("缺少参数")
	}

	path, found := doc.OpenAPI.Paths[strs[1]]
	if !found {
		path = &openapi.PathItem{}
		doc.OpenAPI.Paths[strs[1]] = path
	}

	switch strings.ToUpper(strs[0]) {
	case http.MethodGet:
		if path.Get != nil {
			return nil, errors.New("已经存在一个相同的 Get 路由")
		}
		path.Get = &openapi.Operation{}
		return path.Get, nil
	case http.MethodDelete:
		if path.Delete != nil {
			return nil, errors.New("已经存在一个相同的 Delete 路由")
		}
		path.Delete = &openapi.Operation{}
		return path.Delete, nil
	case http.MethodPost:
		if path.Post != nil {
			return nil, errors.New("已经存在一个相同的 Post 路由")
		}
		path.Post = &openapi.Operation{}
		return path.Post, nil
	case http.MethodPut:
		if path.Put != nil {
			return nil, errors.New("已经存在一个相同的 Put 路由")
		}
		path.Put = &openapi.Operation{}
		return path.Put, nil
	case http.MethodPatch:
		if path.Patch != nil {
			return nil, errors.New("已经存在一个相同的 Patch 路由")
		}
		path.Patch = &openapi.Operation{}
		return path.Patch, nil
	case http.MethodOptions:
		if path.Options != nil {
			return nil, errors.New("已经存在一个相同的 Options 路由")
		}
		path.Options = &openapi.Operation{}
		return path.Options, nil
	case http.MethodHead:
		if path.Head != nil {
			return nil, errors.New("已经存在一个相同的 Head 路由")
		}
		path.Head = &openapi.Operation{}
		return path.Head, nil
	case http.MethodTrace:
		if path.Trace != nil {
			return nil, errors.New("已经存在一个相同的 Trace 路由")
		}
		path.Trace = &openapi.Operation{}
		return path.Trace, nil
	}

	return nil, errors.New("无效的请法语方法")
}
