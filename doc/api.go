// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"github.com/caixw/apidoc/app"
	"github.com/issue9/is"
)

// Scan 扫描文档，生成一个 Doc 实例。
//
// 若代码块没有 api 文档定义，则会返回空值。
// block 该代码块的内容；
func (d *Doc) Scan(data []rune) *app.SyntaxError {
	l := newLexer(data)

LOOP:
	for {
		switch {
		case l.matchTag("@apidoc"):
			return l.scanAPIDoc(d)
		case l.matchTag("@api"):
			return l.scanAPI(d)
		case l.match("@api"): // 不认识标签
			l.backup()
			return l.syntaxError("不认识的顶层标签" + l.readWord())
		default:
			if l.atEOF() {
				break LOOP
			}
			l.pos++ // 去掉无用的字符。
		}
	} // end for

	return nil
}

// 检测变量 api 是否为空值。
func apiIsEmpty(api *API) bool {
	return api == nil || (len(api.Method) == 0 &&
		len(api.URL) == 0 &&
		len(api.Summary) == 0 &&
		len(api.Description) == 0 &&
		len(api.Group) == 0 &&
		len(api.Queries) == 0 &&
		len(api.Params) == 0 &&
		api.Request == nil &&
		api.Success == nil &&
		api.Error == nil)
}

// 检测 api 的所有基本要素是否齐全。
//
// NOTE: scan* 系列函数负责解析标签，及该标签是否合法，
// 但若整个标签缺失则无能为力，此即 checkAPI 的存在的作用。
func checkAPI(api *API) *app.SyntaxError {
	switch {
	case len(api.URL) == 0 || len(api.Method) == 0:
		return &app.SyntaxError{Message: "缺少必要的元素 @api"}
	case api.Success == nil && api.Error == nil:
		return &app.SyntaxError{Message: "@apiSuccess @apiError 必须得有一个"}
	default:
		return nil
	}
}

// @apidoc title of doc
// @apiVersion 2.0
// @apiBaseURL https://api.caixw.io
// @apiLicense MIT https://opensource.org/licenses/MIT
//
// @apiContent
// content1
// content2
func (l *lexer) scanAPIDoc(d *Doc) *app.SyntaxError {
	if len(d.Title) > 0 || len(d.Version) > 0 {
		return l.syntaxError("重复的 @apidoc 标签:title=" + d.Title + ",Version=" + d.Version)
	}

	t := l.readTag()
	d.Title = t.readLine()
	if len(d.Title) == 0 {
		return l.syntaxError("@apidoc 未指定标题")
	}
	if !t.atEOF() {
		return l.syntaxError("@apidoc 过多的参数")
	}

LOOP:
	for {
		switch {
		case l.matchTag("@apiVersion"):
			t := l.readTag()
			d.Version = t.readLine()
			if len(d.Version) == 0 {
				return t.syntaxError("@apiVersion 未指定参数")
			}
			if !t.atEOF() {
				return t.syntaxError("@apiVersion 过多的参数")
			}
		case l.matchTag("@apiBaseURL"):
			t := l.readTag()
			d.BaseURL = t.readLine()
			if len(d.BaseURL) == 0 {
				return t.syntaxError("@apiBaseURL 未指定参数")
			}
			if !t.atEOF() {
				return t.syntaxError("@apiBaseURL 过多的参数")
			}
		case l.matchTag("@apiLicense"):
			t := l.readTag()
			d.LicenseName = t.readWord()
			d.LicenseURL = t.readLine()
			if len(d.LicenseName) == 0 {
				return t.syntaxError("@apiLicense 缺少必要的参数")
			}
			if len(d.LicenseURL) > 0 && !is.URL(d.LicenseURL) {
				return t.syntaxError("@apiLicense 第二个参数必须为一个 URL")
			}
			if !t.atEOF() {
				return t.syntaxError("@apiLicense 过多的参数")
			}
		case l.matchTag("@apiContent"):
			d.Content = string(l.data[l.pos:])
		case l.match("@api"):
			return l.syntaxError("不认识的标签" + l.readWord())
		default:
			if l.atEOF() {
				break LOOP
			}
			l.pos++ // 去掉无用的字符。
		}
	}

	return nil
}

func (l *lexer) scanGroup(api *API) *app.SyntaxError {
	t := l.readTag()

	api.Group = t.readWord()
	if len(api.Group) == 0 {
		return t.syntaxError("@apiGroup 未指定名称")
	}

	if !t.atEOF() {
		t.syntaxError("@apiGroup 参数过多")
	}

	return nil
}

func (l *lexer) scanAPIQueries(api *API) *app.SyntaxError {
	if api.Queries == nil {
		api.Queries = make([]*Param, 0, 1)
	}

	p, err := l.scanAPIParam()
	if err != nil {
		return err
	}
	api.Queries = append(api.Queries, p)
	return nil
}

func (l *lexer) scanAPIParams(api *API) *app.SyntaxError {
	if api.Params == nil {
		api.Params = make([]*Param, 0, 1)
	}

	p, err := l.scanAPIParam()
	if err != nil {
		return err
	}
	api.Params = append(api.Params, p)
	return nil
}

