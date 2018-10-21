// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package docs

import (
	"bytes"
	"strings"

	"github.com/caixw/apidoc/docs/lexer"
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

	group string
}

func (docs *Docs) parseAPI(l *lexer.Lexer) error {
	api := &API{}

	for tag, eof := l.Tag(); !eof; tag, eof = l.Tag() {
		switch strings.ToLower(tag.Name) {
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

	docs.getDoc(api.group).append(api)

	return nil
}

var separatorTag = []byte{','}

// 分析 @api 以及子标签
func (api *API) parseAPI(l *lexer.Lexer, tag *lexer.Tag) error {
	for tag, eof := l.Tag(); !eof; tag, eof = l.Tag() {
		switch strings.ToLower(tag.Name) {
		case "@api":
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
		case "@apigroup":
			if api.group != "" {
				return tag.ErrDuplicateTag()
			}
			api.group = string(tag.Data)
		case "@apitags":
			if len(api.Tags) > 0 {
				return tag.ErrDuplicateTag()
			}

			tags := bytes.Split(tag.Data, separatorTag)
			api.Tags = make([]string, 0, len(tags))
			for _, tag := range tags {
				api.Tags = append(api.Tags, string(tag))
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

func (api *API) parseRequest(l *lexer.Lexer, tag *lexer.Tag) error {
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
			params := tag.Split(4)
			if len(params) != 4 {
				return tag.ErrInvalidFormat()
			}

			if err := buildSchema(tag, req.Type, params[0], params[1], params[2], params[3]); err != nil {
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
