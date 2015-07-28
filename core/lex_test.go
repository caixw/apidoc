// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLexer_lineNumber(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]byte("\n\n"), 100, "file.go")
	a.NotNil(l)

	a.Equal(100, l.lineNumber())

	l.next()
	a.Equal(101, l.lineNumber())

	l.next()
	a.Equal(102, l.lineNumber())

	l.backup()
	a.Equal(101, l.lineNumber())
}

func TestLexer_next(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]byte("ab\ncd\n"), 100, "file.go")
	a.NotNil(l)

	a.Equal('a', l.next())
	a.Equal('b', l.next())
	a.Equal('\n', l.next())
	a.Equal('c', l.next())

	// 退回一个字符
	l.backup()
	a.Equal('c', l.next())

	// 退回多个字符
	l.backup()
	l.backup()
	l.backup()
	a.Equal('c', l.next())

	a.Equal('d', l.next())
	a.Equal('\n', l.next())
	a.Equal(eof, l.next()) // 文件结束
	a.Equal(eof, l.next())
}

func TestLexer_nextLine(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]byte("line1\n line2 \n"), 100, "file.go")
	a.NotNil(l)

	a.Equal("line1", l.nextLine())
	l.backup()
	l.backup()
	a.Equal('l', l.next())

	a.Equal("ine1", l.nextLine())
	a.Equal(" line2 ", l.nextLine()) // 空格会被过滤
	l.backup()
	a.Equal(" line2 ", l.nextLine())

	a.Equal("", l.nextLine()) // 没有更多内容了
}

func TestLexer_nextWord(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]byte("word1 word2\nword3"), 100, "file.go")
	a.NotNil(l)

	l.next()               // w
	l.next()               // o
	w, eol := l.nextWord() // rd1
	a.False(eol).Equal(w, "rd1")
	l.backup() // 多次调用backup，只启作用一次。
	l.backup()
	a.Equal(l.next(), 'r')

	w, eol = l.nextWord()
	a.False(eol).Equal(w, "d1")

	// 对空格的操作，不返回任何值。
	w, eol = l.nextWord()
	w, eol = l.nextWord()
	w, eol = l.nextWord()
	a.False(eol).Equal(w, "")

	// 第2个单词
	l.skipSpace()
	w, eol = l.nextWord() // word2
	a.True(eol).Equal(w, "word2")

	// eol，不会查找下一行的内容
	w, eol = l.nextWord()
	a.True(eol).Equal(w, "")
	// eol，不会查找下一行的内容
	w, eol = l.nextWord()
	a.True(eol).Equal(w, "")

	// 跳到下一行,eol
	l.next()
	w, eol = l.nextWord()
	a.True(eol).Equal(w, "word3")

	// eol,没有再多的内容了
	w, eol = l.nextWord()
	a.True(eol).Equal(w, "")
}

func TestLexer_match(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]byte("line1\n line2 \n"), 100, "file.go")
	a.NotNil(l)

	a.True(l.match("line"))
	a.Equal('1', l.next())

	l.next() // \n
	l.next() // 空格
	l.next() // l

	a.False(l.match("2222")) // 不匹配，不会移动位置
	a.True(l.match("ine2"))  // 正确匹配
	l.backup()
	l.backup()
	a.Equal('i', l.next())

	// 超过剩余字符的长度。
	a.False(l.match("ne2\n\n"))
}

func TestLexer_skipSpace(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]byte("  ln  1  \n 2 \n"), 100, "file.go")
	a.NotNil(l)

	l.skipSpace() // 跳转起始的2个空格
	l.skipSpace() // 不会跳过ln字符
	l.skipSpace() // 不会跳过ln字符
	a.Equal('l', l.next())
	l.next() // n

	l.skipSpace()
	l.backup() // lexer.backup对lexer.skipSpace()不启作用
	a.Equal('1', l.next())

	l.skipSpace() // 不能跳过\n
	l.skipSpace() // 不能跳过\n
	a.Equal('\n', l.next())

	l.skipSpace()
	a.Equal('2', l.next())

	l.skipSpace()
	a.Equal('\n', l.next())
	l.next()
	l.next()
	a.Equal(eof, l.next())

	// 文件结尾
	l.skipSpace()
	a.Equal(eof, l.next())
}

func TestLexer_scanApiURL(t *testing.T) {
	a := assert.New(t)
	d := &doc{}

	// 正常情况
	l := newLexer([]byte("  api/login"), 100, "file.go")
	a.NotError(l.scanApiURL(d))
	a.Equal(d.URL, "api/login")

	// 缺少参数
	l = newLexer([]byte(" "), 100, "file.go")
	a.Error(l.scanApiURL(d))

	// 多个参数
	l = newLexer([]byte("  api/login abctest/adf"), 100, "file.go")
	a.NotError(l.scanApiURL(d))
	a.Equal(d.URL, "api/login")
}

