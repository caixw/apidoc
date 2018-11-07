// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

func TestBody_parseExample(t *testing.T) {
	a := assert.New(t)
	body := &Body{}

	tag := newTagString(`@apiExample application/json summary text
{
	"id": 1,
	"name": "name"
}`)
	body.parseExample(tag)
	e := body.Examples[0]
	a.Equal(e.Mimetype, "application/json").
		Equal(e.Summary, "summary text").
		Equal(e.Value, `{
	"id": 1,
	"name": "name"
}`)

	// 长度不够
	tag = newTagString("application/json")
	body.parseExample(tag)
}

func TestBody_parseHeader(t *testing.T) {
	a := assert.New(t)
	body := &Body{}

	tag := newTagString(`@apiExample content-type required json 或是 xml`)
	body.parseHeader(tag)
	h := body.Headers[0]
	a.Equal(h.Summary, "json 或是 xml").
		Equal(h.Name, "content-type").
		False(h.Optional)

	tag = newTagString(`@apiExample ETag optional etag`)
	body.parseHeader(tag)
	h = body.Headers[1]
	a.Equal(h.Summary, "etag").
		Equal(h.Name, "ETag").
		True(h.Optional)

	// 长度不够
	tag = newTagString("ETag")
	body.parseHeader(tag)
}

func TestIsOptional(t *testing.T) {
	a := assert.New(t)

	a.False(isOptional(requiredBytes))
	a.False(isOptional(bytes.ToUpper(requiredBytes)))
	a.True(isOptional([]byte("optional")))
	a.True(isOptional([]byte("Optional")))
}

func TestResponses_parseResponse(t *testing.T) {
	a := assert.New(t)
	d := &responses{}

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
	tag := newTagString(`@apiResponse 200 array.object * 通用的返回内容定义`)

	d.parseResponse(l, tag)
	a.Equal(len(d.Responses), 1)
	resp := d.Responses[0]
	a.Equal(resp.Status, 200).
		Equal(resp.Mimetype, "*").
		Equal(resp.Type.Description, "通用的返回内容定义")
	a.Equal(len(resp.Headers), 1).
		Equal(resp.Headers[0].Name, "content-type").
		Equal(resp.Headers[0].Summary, "指定内容类型").
		True(resp.Headers[0].Optional)
	a.NotNil(resp.Type).
		Equal(resp.Type.Type, Array)

	// 可以添加多次。
	d.parseResponse(l, tag)
	a.Equal(len(d.Responses), 2)
	resp = d.Responses[1]
	a.Equal(resp.Status, 200).
		Equal(resp.Mimetype, "*")

	// 忽略可选参数
	tag = newTagString(`@apiResponse 200 array.object * `)
	d.parseResponse(l, tag)
	a.Equal(len(d.Responses), 3)
	resp = d.Responses[2]
	a.Equal(resp.Status, 200).
		Equal(resp.Mimetype, "*").
		Empty(resp.Type.Description)
}

func TestAPICallback_parseRequest(t *testing.T) {
	a := assert.New(t)
	c := &apiCallback{}

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

	c.parseRequest(l, tag)
	a.Equal(len(c.Requests), 1)
	req := c.Requests[0]
	a.Equal(req.Mimetype, "*")
	a.Equal(len(req.Headers), 1).
		Equal(req.Headers[0].Name, "content-type").
		Equal(req.Headers[0].Summary, "指定内容类型").
		True(req.Headers[0].Optional)
	a.NotNil(req.Type).
		Equal(req.Type.Type, Array)

	// 可以添加多次。
	c.parseRequest(l, tag)
	a.Equal(len(c.Requests), 2)
	req = c.Requests[1]
	a.Equal(req.Mimetype, "*")

	// 可选的描述内容
	tag = newTagString(`@reqsResponse array.object application/json `)
	c.parseRequest(l, tag)
	a.Equal(len(c.Requests), 3)
	req = c.Requests[2]
	a.Equal(req.Mimetype, "application/json").
		Empty(req.Type.Description)

	// @reqsRequest 格式错误
	tag = newTagString("xxxx")
	c.parseRequest(l, tag)
}

func TestAPICallback_parseQuery(t *testing.T) {
	a := assert.New(t)
	api := &apiCallback{}

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

func TestAPICallback_parseParam(t *testing.T) {
	a := assert.New(t)
	api := &apiCallback{}

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
