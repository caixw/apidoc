// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"testing"

	"github.com/issue9/assert"
)

var synerr = &SyntaxError{}

func TestscanAPIExample(t *testing.T) {
	a := assert.New(t)

	// 正常测试
	code := ` xml
<root>
    <data>123</data>
</root>`
	matchCode := `<root>
    <data>123</data>
</root>`
	l := newLexer([]rune(code))
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
	l = newLexer([]rune(code))
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
	l := newLexer([]rune("id int optional 用户 id号\n"))
	p, err := l.scanAPIParam()
	a.NotError(err).NotNil(p)
	a.Equal(p.Name, "id").
		Equal(p.Type, "int").
		Equal(p.Summary, "optional 用户 id号")

	// 缺少参数
	l = newLexer([]rune("id int \n"))
	p, err = l.scanAPIParam()
	a.ErrorType(err, synerr).Nil(p)

	// 缺少参数
	l = newLexer([]rune("id  \n"))
	p, err = l.scanAPIParam()
	a.ErrorType(err, synerr).Nil(p)
}

func TestScanAPI(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	// 正常情况
	l := newLexer([]rune(" get test.com/api.json?k=1 summary summary\n api description"))
	a.NotError(l.scanAPI(api))
	a.Equal(api.Method, "get").
		Equal(api.URL, "test.com/api.json?k=1").
		Equal(api.Summary, "summary summary").
		Equal(api.Description, "api description")

	// 多行description
	l = newLexer([]rune(" post test.com/api.json?K=1  summary summary\n api \ndescription\n@api summary"))
	a.NotError(l.scanAPI(api))
	a.Equal(api.URL, "test.com/api.json?K=1").
		Equal(api.Method, "post").
		Equal(api.Summary, "summary summary").
		Equal(api.Description, "api \ndescription")

	// 缺少description参数
	l = newLexer([]rune("get test.com/api.json summary summary"))
	a.NotError(l.scanAPI(api))
	a.Equal(api.Method, "get").
		Equal(api.URL, "test.com/api.json").
		Equal(api.Summary, "summary summary").
		Equal(api.Description, "")

	// 缺少description参数
	l = newLexer([]rune("get test.com/api.json summary summary\n@apiURL"))
	a.NotError(l.scanAPI(api))
	a.Equal(api.Method, "get").
		Equal(api.URL, "test.com/api.json").
		Equal(api.Summary, "summary summary").
		Equal(api.Description, "")

	// 没有任何参数
	l = newLexer([]rune("  "))
	a.ErrorType(l.scanAPI(api), synerr)
}

func TestScanAPIRequest(t *testing.T) {
	a := assert.New(t)
	api := &API{}

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
	l := newLexer([]rune(code))
	a.NotError(l.scanAPIRequest(api))
	a.NotNil(api.Request)
	r := api.Request
	a.Equal(2, len(r.Headers)).
		Equal(r.Headers["h1"], "v1").
		Equal(r.Headers["h2"], "v2").
		Equal(r.Params[0].Name, "p1").
		Equal(r.Params[1].Summary, "p2 summary").
		Equal(r.Examples[0].Type, "json").
		Equal(r.Examples[1].Type, "xml")

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
	l = newLexer([]rune(code))
	a.NotError(l.scanAPIRequest(api))
	a.NotNil(api.Request)
	r = api.Request
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
	l := newLexer([]rune(code))
	resp, err := l.scanResponse()
	a.NotError(err).NotNil(resp)
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
	l = newLexer([]rune(code))
	resp, err = l.scanResponse()
	a.NotError(err).NotNil(resp)
	a.Equal(resp.Code, "200").
		Equal(resp.Summary, "xml  status summary").
		Equal(resp.Headers["h1"], "v1").
		Equal(resp.Params[0].Name, "p1").
		Equal(resp.Examples[0].Code, matchCode)

	// 缺少必要的参数
	code = ` 
@apiSuccess g
`
	l = newLexer([]rune(code))
	resp, err = l.scanResponse()
	a.Error(err).Nil(resp)
}

func TestDoc_Scan(t *testing.T) {
	a := assert.New(t)
	doc1 := New()

	code := `
@api get /baseurl/api/login api summary
api description 1
api description 2
@apiGroup users
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
	err := doc1.Scan(code)
	a.NotError(err)
	d := doc1.Apis[0]

	a.Equal(d.URL, "/baseurl/api/login").
		Equal(d.Group, "users").
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

	// 不包含api定义的代码块，将返回一个error,nil
	code = `
Copyright 2015 by caixw, All rights reserved.
Use of this source code is governed by a MIT
license that can be found in the LICENSE file.
`
	err = doc1.Scan(code)
	a.NotError(err)
}
