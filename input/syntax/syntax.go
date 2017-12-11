// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package syntax 提供对代码块的语法进行解析
package syntax

import (
	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/types"
	"github.com/caixw/apidoc/vars"

	"github.com/issue9/is"
)

// Parse 分析一段代码，并将结果保存到 d 中。
//
// 若代码块没有 api 文档定义，则会返回空值。
// data 该代码块的内容；
func Parse(d *types.Doc, data []rune) *types.SyntaxError {
	l := newLexer(data)

	for {
		switch {
		case l.matchTag(vars.APIDoc):
			return l.scanAPIDoc(d)
		case l.matchTag(vars.API):
			return l.scanAPI(d)
		default:
			if l.atEOF() {
				return nil
			}
			l.pos++ // 去掉无用的字符。
		}
	} // end for
}

// 解析 @apidoc 及其子标签
//
// @apidoc title of doc
// @apiVersion 2.0
// @apiBaseURL https://api.caixw.io
// @apiLicense MIT https://opensource.org/licenses/MIT
//
// @apiContent
// content1
// content2
func (l *lexer) scanAPIDoc(d *types.Doc) *types.SyntaxError {
	if len(d.Title) > 0 || len(d.Version) > 0 {
		return l.syntaxError(locale.ErrDuplicateTag, vars.APIDoc)
	}

	t := l.readTag()
	d.Title = t.readLine()
	if len(d.Title) == 0 {
		return l.syntaxError(locale.ErrTagArgNotEnough, vars.APIDoc)
	}
	if !t.atEOF() {
		return l.syntaxError(locale.ErrTagArgTooMuch, vars.APIDoc)
	}

	for {
		switch {
		case l.matchTag(vars.APIVersion):
			t := l.readTag()
			d.Version = t.readLine()
			if len(d.Version) == 0 {
				return t.syntaxError(locale.ErrTagArgNotEnough, vars.APIVersion)
			}
			if !t.atEOF() {
				return t.syntaxError(locale.ErrTagArgTooMuch, vars.APIVersion)
			}
		case l.matchTag(vars.APIBaseURL):
			t := l.readTag()
			d.BaseURL = t.readLine()
			if len(d.BaseURL) == 0 {
				return t.syntaxError(locale.ErrTagArgNotEnough, vars.APIBaseURL)
			}
			if !t.atEOF() {
				return t.syntaxError(locale.ErrTagArgTooMuch, vars.APIBaseURL)
			}
		case l.matchTag(vars.APILicense):
			t := l.readTag()
			d.LicenseName = t.readWord()
			d.LicenseURL = t.readLine()
			if len(d.LicenseName) == 0 {
				return t.syntaxError(locale.ErrTagArgNotEnough, vars.APILicense)
			}
			if len(d.LicenseURL) > 0 && !is.URL(d.LicenseURL) {
				return t.syntaxError(locale.ErrSecondArgMustURL)
			}
			if !t.atEOF() {
				return t.syntaxError(locale.ErrTagArgTooMuch, vars.APILicense)
			}
		case l.matchTag(vars.APIContent):
			d.Content = string(l.data[l.pos:])
		case l.match(vars.API): // 不认识的标签
			l.backup()
			return l.syntaxError(locale.ErrUnknownTag, l.readWord())
		default:
			if l.atEOF() {
				return nil
			}
			l.pos++ // 去掉无用的字符。
		}
	} // end for
}

