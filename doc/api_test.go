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

func TestAPI_parseCallback(t *testing.T) {
	a := assert.New(t)
	api := &API{}

	l := newLexerString(`@apiCallback GET description
@apiQuery page int page
@apiParam no string no
@apiRequest application/json * desc
@apiResponse 200 application/json * desc
@apiUnknown xx`)
	tag := l.tag()
	api.parseCallback(l, tag)
	a.Equal(api.Callback.Method, "GET").
		Equal(api.Callback.Summary, "description")
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
