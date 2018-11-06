// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"testing"

	"github.com/issue9/assert"
)

func TestAPI_getPathParams(t *testing.T) {
	a := assert.New(t)

	api := &API{Path: "/path"}
	a.NotError(api.getPathParams())

	api = &API{Path: "/path/{id}"}
	names, err := api.getPathParams()
	a.NotError(err).
		Equal(names, map[string]bool{"id": true})

	api = &API{Path: "/path/{id}/{name}"}
	names, err = api.getPathParams()
	a.NotError(err).
		Equal(names, map[string]bool{"id": true, "name": true})

	// 缺少 }
	api = &API{Path: "/path/{id"}
	names, err = api.getPathParams()
	a.Error(err).Nil(names)

	// 缺少 {
	api = &API{Path: "/path/id}"}
	names, err = api.getPathParams()
	a.Error(err).Nil(names)

	// 嵌套
	api = &API{Path: "/path/{{id}"}
	names, err = api.getPathParams()
	a.Error(err).Nil(names)
}

func TestDoc_parseAPI(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	l := newLexerString(`@api get /path xxxx
	@apiRequest object application/json summary
	@apiheader content-type required mimetype value
	@apiresponse 200 object application/json summary`)
	d.parseAPI(l)
	a.Equal(len(d.Apis), 1)

	// 重复内容。
	l = newLexerString(`@api get /path xxxx
	@apiNotExists object application/json summary
	@apiheader content-type required mimetype value
	@apiresponse 200 object application/json summary`)
	d.parseAPI(l)
	a.Equal(len(d.Apis), 1)
}

func TestAPI_parseAPI(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	l := newLexerString(`@api get /path xxxx
	@apirequest xxx`)
	tag := l.tag()
	api.parseAPI(l, tag)
	a.Equal(api.Method, "GET")

	api = &API{}
	l = newLexerString(`@apinotexists get /path xxxx
	@apirequest xxx`)
	tag = l.tag()
	api.parseAPI(l, tag)
}

func TestAPI_parseapi(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	// 缺少参数
	tag := newTagString("@api get /path")
	api.parseapi(nil, tag)

	tag = newTagString("@api get /path summary content")
	tag.File = "file.go"
	tag.Line = 111
	api.parseapi(nil, tag)
	a.Equal(api.Method, "GET").
		Equal(api.Path, "/path").
		Equal(api.Summary, "summary content").
		Equal(api.file, "file.go").
		Equal(api.line, 111)

	// 多次调用
	tag = newTagString("get /path summary content")
	api.parseapi(nil, tag)
}

func TestAPI_parseServer(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	// 不能为空
	tag := newTagString("@apiServers")
	api.parseServers(nil, tag)

	// 正常情况
	tag = newTagString("@apiServers server1, server2")
	api.parseServers(nil, tag)
	a.Equal(api.Servers, []string{"server1", "server2"})

	// 不能多次调用
	tag = newTagString("@apiServers s1")
	api.parseServers(nil, tag)
}

func TestAPI_parseTags(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	// 不能为空
	tag := newTagString("@apiTags")
	api.parseTags(nil, tag)

	tag = newTagString("@apiTags t1, t2")
	api.parseTags(nil, tag)
	a.Equal(api.Tags, []string{"t1", "t2"})

	// 不能多次调用
	tag = newTagString("@apiTags t1,t2")
	api.parseTags(nil, tag)

	// 两次同名
	api = &API{}
	tag = newTagString("@apiTags t1,t1")
	api.parseTags(nil, tag)
}

func TestAPI_parseDeperecated(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	// 不能为空
	tag := newTagString("@apiDeprecated")
	api.parseDeprecated(nil, tag)

	// 正常
	tag = newTagString("@apiDeprecated 临时取消")
	api.parseDeprecated(nil, tag)
	a.Equal(api.Deprecated, "临时取消")

	// 不能多次指定
	tag = newTagString("@apiDeprecated 临时取消")
	api.parseDeprecated(nil, tag)
}