func TestLexer_scanApiMethods(t *testing.T) {
	a := assert.New(t)
	d := &doc{}

	// 正常情况
	l := newLexer([]byte("  get"), 100, "file.go")
	a.NotError(l.scanApiMethods(d))
	a.Equal(d.Methods, "get")

	// 缺少参数
	l = newLexer([]byte(" "), 100, "file.go")
	a.Error(l.scanApiMethods(d))

	// 多个参数
	l = newLexer([]byte("  get post"), 100, "file.go")
	a.NotError(l.scanApiMethods(d))
	a.Equal(d.Methods, "get post")

	// 多个参数
	l = newLexer([]byte("  get post\n@api"), 100, "file.go")
	a.NotError(l.scanApiMethods(d))
	a.Equal(d.Methods, "get post")
}

func TestLexer_scanApiVersion(t *testing.T) {
	a := assert.New(t)
	d := &doc{}

	// 正常情况
	l := newLexer([]byte("  0.1.1"), 100, "file.go")
	a.NotError(l.scanApiVersion(d))
	a.Equal(d.Version, "0.1.1")

	// 缺少参数
	l = newLexer([]byte(" "), 100, "file.go")
	a.Error(l.scanApiVersion(d))

	// 多个参数
	l = newLexer([]byte("  0.1.1  abcd"), 100, "file.go")
	a.NotError(l.scanApiVersion(d))
	a.Equal(d.Version, "0.1.1")
}

func TestLexer_scanApiGroup(t *testing.T) {
	a := assert.New(t)
	d := &doc{}

	// 正常情况
	l := newLexer([]byte("  g1"), 100, "file.go")
	a.NotError(l.scanApiGroup(d))
	a.Equal(d.Group, "g1")

	// 缺少参数
	l = newLexer([]byte(" "), 100, "file.go")
	a.Error(l.scanApiGroup(d))

	// 多个参数
	l = newLexer([]byte("  g1  abcd"), 100, "file.go")
	a.NotError(l.scanApiGroup(d))
	a.Equal(d.Group, "g1")
}

func TestLexer_scanApiQuery(t *testing.T) {
	a := assert.New(t)
	d := &doc{Queries: []*param{}}

	// 正常情况
	l := newLexer([]byte("id int user id"), 100, "file.go")
	a.NotError(l.scanApiQuery(d))
	q0 := d.Queries[0]
	a.Equal(q0.Name, "id").
		Equal(q0.Type, "int").
		False(q0.Optional).
		Equal(q0.Description, "user id")

	// 再添加一个参数
	l = newLexer([]byte("name string user name"), 100, "file.go")
	a.NotError(l.scanApiQuery(d))
	q1 := d.Queries[1]
	a.Equal(q1.Name, "name").
		Equal(q1.Type, "string").
		False(q1.Optional).
		Equal(q1.Description, "user name")
}

func TestLexer_scanApiExample(t *testing.T) {
	a := assert.New(t)

	// 正常测试
	code := ` xml
<root>
    <data>123</data>
</root>`
	matchCode := `
<root>
    <data>123</data>
</root>`
	l := newLexer([]byte(code), 100, "file.go")
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
	l = newLexer([]byte(code), 100, "file.go")
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
	l := newLexer([]byte("id int optional 用户 id号\n"), 100, "file.go")
	p, err := l.scanApiParam()
	a.NotError(err).NotNil(p)
	a.Equal(p.Name, "id").
		Equal(p.Type, "int").
		True(p.Optional).
		Equal(p.Description, "用户 id号")

	// 大小写混合的optional
	l = newLexer([]byte("id int OptionAl 用户 id号\n"), 100, "file.go")
	p, err = l.scanApiParam()
	a.NotError(err).NotNil(p)
	a.Equal(p.Name, "id").
		Equal(p.Type, "int").
		True(p.Optional).
		Equal(p.Description, "用户 id号")

	// 缺少optional参数
	l = newLexer([]byte("id int optional1 用户 id号\n"), 100, "file.go")
	p, err = l.scanApiParam()
	a.NotError(err).NotNil(p)
	a.Equal(p.Name, "id").
		Equal(p.Type, "int").
		False(p.Optional).
		Equal(p.Description, "optional1 用户 id号")

	// 缺少参数
	l = newLexer([]byte("id int \n"), 100, "file.go")
	p, err = l.scanApiParam()
	a.Error(err).Nil(p)

	// 缺少参数
	l = newLexer([]byte("id  \n"), 100, "file.go")
	p, err = l.scanApiParam()
	a.Error(err).Nil(p)
}

