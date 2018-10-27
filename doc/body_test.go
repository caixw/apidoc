// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/doc/schema"
)

func TestBody_parseExample(t *testing.T) {
	a := assert.New(t)
	body := &Body{}

	a.NotError(body.parseExample(newTag(`application/json summary text
{
	"id": 1,
	"name": "name"
}`)))
	e := body.Examples[0]
	a.Equal(e.Mimetype, "application/json").
		Equal(e.Summary, "summary text").
		Equal(e.Value, `{
	"id": 1,
	"name": "name"
}`)

	// 长度不够
	a.Error(body.parseExample(newTag("application/json")))
}

func TestBody_parseHeader(t *testing.T) {
	a := assert.New(t)
	body := &Body{}

	a.NotError(body.parseHeader(newTag(`content-type required json 或是 xml`)))
	h := body.Headers[0]
	a.Equal(h.Summary, "json 或是 xml").
		Equal(h.Name, "content-type").
		False(h.Optional)

	a.NotError(body.parseHeader(newTag(`ETag optional etag`)))
	h = body.Headers[1]
	a.Equal(h.Summary, "etag").
		Equal(h.Name, "ETag").
		True(h.Optional)

	// 长度不够
	a.Error(body.parseHeader(newTag("ETag")))
}

func TestIsOptional(t *testing.T) {
	a := assert.New(t)

	a.False(isOptional(requiredBytes))
	a.False(isOptional(bytes.ToUpper(requiredBytes)))
	a.True(isOptional([]byte("optional")))
	a.True(isOptional([]byte("Optional")))
}

func TestNewResponse(t *testing.T) {
	a := assert.New(t)
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
	tag := newTag(`200 array.object * 通用的返回内容定义`)

	resp, err := newResponse(l, tag)
	a.NotError(err).NotNil(resp)
	a.Equal(resp.Status, 200).
		Equal(resp.Mimetype, "*")
	a.Equal(len(resp.Headers), 1).
		Equal(resp.Headers[0].Name, "content-type").
		Equal(resp.Headers[0].Summary, "指定内容类型").
		True(resp.Headers[0].Optional)
	a.NotNil(resp.Type).
		Equal(resp.Type.Type, schema.Array)
}
