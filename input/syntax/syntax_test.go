// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package syntax

import (
	"log"
	"testing"

	"github.com/caixw/apidoc/types"
	"github.com/caixw/apidoc/vars"
	"github.com/issue9/assert"
)

// 生成一个 Input 实例
func newInput(data []rune, logs ...*log.Logger) *Input {
	i := &Input{
		Data: data,
	}

	if len(logs) > 0 {
		i.Error = logs[0]
	}

	if len(logs) > 1 {
		i.Warn = logs[1]
	}

	return i
}

func TestScanAPIExample(t *testing.T) {
	a := assert.New(t)

	// 正常测试
	code := ` xml
<root>
    <data>123</data>
</root>`
	matchCode := `<root>
    <data>123</data>
</root>`
	l := newLexerString(code)
	a.NotNil(l)
	e, err := l.scanAPIExample()
	a.NotError(err).
		Equal(e.Type, "xml").
		Equal(e.Code, matchCode)

	code = ` xml <root>
    <data>123</data>
</root>
 @apiURL abc/test`
	matchCode = `<root>
    <data>123</data>
</root>`
	l = newLexerString(code)
	a.NotNil(l)
	e, err = l.scanAPIExample()
	a.NotError(err).
		Equal(e.Type, "xml").
		Equal(len(e.Code), len(matchCode)).
		Equal(e.Code, matchCode)
}

func TestScanAPIParam(t *testing.T) {
	a := assert.New(t)

	// 正常语法测试
	l := newLexerString("id int optional 用户 id号\n")
	p, ok := l.scanAPIParam(vars.APIQuery)
	a.True(ok).NotNil(p)
	a.Equal(p.Name, "id").
		Equal(p.Type, "int").
		Equal(p.Summary, "optional 用户 id号")

	// 缺少参数
	l = newLexerString("id int \n")
	p, ok = l.scanAPIParam(vars.APIQuery)
	a.False(ok).Nil(p)

	// 缺少参数
	l = newLexerString("id  \n")
	p, ok = l.scanAPIParam(vars.APIQuery)
	a.False(ok).Nil(p)
}

func TestScanAPI(t *testing.T) {
	a := assert.New(t)

	// 正常情况
	l := newLexerString(" get test.com/api.json?k=1 summary summary\n api description\n@apiSuccess 200 OK")

	api, ok := l.scanAPI()
	a.True(ok).NotNil(api)
	a.Equal(api.Method, "get").
		Equal(api.URL, "test.com/api.json?k=1").
		Equal(api.Summary, "summary summary").
		Equal(api.Description, "api description")

	// 多行description
	l = newLexerString(" post test.com/api.json?K=1  summary summary\n api \ndescription\n@apiSuccess 400 summary")
	api, ok = l.scanAPI()
	a.True(ok)
	a.Equal(api.URL, "test.com/api.json?K=1").
		Equal(api.Method, "post").
		Equal(api.Summary, "summary summary").
		Equal(api.Description, "api \ndescription")

	// 缺少description参数
	l = newLexerString("get test.com/api.json summary summary \n@apiSuccess 400 error")
	api, ok = l.scanAPI()
	a.True(ok).NotNil(api)
	a.Equal(api.Method, "get").
		Equal(api.URL, "test.com/api.json").
		Equal(api.Summary, "summary summary").
		Equal(api.Description, "")

	// 缺少description参数
	l = newLexerString("get test.com/api.json summary summary\n \n@apiSuccess 400 error")
	api, ok = l.scanAPI()
	a.True(ok).NotNil(api)
	a.Equal(api.Method, "get").
		Equal(api.URL, "test.com/api.json").
		Equal(api.Summary, "summary summary").
		Equal(api.Description, "")

	// 没有任何参数
	l = newLexerString("  ")
	api, ok = l.scanAPI()
	a.False(ok).Nil(api)
}