func TestLexer_scanApi(t *testing.T) {
	a := assert.New(t)
	d := &doc{}

	// 正常情况
	l := newLexer([]byte("  summary summary\n api description"), 100, "file.go")
	a.NotError(l.scanApi(d))
	a.Equal(d.Summary, "summary summary").
		Equal(d.Description, " api description")

	// 多行description
	l = newLexer([]byte("  summary summary\n api \ndescription\n@api summary"), 100, "file.go")
	a.NotError(l.scanApi(d))
	a.Equal(d.Summary, "summary summary").
		Equal(d.Description, " api \ndescription\n")

	// 缺少description参数
	l = newLexer([]byte("summary summary"), 100, "file.go")
	a.NotError(l.scanApi(d))
	a.Equal(d.Summary, "summary summary").
		Equal(d.Description, "")

	// 缺少description参数
	l = newLexer([]byte("summary summary\n@apiURL"), 100, "file.go")
	a.NotError(l.scanApi(d))
	a.Equal(d.Summary, "summary summary").
		Equal(d.Description, "")

	// 没有任何参数
	l = newLexer([]byte("  "), 100, "file.go")
	a.Error(l.scanApi(d))
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
	l := newLexer([]byte(code), 100, "file.go")
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
@apiStatus
`
	matchCode := `
<root>
    <p1>v1</p1>
</root>`
	l = newLexer([]byte(code), 100, "file.go")
	a.NotError(l.scanApiRequest(d))
	a.NotNil(d.Request)
	r = d.Request
	a.Equal(1, len(r.Headers)).
		Equal(r.Headers["h1"], "v1").
		Equal(r.Params[0].Name, "p1").
		Equal(r.Examples[0].Type, "xml").
		Equal(r.Examples[0].Code, matchCode)
}

func TestLexer_scanApiStatus(t *testing.T) {
	a := assert.New(t)
	d := &doc{}

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
	l := newLexer([]byte(code), 100, "file.go")
	a.NotError(l.scanApiStatus(d))
	a.Equal(1, len(d.Status))
	s := d.Status[0]
	a.Equal(s.Code, "200").
		Equal(s.Type, "json").
		Equal(s.Summary, "").
		Equal(s.Headers["h1"], "v1").
		Equal(s.Headers["h2"], "v2").
		Equal(s.Params[0].Name, "p1").
		Equal(s.Params[1].Description, "p2 summary").
		Equal(s.Examples[0].Type, "json").
		Equal(s.Examples[1].Type, "xml")

	code = ` 200 xml  status summary
@apiHeader h1 v1
@apiParam p1 int p1 summary

@apiExample xml
<root>
    <p1>v1</p1>
</root>
@apiStatus
`
	matchCode := `
<root>
    <p1>v1</p1>
</root>`
	l = newLexer([]byte(code), 100, "file.go")
	a.NotError(l.scanApiStatus(d))
	a.Equal(2, len(d.Status))
	s = d.Status[1]
	a.Equal(s.Code, "200").
		Equal(s.Type, "xml").
		Equal(s.Summary, "status summary").
		Equal(s.Headers["h1"], "v1").
		Equal(s.Params[0].Name, "p1").
		Equal(s.Examples[0].Code, matchCode)

	// 缺少必要的参数
	code = ` 200
@apiStatus
`
	l = newLexer([]byte(code), 100, "file.go")
	a.Error(l.scanApiStatus(d))
}

func TestLexer_scan(t *testing.T) {
	a := assert.New(t)

	code := `
@api api summary
api description 1
api description 2
@apiURL /baseurl/api/login
@apiMethods get/post
@apiVersion 1.0
@apiGroup users
@apiQuery q1 int q1 summary
@apiQuery q2 int q2 summary
@apiParam p1 int p1 summary
@apiParam p2 int p2 summary
@apiStatus 200 json
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
	l := newLexer([]byte(code), 100, "file.go")
	d, err := l.scan()
	a.NotError(err).NotNil(d)

	a.Equal(d.Version, "1.0").
		Equal(d.URL, "/baseurl/api/login").
		Equal(d.Group, "users").
		Equal(d.Summary, "api summary").
		Equal(d.Description, "api description 1\napi description 2\n")

	a.Equal(2, len(d.Queries)).Equal(2, len(d.Params)).Equal(1, len(d.Status))

	q := d.Queries
	a.Equal(q[0].Name, "q1").Equal(q[0].Description, "q1 summary")

	p := d.Params
	a.Equal(p[0].Name, "p1").Equal(p[0].Description, "p1 summary")

	s := d.Status[0]
	a.Equal(s.Code, "200").
		Equal(s.Type, "json").
		Equal(s.Summary, "").
		Equal(s.Headers["h1"], "v1").
		Equal(s.Headers["h2"], "v2").
		Equal(s.Params[0].Name, "p1").
		Equal(s.Params[1].Description, "p2 summary").
		Equal(s.Examples[0].Type, "json").
		Equal(s.Examples[1].Type, "xml")
}
