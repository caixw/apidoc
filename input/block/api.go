// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package block

import (
	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/locale"
	"github.com/issue9/is"
)

// Scan 扫描一段代码，生成一个 API 实例。
//
// 若代码块没有 api 文档定义，则会返回空值。
// data 该代码块的内容；
func Scan(d *doc.Doc, data []rune) *app.SyntaxError {
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
			return l.syntaxError(locale.ErrUnknownTopTag, l.readWord())
		default:
			if l.atEOF() {
				break LOOP
			}
			l.pos++ // 去掉无用的字符。
		}
	} // end for

	return nil
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
func (l *lexer) scanAPIDoc(d *doc.Doc) *app.SyntaxError {
	if len(d.Title) > 0 || len(d.Version) > 0 {
		return l.syntaxError(locale.ErrDuplicateTag, "@apidoc")
	}

	t := l.readTag()
	d.Title = t.readLine()
	if len(d.Title) == 0 {
		return l.syntaxError(locale.ErrTagArgNotEnough, "@apidoc")
	}
	if !t.atEOF() {
		return l.syntaxError(locale.ErrTagArgTooMuch, "@apidoc")
	}

LOOP:
	for {
		switch {
		case l.matchTag("@apiVersion"):
			t := l.readTag()
			d.Version = t.readLine()
			if len(d.Version) == 0 {
				return t.syntaxError(locale.ErrTagArgNotEnough, "@apiVersion")
			}
			if !t.atEOF() {
				return t.syntaxError(locale.ErrTagArgTooMuch, "@apiVersion")
			}
		case l.matchTag("@apiBaseURL"):
			t := l.readTag()
			d.BaseURL = t.readLine()
			if len(d.BaseURL) == 0 {
				return t.syntaxError(locale.ErrTagArgNotEnough, "@apiBaseURL")
			}
			if !t.atEOF() {
				return t.syntaxError(locale.ErrTagArgTooMuch, "@apiBaseURL")
			}
		case l.matchTag("@apiLicense"):
			t := l.readTag()
			d.LicenseName = t.readWord()
			d.LicenseURL = t.readLine()
			if len(d.LicenseName) == 0 {
				return t.syntaxError(locale.ErrTagArgNotEnough, "@apiLicense")
			}
			if len(d.LicenseURL) > 0 && !is.URL(d.LicenseURL) {
				return t.syntaxError(locale.ErrSecondArgMustURL)
			}
			if !t.atEOF() {
				return t.syntaxError(locale.ErrTagArgTooMuch, "@apiLicense")
			}
		case l.matchTag("@apiContent"):
			d.Content = string(l.data[l.pos:])
		case l.match("@api"):
			l.backup()
			return l.syntaxError(locale.ErrUnknownTag, l.readWord())
		default:
			if l.atEOF() {
				break LOOP
			}
			l.pos++ // 去掉无用的字符。
		}
	}

	return nil
}