func TestScanAPIDoc(t *testing.T) {
	a := assert.New(t)

	code := ` title of apidoc
@apiVersion 2.0.1
@apiBaseURL https://api.caixw.io
@apiLicense MIT https://opensource.org/licenses/MIT
@apiContent
line1
line2`

	l := newLexerString(code)
	a.NotNil(l)
	d := &types.Doc{}

	a.NotError(l.scanAPIDoc(d))
	a.Equal(d.Version, "2.0.1").
		Equal(d.Title, "title of apidoc").
		Equal(d.BaseURL, "https://api.caixw.io").
		Equal(d.LicenseName, "MIT").
		Equal(d.LicenseURL, "https://opensource.org/licenses/MIT").
		Equal(d.Content, "\nline1\nline2")

		// 重复内容，报错
	code = `title of apidoc2
		@apiBaseURL http://api.caixw.io
line1
line2`
	l = newLexerString(code)
	a.False(l.scanAPIDoc(d))

	// 检测各个都是空值的情况
	code = `2.9 title of apidoc
`
	l = newLexerString(code)
	a.NotNil(l)
	d = &types.Doc{}
	a.NotError(l.scanAPIDoc(d))
	a.Equal(d.Title, "2.9 title of apidoc").
		Equal(d.Version, "").
		Equal(d.BaseURL, "").
		Equal(d.LicenseName, "").
		Equal(d.LicenseURL, "")
}

func TestScanAPIRequest(t *testing.T) {
	a := assert.New(t)

	code := ` xml
 @apiHeader h1 v1
 @apiHeader h2 v2
@apiParam p1 int optional p1 summary
@apiParam p2 int p2 summary
@apiExample json
{
    p1:v1,
    p2:v2
}
@apiExample xml
<root>
    <p1>v1</p1>
    <p2>v2</p2>
</root>
`
	l := newLexerString(code)
	r, ok := l.scanAPIRequest()
	a.True(ok).NotNil(r)
	a.Equal(2, len(r.Headers)).
		Equal(r.Headers["h1"], "v1").
		Equal(r.Headers["h2"], "v2").
		Equal(r.Params[0].Name, "p1").
		Equal(r.Params[1].Summary, "p2 summary").
		Equal(r.Examples[0].Type, "json").
		Equal(r.Examples[1].Type, "xml").
		Equal(r.Examples[1].Code, `<root>
    <p1>v1</p1>
    <p2>v2</p2>
</root>`)

	code = ` xml
@apiHeader h1 v1
@apiParam p1 int p1 summary

@apiExample xml
<root>
    <p1>v1</p1>
</root>
@apiGroup abc
`
	matchCode := `<root>
    <p1>v1</p1>
</root>`
	l = newLexerString(code)
	r, ok = l.scanAPIRequest()
	a.True(ok).NotNil(r)
	a.Equal(1, len(r.Headers)).
		Equal(r.Headers["h1"], "v1").
		Equal(r.Params[0].Name, "p1").
		Equal(r.Examples[0].Type, "xml").
		Equal(r.Examples[0].Code, matchCode)
}

func TestScanResponse(t *testing.T) {
	a := assert.New(t)

	code := ` 200 json
@apiHeader h1 v1
@apiHeader h2 v2
@apiParam p1 int optional p1 summary
@apiParam p2 int p2 summary
@apiExample json
{
    p1:v1,
    p2:v2
}
@apiExample xml
<root>
    <p1>v1</p1>
    <p2>v2</p2>
</root>
`
	l := newLexerString(code)
	resp, ok := l.scanResponse(vars.APIError)
	a.True(ok).NotNil(resp)
	a.Equal(resp.Code, "200").
		Equal(resp.Summary, "json").
		Equal(resp.Headers["h1"], "v1").
		Equal(resp.Headers["h2"], "v2").
		Equal(resp.Params[0].Name, "p1").
		Equal(resp.Params[1].Summary, "p2 summary").
		Equal(resp.Examples[0].Type, "json").
		Equal(resp.Examples[1].Type, "xml")

	code = ` 200 xml  status summary
 @apiHeader h1 v1
 @apiParam p1 int p1 summary

@apiExample xml
<root>
    <p1>v1</p1>
</root>
@apiError
`
	matchCode := `<root>
    <p1>v1</p1>
</root>`
	l = newLexerString(code)
	resp, ok = l.scanResponse(vars.APIError)
	a.True(ok).NotNil(resp)
	a.Equal(resp.Code, "200").
		Equal(resp.Summary, "xml  status summary").
		Equal(resp.Headers["h1"], "v1").
		Equal(resp.Params[0].Name, "p1").
		Equal(resp.Examples[0].Code, matchCode)

	// 缺少必要的参数
	code = `
@apiSuccess g
`
	l = newLexerString(code)
	resp, ok = l.scanResponse(vars.APIError)
	a.False(ok).Nil(resp)
}