// @apiRequest json,xml
func (l *lexer) scanAPIRequest(api *API) *app.SyntaxError {
	t := l.readTag()
	r := &Request{
		Type:     t.readLine(),
		Headers:  map[string]string{},
		Params:   []*Param{},
		Examples: []*Example{},
	}
	if !t.atEOF() {
		return t.syntaxError("@apiRequest 过多的参数:" + t.readEnd())
	}

LOOP:
	for {
		switch {
		case l.matchTag("@apiHeader"):
			t := l.readTag()
			key := t.readWord()
			val := t.readLine()
			if len(key) == 0 || len(val) == 0 {
				return t.syntaxError("@apiHeader 缺少必要的参数")
			}
			if !t.atEOF() {
				return t.syntaxError("@apiHeader 参数过多")
			}
			r.Headers[string(key)] = string(val)
		case l.matchTag("@apiParam"):
			p, err := l.scanAPIParam()
			if err != nil {
				return err
			}
			r.Params = append(r.Params, p)
		case l.matchTag("@apiExample"):
			e, err := l.scanAPIExample()
			if err != nil {
				return err
			}
			r.Examples = append(r.Examples, e)
		case l.match("@api"): // 其它api*，退出。
			l.backup()
			break LOOP
		default:
			if l.atEOF() {
				break LOOP
			}
			l.pos++ // 去掉无用的字符。

		} // end switch
	} // end for

	api.Request = r
	return nil
}

func (l *lexer) scanResponse() (*Response, *app.SyntaxError) {
	tag := l.readTag()
	resp := &Response{
		Code:     tag.readWord(),
		Summary:  tag.readLine(),
		Headers:  map[string]string{},
		Params:   []*Param{},
		Examples: []*Example{},
	}

	if len(resp.Code) == 0 || len(resp.Summary) == 0 {
		return nil, tag.syntaxError("@apiSuccess 或是 @apiError 缺少必要的元素")
	}
	if !tag.atEOF() {
		return nil, tag.syntaxError("@apiSuccess 或是 @apiError 参数过多")
	}

LOOP:
	for {
		switch {
		case l.matchTag("@apiHeader"):
			t := l.readTag()
			key := t.readWord()
			val := t.readLine()
			if len(key) == 0 || len(val) == 0 {
				return nil, t.syntaxError("@apiHeader 缺少必要的参数")
			}
			if !t.atEOF() {
				return nil, t.syntaxError("@apiHeader 参数过多")
			}
			resp.Headers[key] = val
		case l.matchTag("@apiParam"):
			p, err := l.scanAPIParam()
			if err != nil {
				return nil, err
			}
			resp.Params = append(resp.Params, p)
		case l.matchTag("@apiExample"):
			e, err := l.scanAPIExample()
			if err != nil {
				return nil, err
			}
			resp.Examples = append(resp.Examples, e)
		case l.match("@api"): // 其它api*，退出。
			l.backup()
			break LOOP
		default:
			if l.atEOF() {
				break LOOP
			}
			l.pos++ // 去掉无用的字符。
		}
	}

	return resp, nil
}

func (l *lexer) scanAPIExample() (*Example, *app.SyntaxError) {
	tag := l.readTag()
	example := &Example{
		Type: tag.readWord(),
		Code: tag.readEnd(),
	}

	if len(example.Type) == 0 || len(example.Code) == 0 {
		return nil, tag.syntaxError("@apiExample 缺少必要的参数")
	}

	return example, nil
}

func (l *lexer) scanAPIParam() (*Param, *app.SyntaxError) {
	p := &Param{}

	tag := l.readTag()
	p.Name = tag.readWord()
	p.Type = tag.readWord()
	p.Summary = tag.readEnd()
	if len(p.Name) == 0 || len(p.Type) == 0 || len(p.Summary) == 0 {
		return nil, tag.syntaxError("@apiParam 或是 @apiQuery 缺少必要的参数")
	}
	return p, nil
}

// 若存在 description 参数，会原样输出，不会像其它一样去掉前导空格。
// 扫描以下格式内容：
//  @api get /test.com/api/user.json api summary
//  api description
//  api description
func (l *lexer) scanAPI(d *Doc) (err *app.SyntaxError) {
	api := &API{}
	t := l.readTag()
	api.Method = t.readWord()
	api.URL = t.readWord()
	api.Summary = t.readLine()

	if len(api.Method) == 0 || len(api.URL) == 0 || len(api.Summary) == 0 {
		return t.syntaxError("@api 缺少必要的参数")
	}

	api.Description = t.readEnd()

	ignore := false
LOOP:
	for {
		switch {
		case l.matchTag("@apiIgnore"):
			ignore = true
			break LOOP
		case l.matchTag("@apiGroup"):
			err = l.scanGroup(api)
		case l.matchTag("@apiQuery"):
			err = l.scanAPIQueries(api)
		case l.matchTag("@apiParam"):
			err = l.scanAPIParams(api)
		case l.matchTag("@apiRequest"):
			err = l.scanAPIRequest(api)
		case l.matchTag("@apiError"):
			api.Error, err = l.scanResponse()
		case l.matchTag("@apiSuccess"):
			api.Success, err = l.scanResponse()
		case l.match("@api"): // 不认识的标签，抛给外层解决
			return l.syntaxError("不认识的标签" + l.readWord())
		default:
			if l.atEOF() {
				break LOOP
			}
			l.pos++ // 去掉无用的字符。
		}

		if err != nil {
			return err
		}
	}

	if ignore {
		return nil
	}
	if err := checkAPI(api); err != nil {
		return err
	}

	if len(api.Group) == 0 {
		api.Group = app.DefaultGroupName
	}

	d.mux.Lock()
	d.Apis = append(d.Apis, api)
	d.mux.Unlock()
	return nil
}
