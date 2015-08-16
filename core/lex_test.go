// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	"github.com/issue9/assert"
)

var synerr = &SyntaxError{}

func TestLexer_lineNumber(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]rune("\n\n"), 100, "file.go")
	a.NotNil(l)

	a.Equal(100, l.lineNumber())

	l.pos++
	a.Equal(101, l.lineNumber())

	l.pos++
	a.Equal(102, l.lineNumber())
}

func TestLexer_read(t *testing.T) {
	a := assert.New(t)

	l := newLexer([]rune("line1\n @api line2 \n"), 100, "file.go")
	a.NotNil(l)

	word := l.read("@api")
	a.Equal(word, "line1")

	l.match("@api")
	word = l.read("@api")
	a.Equal(word, "line2")
}

func TestLexer_readN(t *testing.T) {
	a := assert.New(t)

	l := newLexer([]rune("line1\n @api line2 \n"), 100, "file.go")
	a.NotNil(l)

	words, err := l.readN(1, "@api")
	a.NotError(err).Equal(words, []string{"line1"})

	l.match("@api")
	words, err = l.readN(1, "@api") // 行尾并没有@api,匹配eof
	a.NotError(err).Equal(words, []string{"line2"})

	// 多词匹配
	l = newLexer([]rune("word1 word2 word3 word4\n @api word5 word6 \n"), 100, "file.go")
	words, err = l.readN(2, "\n")
	a.NotError(err).Equal(words, []string{"word1", "word2 word3 word4"})

	l.match("@api")
	words, err = l.readN(5, "\n")
	a.Error(err)

	l = newLexer([]rune("word1 word2 word3 word4\n"), 100, "file.go")
	words, err = l.readN(1, "\n")
	a.NotError(err).Equal(words, []string{"word1 word2 word3 word4"})
}

func TestLexer_match(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]rune("line1\n line2 \n"), 100, "file.go")
	a.NotNil(l)

	a.True(l.match("line"))
	a.Equal('1', l.data[l.pos])
	l.pos++

	l.pos++ // \n
	l.pos++ // 空格
	l.pos++ // l

	a.False(l.match("2222")) // 不匹配，不会移动位置
	a.True(l.match("ine2"))  // 正确匹配
	l.backup()
	l.backup()
	a.Equal('i', l.data[l.pos])
	l.pos++

	// 超过剩余字符的长度。
	a.False(l.match("ne2\n\n"))
}

func TestLexer_scanApiGroup(t *testing.T) {
	a := assert.New(t)
	d := &doc{}

	// 正常情况
	l := newLexer([]rune("  g1"), 100, "file.go")
	a.NotError(l.scanApiGroup(d))
	a.Equal(d.Group, "g1")

	// 缺少参数
	l = newLexer([]rune(" "), 100, "file.go")
	a.ErrorType(l.scanApiGroup(d), synerr)

	// 带空格
	l = newLexer([]rune("  g1  abcd"), 100, "file.go")
	a.NotError(l.scanApiGroup(d))
	a.Equal(d.Group, "g1  abcd")
}

func TestLexer_scanApiQuery(t *testing.T) {
	a := assert.New(t)
	d := &doc{Queries: []*param{}}

	// 正常情况
	l := newLexer([]rune("id int user id"), 100, "file.go")
	a.NotError(l.scanApiQuery(d))
	q0 := d.Queries[0]
	a.Equal(q0.Name, "id").
		Equal(q0.Type, "int").
		Equal(q0.Description, "user id")

	// 再添加一个参数
	l = newLexer([]rune("name string user name"), 100, "file.go")
	a.NotError(l.scanApiQuery(d))
	q1 := d.Queries[1]
	a.Equal(q1.Name, "name").
		Equal(q1.Type, "string").
		Equal(q1.Description, "user name")
}

func TestLexer_scanApiExample(t *testing.T) {
	a := assert.New(t)

	// 正常测试
	code := ` xml
<root>
    <data>123</data>
</root>`
	matchCode := `<root>
    <data>123</data>
</root>`
	l := newLexer([]rune(code), 100, "file.go")
	a.NotNil(l)
	e, err := l.scanApiExample()
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
	l = newLexer([]rune(code), 100, "file.go")
	a.NotNil(l)
	e, err = l.scanApiExample()
	a.NotError(err).
		Equal(e.Type, "xml").
		Equal(len(e.Code), len(matchCode)).
		Equal(e.Code, matchCode)
}

func TestLexer_scanApiParam(t *testing.T) {
	a := assert.New(t)

	// 正常语法测试
	l := newLexer([]rune("id int optional 用户 id号\n"), 100, "file.go")
	p, err := l.scanApiParam()
	a.NotError(err).NotNil(p)
	a.Equal(p.Name, "id").
		Equal(p.Type, "int").
		Equal(p.Description, "optional 用户 id号")

	// 缺少参数
	l = newLexer([]rune("id int \n"), 100, "file.go")
	p, err = l.scanApiParam()
	a.ErrorType(err, synerr).Nil(p)

	// 缺少参数
	l = newLexer([]rune("id  \n"), 100, "file.go")
	p, err = l.scanApiParam()
	a.ErrorType(err, synerr).Nil(p)
}