func TestAPI_parseQuery(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	tag := newTagString("@apiQuery name string required  名称")
	api.parseQuery(nil, tag)
	a.Equal(len(api.Queries), 1)
	q := api.Queries[0]
	a.Equal(q.Name, "name").
		Equal(q.Type.Type, String).
		False(q.Optional).
		Equal(q.Summary, "名称")

	tag = newTagString("@apiQuery name1 string optional.v1  名称")
	api.parseQuery(nil, tag)
	a.Equal(len(api.Queries), 2)
	q = api.Queries[1]
	a.Equal(q.Name, "name1").
		Equal(q.Type.Type, String).
		True(q.Optional).
		Equal(q.Summary, "名称")

	// 同名参数
	tag = newTagString("@apiQuery name string optional.v1  名称")
	api.parseQuery(nil, tag)
}

func TestAPI_parseParam(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	tag := newTagString("@apiParam name string required  名称")
	api.parseParam(nil, tag)
	a.Equal(len(api.Params), 1)
	p := api.Params[0]
	a.Equal(p.Name, "name").
		Equal(p.Type.Type, String).
		False(p.Optional).
		Equal(p.Summary, "名称")

	tag = newTagString("@apiParam name1 string optional.v1  名称")
	api.parseParam(nil, tag)
	a.Equal(len(api.Params), 2)
	p = api.Params[1]
	a.Equal(p.Name, "name1").
		Equal(p.Type.Type, String).
		True(p.Optional).
		Equal(p.Summary, "名称")

	// 同名参数
	tag = newTagString("@apiParam name string optional.v1  名称")
	api.parseParam(nil, tag)
}

func TestAPI_parseResponse(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	l := newLexerString(`@apiHeader content-type optional 指定内容类型
	@apiParam id int required 唯一 ID
	@apiParam name string required 名称
	@apiParam nickname string optional 昵称
	@apiExample json 默认返回示例
	{
		"id": 1,
		"name": "name",
		"nickname": "nickname"
	}
	@apiUnknown xxx`)
	tag := newTagString(`@apiResponse array.object * 通用的返回内容定义`)

	api.parseRequest(l, tag)
	a.Equal(len(api.Requests), 1)
	req := api.Requests[0]
	a.Equal(req.Mimetype, "*")
	a.Equal(len(req.Headers), 1).
		Equal(req.Headers[0].Name, "content-type").
		Equal(req.Headers[0].Summary, "指定内容类型").
		True(req.Headers[0].Optional)
	a.NotNil(req.Type).
		Equal(req.Type.Type, Array)

	// 可以添加多次。
	api.parseRequest(l, tag)
	a.Equal(len(api.Requests), 2)
	req = api.Requests[1]
	a.Equal(req.Mimetype, "*")

	// 可选的描述内容
	tag = newTagString(`@apiResponse array.object application/json `)
	api.parseRequest(l, tag)
	a.Equal(len(api.Requests), 3)
	req = api.Requests[2]
	a.Equal(req.Mimetype, "application/json").
		Empty(req.Type.Description)

	// @apiRequest 格式错误
	tag = newTagString("xxxx")
	api.parseRequest(l, tag)
}

func TestNewParam(t *testing.T) {
	a := assert.New(t)

	tag := newTagString("@apiParam name string required  名称")
	p, ok := newParam(tag)
	a.True(ok).
		NotNil(p).
		Equal(p.Name, "name").
		Equal(p.Type.Type, String).
		False(p.Optional).
		Equal(p.Summary, "名称")

	tag = newTagString("@apiParam name string optional.v1  名称")
	p, ok = newParam(tag)
	a.True(ok).
		NotNil(p).
		Equal(p.Name, "name").
		Equal(p.Type.Type, String).
		True(p.Optional).
		Equal(p.Summary, "名称")

	tag = newTagString("@apiParam name string optional  名称")
	p, ok = newParam(tag)
	a.True(ok).
		NotNil(p).
		Equal(p.Name, "name").
		Equal(p.Type.Type, String).
		True(p.Optional).
		Equal(p.Summary, "名称")

	// 参数不够
	tag = newTagString("name ")
	p, ok = newParam(tag)
	a.False(ok).Nil(p)
}
