// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package syntax 负责对标签语法的解析操作。
package syntax

import (
	"log"

	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/types"
	"github.com/caixw/apidoc/vars"

	"github.com/issue9/is"
)

// Input 由外界提供的与标签语法相关的内容。
type Input struct {
	File  string      // 该段代码所在的文件
	Line  int         // 该段代码在文件中的行号
	Data  []rune      // 需要解析的代码段
	Error *log.Logger // 出错时的输出通道
	Warn  *log.Logger // 警告信息的输出通道
}

// Parse 分析一段代码，并将结果保存到 d 中。
func Parse(input *Input, d *types.Doc) {
	l := newLexer(input)

	for {
		switch {
		case l.matchTag(vars.APIDoc):
			if !l.scanAPIDoc(d) {
				return
			}
		case l.matchTag(vars.API):
			if api, ok := l.scanAPI(); ok {
				if api == nil { // @apiIgnore
					return
				}
				d.NewAPI(api)
			} else {
				return
			}
		case l.match(vars.API):
			l.backup()
			l.syntaxWarn(locale.ErrUnknownTag, l.readWord())
			l.readTag() // 指针移到下一个标签处
		default:
			if l.atEOF() {
				return
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
func (l *lexer) scanAPIDoc(d *types.Doc) bool {
	if len(d.Title) > 0 || len(d.Version) > 0 {
		l.syntaxError(locale.ErrDuplicateTag, vars.APIDoc)
		return false
	}

	t := l.readTag()
	d.Title = t.readLine()
	if len(d.Title) == 0 {
		l.syntaxError(locale.ErrTagArgNotEnough, vars.APIDoc)
		return false
	}
	if !t.atEOF() {
		l.syntaxError(locale.ErrTagArgTooMuch, vars.APIDoc)
		return false
	}

	for {
		switch {
		case l.matchTag(vars.APIVersion):
			t := l.readTag()
			d.Version = t.readLine()
			if len(d.Version) == 0 {
				t.syntaxError(locale.ErrTagArgNotEnough, vars.APIVersion)
				return false
			}
			if !t.atEOF() {
				t.syntaxError(locale.ErrTagArgTooMuch, vars.APIVersion)
				return false
			}
		case l.matchTag(vars.APIBaseURL):
			t := l.readTag()
			d.BaseURL = t.readLine()
			if len(d.BaseURL) == 0 {
				t.syntaxError(locale.ErrTagArgNotEnough, vars.APIBaseURL)
				return false
			}
			if !t.atEOF() {
				t.syntaxError(locale.ErrTagArgTooMuch, vars.APIBaseURL)
				return false
			}
		case l.matchTag(vars.APILicense):
			t := l.readTag()
			d.LicenseName = t.readWord()
			d.LicenseURL = t.readLine()
			if len(d.LicenseName) == 0 {
				t.syntaxError(locale.ErrTagArgNotEnough, vars.APILicense)
				return false
			}
			if len(d.LicenseURL) > 0 && !is.URL(d.LicenseURL) {
				t.syntaxError(locale.ErrSecondArgMustURL)
				return false
			}
			if !t.atEOF() {
				t.syntaxError(locale.ErrTagArgTooMuch, vars.APILicense)
				return false
			}
		case l.matchTag(vars.APIContent):
			d.Content = l.readEnd()
		case l.match(vars.API): // 不认识的标签
			l.backup()
			l.syntaxError(locale.ErrUnknownTag, l.readWord())
			return false
		default:
			if l.atEOF() {
				return true
			}
			l.pos++ // 去掉无用的字符。
		}
	} // end for
}

// 解析 @api 及其子标签
func (l *lexer) scanAPI() (*types.API, bool) {
	api := &types.API{}
	t := l.readTag()
	api.Method = t.readWord()
	api.URL = t.readWord()
	api.Summary = t.readLine()

	if len(api.Method) == 0 || len(api.URL) == 0 || len(api.Summary) == 0 {
		t.syntaxError(locale.ErrTagArgNotEnough, vars.API)
		return nil, false
	}

	api.Description = t.readEnd()
LOOP:
	for {
		switch {
		case l.matchTag(vars.APIIgnore):
			return nil, true
		case l.matchTag(vars.APIGroup):
			if !l.scanGroup(api) {
				return nil, false
			}
		case l.matchTag(vars.APIQuery):
			if !l.scanAPIQueries(api) {
				return nil, false
			}
		case l.matchTag(vars.APIParam):
			if !l.scanAPIParams(api) {
				return nil, false
			}
		case l.matchTag(vars.APIRequest):
			req, ok := l.scanAPIRequest()
			if !ok {
				return nil, false
			}
			api.Request = req
		case l.matchTag(vars.APIError):
			resp, ok := l.scanResponse(vars.APIError)
			if !ok {
				return nil, false
			}
			api.Error = resp
		case l.matchTag(vars.APISuccess):
			resp, ok := l.scanResponse(vars.APISuccess)
			if !ok {
				return nil, false
			}
			api.Success = resp
		case l.match(vars.API): // 不认识的标签
			l.backup()
			l.syntaxWarn(locale.ErrUnknownTag, l.readWord())
			l.readTag() // 指针移到下一个标签处
		default:
			if l.atEOF() {
				break LOOP
			}
			l.pos++ // 去掉无用的字符。
		}
	} // end for

	if api.Success == nil {
		l.syntaxError(locale.ErrSuccessNotEmpty)
		return nil, false
	}

	if len(api.Group) == 0 {
		api.Group = vars.DefaultGroupName
	}

	return api, true
}

func (l *lexer) scanGroup(api *types.API) bool {
	t := l.readTag()

	api.Group = t.readWord()
	if len(api.Group) == 0 {
		t.syntaxError(locale.ErrTagArgNotEnough, vars.APIGroup)
		return false
	}

	if !t.atEOF() {
		t.syntaxError(locale.ErrTagArgTooMuch, vars.APIGroup)
		return false
	}

	return true
}

func (l *lexer) scanAPIQueries(api *types.API) bool {
	if api.Queries == nil {
		api.Queries = make([]*types.Param, 0, 1)
	}

	if p, ok := l.scanAPIParam(vars.APIQuery); ok {
		api.Queries = append(api.Queries, p)
		return true
	}
	return false
}

func (l *lexer) scanAPIParams(api *types.API) bool {
	if api.Params == nil {
		api.Params = make([]*types.Param, 0, 1)
	}

	if p, ok := l.scanAPIParam(vars.APIParam); ok {
		api.Params = append(api.Params, p)
		return true
	}
	return false
}

// 解析 @apiRequest 及其子标签
func (l *lexer) scanAPIRequest() (*types.Request, bool) {
	t := l.readTag()
	r := &types.Request{
		Type:     t.readLine(),
		Headers:  map[string]string{},
		Params:   []*types.Param{},
		Examples: []*types.Example{},
	}
	if !t.atEOF() {
		t.syntaxError(locale.ErrTagArgTooMuch, vars.APIRequest)
		return nil, false
	}

LOOP:
	for {
		switch {
		case l.matchTag(vars.APIHeader):
			t := l.readTag()
			key := t.readWord()
			val := t.readLine()
			if len(key) == 0 || len(val) == 0 {
				t.syntaxError(locale.ErrTagArgNotEnough, vars.APIHeader)
				return nil, false
			}
			if !t.atEOF() {
				t.syntaxError(locale.ErrTagArgTooMuch, vars.APIHeader)
				return nil, false
			}
			r.Headers[string(key)] = string(val)
		case l.matchTag(vars.APIParam):
			p, ok := l.scanAPIParam(vars.APIParam)
			if !ok {
				return nil, false
			}
			r.Params = append(r.Params, p)
		case l.matchTag(vars.APIExample):
			e, ok := l.scanAPIExample()
			if !ok {
				return nil, false
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

	return r, true
}

// 解析 @apiSuccess 或是 @apiError 及其子标签。
func (l *lexer) scanResponse(tagName string) (*types.Response, bool) {
	tag := l.readTag()
	resp := &types.Response{
		Code:     tag.readWord(),
		Summary:  tag.readLine(),
		Headers:  map[string]string{},
		Params:   []*types.Param{},
		Examples: []*types.Example{},
	}

	if len(resp.Code) == 0 || len(resp.Summary) == 0 {
		tag.syntaxError(locale.ErrTagArgNotEnough, tagName)
		return nil, false
	}
	if !tag.atEOF() {
		tag.syntaxError(locale.ErrTagArgTooMuch, tagName)
		return nil, false
	}

LOOP:
	for {
		switch {
		case l.matchTag(vars.APIHeader):
			t := l.readTag()
			key := t.readWord()
			val := t.readLine()
			if len(key) == 0 || len(val) == 0 {
				t.syntaxError(locale.ErrTagArgNotEnough, vars.APIHeader)
				return nil, false
			}
			if !t.atEOF() {
				t.syntaxError(locale.ErrTagArgTooMuch, vars.APIHeader)
				return nil, false
			}
			resp.Headers[key] = val
		case l.matchTag(vars.APIParam):
			p, ok := l.scanAPIParam(vars.APIParam)
			if !ok {
				return nil, false
			}
			resp.Params = append(resp.Params, p)
		case l.matchTag(vars.APIExample):
			e, ok := l.scanAPIExample()
			if !ok {
				return nil, false
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

	return resp, true
}

// 解析 @apiExample 标签
func (l *lexer) scanAPIExample() (*types.Example, bool) {
	tag := l.readTag()
	example := &types.Example{
		Type: tag.readWord(),
		Code: tag.readEnd(),
	}

	if len(example.Type) == 0 || len(example.Code) == 0 {
		tag.syntaxError(locale.ErrTagArgNotEnough, vars.APIExample)
		return nil, false
	}

	return example, true
}

// 解析 @apiParam 标签
func (l *lexer) scanAPIParam(tagName string) (*types.Param, bool) {
	p := &types.Param{}

	tag := l.readTag()
	p.Name = tag.readWord()
	p.Type = tag.readWord()
	p.Summary = tag.readEnd()
	if len(p.Name) == 0 || len(p.Type) == 0 || len(p.Summary) == 0 {
		tag.syntaxError(locale.ErrTagArgNotEnough, tagName)
		return nil, false
	}
	return p, true
}