func TestParse(t *testing.T) {
	a := assert.New(t)
	doc := types.NewDoc()

	code := `
@api get /baseurl/api/login api summary
api description 1
api description 2
@apiQuery q1 int q1 summary
@apiQuery q2 int q2 summary
@apiParam p1 int p1 summary
@apiParam p2 int p2 summary

@apiSuccess 200 json
@apiHeader h1 v1
@apiHeader h2 v2
@apiParam p1 int optional p1 summary
@apiParam p2 int p2 summary
@apiExample json
{
    p1:v1,
    p2:v2
}
@apiExample xml
<root>
    <p1>v1</p1>
    <p2>v2</p2>
</root>

@apiError 200 json
@apiHeader h1 v1
@apiHeader h2 v2
`
	Parse(&Input{Data: []rune(code)}, doc)
	d := doc.Apis[0]

	a.Equal(d.URL, "/baseurl/api/login").
		Equal(d.Group, vars.DefaultGroupName).
		Equal(d.Summary, "api summary").
		Equal(d.Description, "api description 1\napi description 2")

	a.Equal(2, len(d.Queries)).
		Equal(2, len(d.Params))

	q := d.Queries
	a.Equal(q[0].Name, "q1").Equal(q[0].Summary, "q1 summary")

	p := d.Params
	a.Equal(p[0].Name, "p1").Equal(p[0].Summary, "p1 summary")

	s := d.Success
	a.Equal(s.Code, "200").
		Equal(s.Summary, "json").
		Equal(s.Headers["h1"], "v1").
		Equal(s.Headers["h2"], "v2").
		Equal(s.Params[0].Name, "p1").
		Equal(s.Params[1].Summary, "p2 summary").
		Equal(s.Examples[0].Type, "json").
		Equal(s.Examples[1].Type, "xml")

	s = d.Error
	a.Equal(s.Code, "200")

	// 不包含api定义的代码块，将返回 nil 且不会有内容添加到 doc1.Apis
	code = `
Copyright 2015 by caixw, All rights reserved.
Use of this source code is governed by a MIT
license that can be found in the LICENSE file.
`
	l := len(doc.Apis)
	Parse(&Input{Data: []rune(code)}, doc)
	a.Equal(l, len(doc.Apis)) // 没有增加新内容

	// @apiGroup
	code = `
@api delete /admin/users/{id} delete users
@apiGroup
@apiParam id int user id
@apiError 400 error
`
	l = len(doc.Apis)
	Parse(&Input{Data: []rune(code)}, doc)
	a.Equal(l, len(doc.Apis)) // 不会增加新的内容

	// @apiIgnore
	code = `
@api delete /admin/users/{id} delete users
@apiIgnore
@apiParam id int user id
@apiError 400 error
`
	l = len(doc.Apis)
	Parse(&Input{Data: []rune(code)}, doc)
	a.Equal(l, len(doc.Apis))

	// @apiUnknownTag 非正常标签
	code = `
@apiUnknownTag xxx
@api delete /admin/users/{id} delete users
@apiGroup group
ab @apiIgnore
@apiParam id int user id
@apiSuccess 200 OK
`
	l = len(doc.Apis)
	Parse(&Input{Data: []rune(code)}, doc)
	a.Equal(l, len(doc.Apis))

	// 不认识的标签
	code = `
@apiUnknownTag1 不认识的标签

@api delete /admin/users/{id} delete users
@apiGroup group
@apiParam id int user id
@apiSuccess 200 OK
@apiUnknownTag2 不认识的标签
`
	l = len(doc.Apis)
	Parse(&Input{Data: []rune(code)}, doc)
	a.Equal(l+1, len(doc.Apis))

	// @apiIgno 不认识的标签，会被过滤
	code = `
@api delete /admin/users/{id} delete users
@apiIgno
@apiGroup group
@apiParam id int user id
@apiSuccess 200 OK
`
	l = len(doc.Apis)
	Parse(&Input{Data: []rune(code)}, doc)
	a.Equal(l+1, len(doc.Apis))
}
