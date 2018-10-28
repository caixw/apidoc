// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/doc/schema"
)

func TestDoc_parseAPI(t *testing.T) {
	a := assert.New(t)
	d := &Doc{}

	l := newLexer(`@api get /path xxxx
	@apiRequest object application/json summary
	@apiheader content-type required mimetype value
	@apiresponse 200 object application/json summary`)
	a.NotError(d.parseAPI(l))

	l = newLexer(`@api get /path xxxx
	@apiNotExists object application/json summary
	@apiheader content-type required mimetype value
	@apiresponse 200 object application/json summary`)
	a.Error(d.parseAPI(l))
}

func TestAPI_parseAPI(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	l := newLexer(`@api get /path xxxx
	@apirequest xxx`)
	tag := l.Tag()
	a.NotError(api.parseAPI(l, tag))
	a.Equal(api.Method, "GET")

	api = &API{}
	l = newLexer(`@apinotexists get /path xxxx
	@apirequest xxx`)
	tag = l.Tag()
	a.NotError(api.parseAPI(l, tag))
}

func TestAPI_parseapi(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	// 缺少参数
	a.Error(api.parseapi(nil, newTag("get /path")))

	a.NotError(api.parseapi(nil, newTag("get /path summary content")))
	a.Equal(api.Method, "GET").
		Equal(api.Path, "/path").
		Equal(api.Summary, "summary content")

	// 多次调用
	a.Error(api.parseapi(nil, newTag("get /path summary content")))
}

func TestAPI_parseServer(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	// 不能为空
	a.Error(api.parseServer(nil, newTag("")))

	a.NotError(api.parseServer(nil, newTag("server1")))
	a.Equal(api.Server, "server1")

	// 不能多次调用
	a.Error(api.parseServer(nil, newTag("s1")))
}

func TestAPI_parseTags(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	// 不能为空
	a.Error(api.parseTags(nil, newTag("")))

	a.NotError(api.parseTags(nil, newTag("t1,t2")))
	a.Equal(api.Tags, []string{"t1", "t2"})

	// 不能多次调用
	a.Error(api.parseTags(nil, newTag("t1,t2")))
}

func TestAPI_parseDeperecated(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	// 不能为空
	a.Error(api.parseDeprecated(nil, newTag("")))

	// 正常
	a.NotError(api.parseDeprecated(nil, newTag("临时取消")))
	a.Equal(api.Deprecated, "临时取消")

	// 不能多次指定
	a.Error(api.parseDeprecated(nil, newTag("临时取消")))
}

func TestAPI_parseQuery(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	a.NotError(api.parseQuery(nil, newTag("name string required  名称")))
	a.Equal(len(api.Queries), 1)
	q := api.Queries[0]
	a.Equal(q.Name, "name").
		Equal(q.Type.Type, schema.String).
		False(q.Optional).
		Equal(q.Summary, "名称")

	a.NotError(api.parseQuery(nil, newTag("name string optional.v1  名称")))
	a.Equal(len(api.Queries), 2)
	q = api.Queries[1]
	a.Equal(q.Name, "name").
		Equal(q.Type.Type, schema.String).
		True(q.Optional).
		Equal(q.Summary, "名称")
}

func TestAPI_parseParam(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	a.NotError(api.parseParam(nil, newTag("name string required  名称")))
	a.Equal(len(api.Params), 1)
	p := api.Params[0]
	a.Equal(p.Name, "name").
		Equal(p.Type.Type, schema.String).
		False(p.Optional).
		Equal(p.Summary, "名称")

	a.NotError(api.parseParam(nil, newTag("name string optional.v1  名称")))
	a.Equal(len(api.Params), 2)
	p = api.Params[1]
	a.Equal(p.Name, "name").
		Equal(p.Type.Type, schema.String).
		True(p.Optional).
		Equal(p.Summary, "名称")
}

func TestAPI_parseResponse(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	l := newLexer(`@apiHeader content-type optional 指定内容类型
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
	tag := newTag(`array.object * 通用的返回内容定义`)

	a.NotError(api.parseRequest(l, tag)).
		Equal(len(api.Requests), 1)
	req := api.Requests[0]
	a.Equal(req.Mimetype, "*")
	a.Equal(len(req.Headers), 1).
		Equal(req.Headers[0].Name, "content-type").
		Equal(req.Headers[0].Summary, "指定内容类型").
		True(req.Headers[0].Optional)
	a.NotNil(req.Type).
		Equal(req.Type.Type, schema.Array)

	// 可以添加多次。
	a.NotError(api.parseRequest(l, tag)).
		Equal(len(api.Requests), 2)
	req = api.Requests[0]
	a.Equal(req.Mimetype, "*")

	// @apiRequest 格式错误
	a.Error(api.parseRequest(l, newTag("xxxx")))
}

func TestNewParam(t *testing.T) {
	a := assert.New(t)

	p, err := newParam(newTag("name string required  名称"))
	a.NotError(err).
		NotNil(p).
		Equal(p.Name, "name").
		Equal(p.Type.Type, schema.String).
		False(p.Optional).
		Equal(p.Summary, "名称")

	p, err = newParam(newTag("name string optional.v1  名称"))
	a.NotError(err).
		NotNil(p).
		Equal(p.Name, "name").
		Equal(p.Type.Type, schema.String).
		True(p.Optional).
		Equal(p.Summary, "名称")

	p, err = newParam(newTag("name string optional  名称"))
	a.NotError(err).
		NotNil(p).
		Equal(p.Name, "name").
		Equal(p.Type.Type, schema.String).
		True(p.Optional).
		Equal(p.Summary, "名称")

	// 参数不够
	p, err = newParam(newTag("name "))
	a.Error(err).Nil(p)
}