func TestLexer_scanApi(t *testing.T) {
	a := assert.New(t)
	d := &doc{}

	// 正常情况
	l := newLexer([]rune(" get test.com/api.json?k=1 summary summary\n api description"), 100, "file.go")
	a.NotError(l.scanApi(d))
	a.Equal(d.Method, "get").
		Equal(d.URL, "test.com/api.json?k=1").
		Equal(d.Summary, "summary summary").
		Equal(d.Description, "api description")

	// 多行description
	l = newLexer([]rune(" post test.com/api.json?K=1  summary summary\n api \ndescription\n@api summary"), 100, "file.go")
	a.NotError(l.scanApi(d))
	a.Equal(d.URL, "test.com/api.json?K=1").
		Equal(d.Method, "post").
		Equal(d.Summary, "summary summary").
		Equal(d.Description, "api \ndescription")

	// 缺少description参数
	l = newLexer([]rune("get test.com/api.json summary summary"), 100, "file.go")
	a.NotError(l.scanApi(d))
	a.Equal(d.Method, "get").
		Equal(d.URL, "test.com/api.json").
		Equal(d.Summary, "summary summary").
		Equal(d.Description, "")

	// 缺少description参数
	l = newLexer([]rune("get test.com/api.json summary summary\n@apiURL"), 100, "file.go")
	a.NotError(l.scanApi(d))
	a.Equal(d.Method, "get").
		Equal(d.URL, "test.com/api.json").
		Equal(d.Summary, "summary summary").
		Equal(d.Description, "")

	// 没有任何参数
	l = newLexer([]rune("  "), 100, "file.go")
	a.ErrorType(l.scanApi(d), synerr)
}

func TestLexer_scanApiRequest(t *testing.T) {
	a := assert.New(t)
	d := &doc{}

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
	l := newLexer([]rune(code), 100, "file.go")
	a.NotError(l.scanApiRequest(d))
	a.NotNil(d.Request)
	r := d.Request
	a.Equal(2, len(r.Headers)).
		Equal(r.Headers["h1"], "v1").
		Equal(r.Headers["h2"], "v2").
		Equal(r.Params[0].Name, "p1").
		Equal(r.Params[1].Description, "p2 summary").
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
	l = newLexer([]rune(code), 100, "file.go")
	a.NotError(l.scanApiRequest(d))
	a.NotNil(d.Request)
	r = d.Request
	a.Equal(1, len(r.Headers)).
		Equal(r.Headers["h1"], "v1").
		Equal(r.Params[0].Name, "p1").
		Equal(r.Examples[0].Type, "xml").
		Equal(r.Examples[0].Code, matchCode)
}

func TestLexer_scanResponse(t *testing.T) {
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
	l := newLexer([]rune(code), 100, "file.go")
	resp, err := l.scanResponse()
	a.NotError(err).NotNil(resp)
	a.Equal(resp.Code, "200").
		Equal(resp.Summary, "json").
		Equal(resp.Headers["h1"], "v1").
		Equal(resp.Headers["h2"], "v2").
		Equal(resp.Params[0].Name, "p1").
		Equal(resp.Params[1].Description, "p2 summary").
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
	l = newLexer([]rune(code), 100, "file.go")
	resp, err = l.scanResponse()
	a.NotError(err).NotNil(resp)
	a.Equal(resp.Code, "200").
		Equal(resp.Summary, "xml  status summary").
		Equal(resp.Headers["h1"], "v1").
		Equal(resp.Params[0].Name, "p1").
		Equal(resp.Examples[0].Code, matchCode)

	// 缺少必要的参数
	code = ` 
@apiGroup g
`
	l = newLexer([]rune(code), 100, "file.go")
	resp, err = l.scanResponse()
	a.Error(err).Nil(resp)
}

func TestLexer_scan(t *testing.T) {
	a := assert.New(t)

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
	l := newLexer([]rune(code), 100, "file.go")
	d, err := l.scan()
	a.NotError(err).NotNil(d)

	a.Equal(d.URL, "/baseurl/api/login").
		Equal(d.Group, "users").
		Equal(d.Summary, "api summary").
		Equal(d.Description, "api description 1\napi description 2")

	a.Equal(2, len(d.Queries)).
		Equal(2, len(d.Params))

	q := d.Queries
	a.Equal(q[0].Name, "q1").Equal(q[0].Description, "q1 summary")

	p := d.Params
	a.Equal(p[0].Name, "p1").Equal(p[0].Description, "p1 summary")

	s := d.Success
	a.Equal(s.Code, "200").
		Equal(s.Summary, "json").
		Equal(s.Headers["h1"], "v1").
		Equal(s.Headers["h2"], "v2").
		Equal(s.Params[0].Name, "p1").
		Equal(s.Params[1].Description, "p2 summary").
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
	l = newLexer([]rune(code), 100, "file.go")
	d, err = l.scan()
	a.Error(err).Nil(d)
}