// 解析 @api 及其子标签
func (l *lexer) scanAPI(d *types.Doc) (err *types.SyntaxError) {
	api := &types.API{}
	t := l.readTag()
	api.Method = t.readWord()
	api.URL = t.readWord()
	api.Summary = t.readLine()

	if len(api.Method) == 0 || len(api.URL) == 0 || len(api.Summary) == 0 {
		return t.syntaxError(locale.ErrTagArgNotEnough, vars.API)
	}

	api.Description = t.readEnd()
LOOP:
	for {
		switch {
		case l.matchTag(vars.APIIgnore):
			return nil
		case l.matchTag(vars.APIGroup):
			err = l.scanGroup(api)
		case l.matchTag(vars.APIQuery):
			err = l.scanAPIQueries(api)
		case l.matchTag(vars.APIParam):
			err = l.scanAPIParams(api)
		case l.matchTag(vars.APIRequest):
			err = l.scanAPIRequest(api)
		case l.matchTag(vars.APIError):
			api.Error, err = l.scanResponse(vars.APIError)
		case l.matchTag(vars.APISuccess):
			api.Success, err = l.scanResponse(vars.APISuccess)
		case l.match(vars.API): // 不认识的标签
			l.backup()
			return l.syntaxError(locale.ErrUnknownTag, l.readWord())
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

	if api.Success == nil {
		return &types.SyntaxError{Message: locale.ErrSuccessNotEmpty}
	}

	if len(api.Group) == 0 {
		api.Group = vars.DefaultGroupName
	}

	d.NewAPI(api)
	return nil
}

func (l *lexer) scanGroup(api *types.API) *types.SyntaxError {
	t := l.readTag()

	api.Group = t.readWord()
	if len(api.Group) == 0 {
		return t.syntaxError(locale.ErrTagArgNotEnough, vars.APIGroup)
	}

	if !t.atEOF() {
		t.syntaxError(locale.ErrTagArgTooMuch, vars.APIGroup)
	}

	return nil
}

func (l *lexer) scanAPIQueries(api *types.API) *types.SyntaxError {
	if api.Queries == nil {
		api.Queries = make([]*types.Param, 0, 1)
	}

	p, err := l.scanAPIParam(vars.APIQuery)
	if err != nil {
		return err
	}
	api.Queries = append(api.Queries, p)
	return nil
}

func (l *lexer) scanAPIParams(api *types.API) *types.SyntaxError {
	if api.Params == nil {
		api.Params = make([]*types.Param, 0, 1)
	}

	p, err := l.scanAPIParam(vars.APIParam)
	if err != nil {
		return err
	}
	api.Params = append(api.Params, p)
	return nil
}

// 解析 @apiRequest 及其子标签
func (l *lexer) scanAPIRequest(api *types.API) *types.SyntaxError {
	t := l.readTag()
	r := &types.Request{
		Type:     t.readLine(),
		Headers:  map[string]string{},
		Params:   []*types.Param{},
		Examples: []*types.Example{},
	}
	if !t.atEOF() {
		return t.syntaxError(locale.ErrTagArgTooMuch, vars.APIRequest)
	}

LOOP:
	for {
		switch {
		case l.matchTag(vars.APIHeader):
			t := l.readTag()
			key := t.readWord()
			val := t.readLine()
			if len(key) == 0 || len(val) == 0 {
				return t.syntaxError(locale.ErrTagArgNotEnough, vars.APIHeader)
			}
			if !t.atEOF() {
				return t.syntaxError(locale.ErrTagArgTooMuch, vars.APIHeader)
			}
			r.Headers[string(key)] = string(val)
		case l.matchTag(vars.APIParam):
			p, err := l.scanAPIParam(vars.APIParam)
			if err != nil {
				return err
			}
			r.Params = append(r.Params, p)
		case l.matchTag(vars.APIExample):
			e, err := l.scanAPIExample()
			if err != nil {
				return err
			}
			r.Examples = append(r.Examples, e)
		case l.match(vars.API): // 其它 api*，退出。
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

// 解析 @apiSuccess 或是 @apiError 及其子标签。
func (l *lexer) scanResponse(tagName string) (*types.Response, *types.SyntaxError) {
	tag := l.readTag()
	resp := &types.Response{
		Code:     tag.readWord(),
		Summary:  tag.readLine(),
		Headers:  map[string]string{},
		Params:   []*types.Param{},
		Examples: []*types.Example{},
	}

	if len(resp.Code) == 0 || len(resp.Summary) == 0 {
		return nil, tag.syntaxError(locale.ErrTagArgNotEnough, tagName)
	}
	if !tag.atEOF() {
		return nil, tag.syntaxError(locale.ErrTagArgTooMuch, tagName)
	}

LOOP:
	for {
		switch {
		case l.matchTag(vars.APIHeader):
			t := l.readTag()
			key := t.readWord()
			val := t.readLine()
			if len(key) == 0 || len(val) == 0 {
				return nil, t.syntaxError(locale.ErrTagArgNotEnough, vars.APIHeader)
			}
			if !t.atEOF() {
				return nil, t.syntaxError(locale.ErrTagArgTooMuch, vars.APIHeader)
			}
			resp.Headers[key] = val
		case l.matchTag(vars.APIParam):
			p, err := l.scanAPIParam(vars.APIParam)
			if err != nil {
				return nil, err
			}
			resp.Params = append(resp.Params, p)
		case l.matchTag(vars.APIExample):
			e, err := l.scanAPIExample()
			if err != nil {
				return nil, err
			}
			resp.Examples = append(resp.Examples, e)
		case l.match(vars.API): // 其它 api*，退出。
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

// 解析 @apiExample 标签
func (l *lexer) scanAPIExample() (*types.Example, *types.SyntaxError) {
	tag := l.readTag()
	example := &types.Example{
		Type: tag.readWord(),
		Code: tag.readEnd(),
	}

	if len(example.Type) == 0 || len(example.Code) == 0 {
		return nil, tag.syntaxError(locale.ErrTagArgNotEnough, vars.APIExample)
	}

	return example, nil
}

// 解析 @apiParam 标签
func (l *lexer) scanAPIParam(tagName string) (*types.Param, *types.SyntaxError) {
	p := &types.Param{}

	tag := l.readTag()
	p.Name = tag.readWord()
	p.Type = tag.readWord()
	p.Summary = tag.readEnd()
	if len(p.Name) == 0 || len(p.Type) == 0 || len(p.Summary) == 0 {
		return nil, tag.syntaxError(locale.ErrTagArgNotEnough, tagName)
	}
	return p, nil
}