// 解析 @api 及其子标签
func (l *lexer) scanAPI(d *doc.Doc) (err *app.SyntaxError) {
	api := &doc.API{}
	t := l.readTag()
	api.Method = t.readWord()
	api.URL = t.readWord()
	api.Summary = t.readLine()

	if len(api.Method) == 0 || len(api.URL) == 0 || len(api.Summary) == 0 {
		return t.syntaxError(locale.ErrTagArgNotEnough, "@api")
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
			api.Error, err = l.scanResponse("@apiError")
		case l.matchTag("@apiSuccess"):
			api.Success, err = l.scanResponse("@apiSuccess")
		case l.match("@api"):
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

	if ignore {
		return nil
	}
	if api.Success == nil {
		return &app.SyntaxError{Message: locale.ErrSuccessNotEmpty}
	}

	if len(api.Group) == 0 {
		api.Group = app.DefaultGroupName
	}

	d.NewAPI(api)
	return nil
}

func (l *lexer) scanGroup(api *doc.API) *app.SyntaxError {
	t := l.readTag()

	api.Group = t.readWord()
	if len(api.Group) == 0 {
		return t.syntaxError(locale.ErrTagArgNotEnough, "@apiGroup")
	}

	if !t.atEOF() {
		t.syntaxError(locale.ErrTagArgTooMuch, "@apiGroup")
	}

	return nil
}

func (l *lexer) scanAPIQueries(api *doc.API) *app.SyntaxError {
	if api.Queries == nil {
		api.Queries = make([]*doc.Param, 0, 1)
	}

	p, err := l.scanAPIParam("@apiQuery")
	if err != nil {
		return err
	}
	api.Queries = append(api.Queries, p)
	return nil
}

func (l *lexer) scanAPIParams(api *doc.API) *app.SyntaxError {
	if api.Params == nil {
		api.Params = make([]*doc.Param, 0, 1)
	}

	p, err := l.scanAPIParam("@apiParam")
	if err != nil {
		return err
	}
	api.Params = append(api.Params, p)
	return nil
}

// 解析 @apiRequest 及其子标签
func (l *lexer) scanAPIRequest(api *doc.API) *app.SyntaxError {
	t := l.readTag()
	r := &doc.Request{
		Type:     t.readLine(),
		Headers:  map[string]string{},
		Params:   []*doc.Param{},
		Examples: []*doc.Example{},
	}
	if !t.atEOF() {
		return t.syntaxError(locale.ErrTagArgTooMuch, "@apiRequest")
	}

LOOP:
	for {
		switch {
		case l.matchTag("@apiHeader"):
			t := l.readTag()
			key := t.readWord()
			val := t.readLine()
			if len(key) == 0 || len(val) == 0 {
				return t.syntaxError(locale.ErrTagArgNotEnough, "@apiHeader")
			}
			if !t.atEOF() {
				return t.syntaxError(locale.ErrTagArgTooMuch, "@apiHeader")
			}
			r.Headers[string(key)] = string(val)
		case l.matchTag("@apiParam"):
			p, err := l.scanAPIParam("@apiParam")
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

// 解析 @apiSuccess 或是 @apiError 及其子标签。
func (l *lexer) scanResponse(tagName string) (*doc.Response, *app.SyntaxError) {
	tag := l.readTag()
	resp := &doc.Response{
		Code:     tag.readWord(),
		Summary:  tag.readLine(),
		Headers:  map[string]string{},
		Params:   []*doc.Param{},
		Examples: []*doc.Example{},
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
		case l.matchTag("@apiHeader"):
			t := l.readTag()
			key := t.readWord()
			val := t.readLine()
			if len(key) == 0 || len(val) == 0 {
				return nil, t.syntaxError(locale.ErrTagArgNotEnough, "@apiHeader")
			}
			if !t.atEOF() {
				return nil, t.syntaxError(locale.ErrTagArgTooMuch, "@apiHeader")
			}
			resp.Headers[key] = val
		case l.matchTag("@apiParam"):
			p, err := l.scanAPIParam("@apiParam")
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

// 解析 @apiExample 标签
func (l *lexer) scanAPIExample() (*doc.Example, *app.SyntaxError) {
	tag := l.readTag()
	example := &doc.Example{
		Type: tag.readWord(),
		Code: tag.readEnd(),
	}

	if len(example.Type) == 0 || len(example.Code) == 0 {
		return nil, tag.syntaxError(locale.ErrTagArgNotEnough, "@apiExample")
	}

	return example, nil
}

// 解析 @apiParam 标签
func (l *lexer) scanAPIParam(tagName string) (*doc.Param, *app.SyntaxError) {
	p := &doc.Param{}

	tag := l.readTag()
	p.Name = tag.readWord()
	p.Type = tag.readWord()
	p.Summary = tag.readEnd()
	if len(p.Name) == 0 || len(p.Type) == 0 || len(p.Summary) == 0 {
		return nil, tag.syntaxError(locale.ErrTagArgNotEnough, tagName)
	}
	return p, nil
}
