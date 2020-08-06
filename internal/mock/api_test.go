// SPDX-License-Identifier: MIT

package mock

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/qheader"

	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

func TestFindRequestByContentType(t *testing.T) {
	a := assert.New(t)
	data := []*struct {
		// 输入参数
		requests []*ast.Request
		ct       string

		// 返回参数
		index int
	}{
		{
			index: -1,
		},
		{
			requests: []*ast.Request{{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "application/json"}}}},
			ct:       "application/json",
			index:    0,
		},
		{
			requests: []*ast.Request{{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "application/json"}}}},
			ct:       "not/exists",
			index:    -1,
		},
		{
			requests: []*ast.Request{
				{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "application/json"}}},
				{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "text/xml"}}},
			},
			ct:    "text/xml",
			index: 1,
		},
		{
			requests: []*ast.Request{{}, {Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "text/xml"}}}},
			ct:       "text/xml",
			index:    1,
		},
		{ // 没有明确匹配，则匹配 none
			requests: []*ast.Request{{}, {Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "application/json"}}}},
			ct:       "text/xml",
			index:    0,
		},
	}

	for index, item := range data {
		req := findRequestByContentType(item.requests, item.ct)

		if item.index == -1 {
			a.Nil(req, "not nil at %d", index)
		} else {
			a.Equal(req, item.requests[item.index], "not equal at %d", index)
		}
	}
}

func TestFindResponseByAccept(t *testing.T) {
	a := assert.New(t)
	data := []*struct {
		// 输入参数
		mimetypes []string
		requests  []*ast.Request
		accepts   []*qheader.Header

		// 返回参数
		index int
		ct    string
	}{
		{ // 0
			index: -1,
		},
		{
			requests: []*ast.Request{},
			index:    -1,
		},
		{
			requests: []*ast.Request{{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "text/xml"}}}},
			accepts:  []*qheader.Header{{Value: "text/xml"}},
			index:    0,
			ct:       "text/xml",
		},
		{
			requests: []*ast.Request{
				{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "text/xml"}}},
				{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "application/json"}}},
			},
			accepts: []*qheader.Header{{Value: "text/xml"}},
			index:   0,
			ct:      "text/xml",
		},
		{
			requests: []*ast.Request{
				{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "text/xml"}}},
				{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "application/json"}}},
			},
			accepts: []*qheader.Header{{Value: "text/*"}},
			index:   0,
			ct:      "text/xml",
		},
		{ // 5
			requests: []*ast.Request{
				{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "text/xml"}}},
				{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "application/json"}}},
			},
			accepts: []*qheader.Header{{Value: "application/*"}},
			index:   1,
			ct:      "application/json",
		},
		{
			requests: []*ast.Request{
				{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "text/xml"}}},
				{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "application/json"}}},
			},
			accepts: []*qheader.Header{{Value: "*/*"}},
			index:   0,
			ct:      "text/xml",
		},
		{
			requests: []*ast.Request{
				{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "text/xml"}}},
				{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "application/json"}}},
			},
			accepts: []*qheader.Header{{Value: "*/*"}, {Value: "application/*"}},
			index:   0,
			ct:      "text/xml", // 第一个元素，匹配 */*
		},
		{
			mimetypes: []string{"text/xml"},
			requests:  []*ast.Request{},
			accepts:   []*qheader.Header{{Value: "*/*"}},
			index:     -1,
			ct:        "",
		},
		{
			mimetypes: []string{"text/xml"},
			requests:  []*ast.Request{{}},
			accepts:   []*qheader.Header{{Value: "*/*"}},
			index:     0,
			ct:        "text/xml",
		},
		{
			mimetypes: []string{"text/xml"},
			requests:  []*ast.Request{{}, {}},
			accepts:   []*qheader.Header{{Value: "*/*"}},
			index:     0,
			ct:        "text/xml",
		},
		{
			mimetypes: []string{"text/xml", "application/json"},
			requests:  []*ast.Request{{}},
			accepts:   []*qheader.Header{{Value: "application/*"}},
			index:     0,
			ct:        "application/json",
		},
		{
			mimetypes: []string{"text/xml", "application/json"},
			requests:  []*ast.Request{{}},
			accepts:   []*qheader.Header{{Value: "application/json"}},
			index:     0,
			ct:        "application/json",
		},
		{
			mimetypes: []string{"text/xml", "application/json"},
			requests:  []*ast.Request{{}},
			accepts:   []*qheader.Header{{Value: "application/*"}, {Value: "text/xml"}},
			index:     0,
			ct:        "text/xml",
		},
		{
			mimetypes: []string{"text/xml", "application/json"},
			requests:  []*ast.Request{{Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "application/json"}}}},
			accepts:   []*qheader.Header{{Value: "application/*"}},
			index:     0,
			ct:        "application/json",
		},
		{ // 任意值，匹配 mimetypes
			mimetypes: []string{"text/xml", "application/json"},
			requests:  []*ast.Request{{}},
			accepts:   []*qheader.Header{{Value: "font/*"}, {Value: "*/*"}},
			index:     0,
			ct:        "text/xml",
		},
	}

	for index, item := range data {
		mts := make([]*ast.Element, 0, len(item.mimetypes))
		for _, mimetype := range item.mimetypes {
			mts = append(mts, &ast.Element{Content: ast.Content{Value: mimetype}})
		}
		req, ct := findResponseByAccept(mts, item.requests, item.accepts)

		a.Equal(ct, item.ct, "not equal at %d,v1: %s,v2:%s", index, ct, item.ct)
		if item.index == -1 {
			a.Nil(req, "not nil at %d", index)
		} else {
			a.Equal(req, item.requests[item.index], "not equal at %d", index)
		}
	}
}

func TestValidRequest(t *testing.T) {
	a := assert.New(t)

	r := httptest.NewRequest(http.MethodGet, "/path", nil)
	a.Error(validRequest(nil, nil, r))

	// 匹配 json
	body := bytes.NewBufferString(dataWithHeader.JSON)
	r = httptest.NewRequest(http.MethodGet, "/path", body)
	r.Header.Set("content-type", "application/json")
	r.Header.Set("encoding", "xxx")
	a.NotError(validRequest(nil, []*ast.Request{dataWithHeader.Type}, r))

	// 匹配 xml
	body = bytes.NewBufferString(dataWithHeader.XML)
	r = httptest.NewRequest(http.MethodGet, "/path", body)
	r.Header.Set("content-type", "application/xml")
	r.Header.Set("encoding", "yyy")
	a.NotError(validRequest(nil, []*ast.Request{dataWithHeader.Type}, r))

	// 无法匹配 content-type
	body = bytes.NewBufferString(dataWithHeader.JSON)
	r = httptest.NewRequest(http.MethodGet, "/path", body)
	r.Header.Set("content-type", "not-exists")
	r.Header.Set("encoding", "xxx")
	a.Error(validRequest(nil, []*ast.Request{dataWithHeader.Type}, r))
}

func TestBuildResponse(t *testing.T) {
	a := assert.New(t)

	m := &mock{
		indent: indent,
		gen:    testOptions,
		doc:    &ast.APIDoc{},
	}

	resp, err := m.buildResponse(nil, nil)
	a.NotError(err).Nil(resp)

	// 匹配 json
	r := httptest.NewRequest(http.MethodGet, "/path", nil)
	r.Header.Set("accept", "application/json")
	r.Header.Set("encoding", "xxx")
	resp, err = m.buildResponse(dataWithHeader.Type, r)
	a.NotError(err).Equal(string(resp), dataWithHeader.JSON)

	// 匹配 xml
	r = httptest.NewRequest(http.MethodGet, "/path", nil)
	r.Header.Set("accept", "application/json;q=0.1,application/xml")
	r.Header.Set("encoding", "yyy")
	resp, err = m.buildResponse(dataWithHeader.Type, r)
	a.NotError(err).Equal(string(resp), dataWithHeader.XML)

	// 无法匹配 content-type
	r = httptest.NewRequest(http.MethodGet, "/path", nil)
	r.Header.Set("content-type", "not-exists")
	r.Header.Set("encoding", "xxx")
	resp, err = m.buildResponse(dataWithHeader.Type, r)
	a.Error(err).Nil(resp)
}

func TestValidSimpleParam(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		title string
		p     *ast.Param
		v     string
		err   bool
	}{
		{
			title: "nil",
			p:     nil,
			v:     "",
		},
		{
			title: "number",
			p:     &ast.Param{Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}}},
			v:     "-10.2",
		},
		{
			title: "number failed",
			p:     &ast.Param{Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}}},
			v:     "-xxx10.2",
			err:   true,
		},
		{
			title: "number with enum",
			p: &ast.Param{
				Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
				Enums: []*ast.Enum{
					{Value: &ast.Attribute{Value: xmlenc.String{Value: "1"}}},
					{Value: &ast.Attribute{Value: xmlenc.String{Value: "2"}}},
					{Value: &ast.Attribute{Value: xmlenc.String{Value: "10"}}},
				},
			},
			v: "1",
		},
		{
			title: "number with enum failed",
			p: &ast.Param{
				Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
				Enums: []*ast.Enum{
					{Value: &ast.Attribute{Value: xmlenc.String{Value: "1"}}},
					{Value: &ast.Attribute{Value: xmlenc.String{Value: "2"}}},
					{Value: &ast.Attribute{Value: xmlenc.String{Value: "10"}}},
				},
			},
			v:   "10001",
			err: true,
		},
		{
			title: "bool",
			p:     &ast.Param{Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeBool}}},
			v:     "false",
		},
		{
			title: "bool failed",
			p:     &ast.Param{Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeBool}}},
			v:     "-xxx-true",
			err:   true,
		},
		{
			title: "bool with optional",
			p: &ast.Param{
				Type:     &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeBool}},
				Optional: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
			},
			v: "",
		},
		{
			title: "bool with empty",
			p:     &ast.Param{Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeBool}}},
			v:     "",
			err:   true,
		},
		{
			title: "string",
			p:     &ast.Param{Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}}},
			v:     "-xxx10.2",
		},
		{
			title: "doc.None",
			p:     &ast.Param{Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNone}}},
			v:     "-xxx10.2",
			err:   true,
		},
	}

	for _, item := range data {
		err := validSimpleParam(item.p, item.title, item.v)
		if item.err {
			a.Error(err, "%s 并未返回错误值", item.title)
		} else {
			a.NotError(err, "%s 返回了非预期的错误值 %s", item.title, err)
		}
	}
}

func TestValidQueries(t *testing.T) {
	a := assert.New(t)
	data := []*struct {
		title string
		p     []*ast.Param
		r     *http.Request
		err   bool
	}{
		{
			title: "nil",
		},
		{
			title: "非数组",
			p: []*ast.Param{
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "k1"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
				},
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "k2"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
				},
			},
			r: httptest.NewRequest(http.MethodGet, "/users?k1=1&k2=2", nil),
		},
		{
			title: "非数组，格式不正确",
			p: []*ast.Param{
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "k1"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
				},
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "k2"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
				},
			},
			r:   httptest.NewRequest(http.MethodGet, "/users?k1=1&k2=not-number", nil),
			err: true,
		},
		{
			title: "数组-form",
			p: []*ast.Param{
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "k1"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
				},
				{
					Name:  &ast.Attribute{Value: xmlenc.String{Value: "k2"}},
					Type:  &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
					Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
				},
			},
			r: httptest.NewRequest(http.MethodGet, "/users?k1=1&k2=2&k2=3&k2=4", nil),
		},
		{
			title: "数组-form，格式不正确",
			p: []*ast.Param{
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "k1"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
				},
				{
					Name:  &ast.Attribute{Value: xmlenc.String{Value: "k2"}},
					Type:  &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
					Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
				},
			},
			r:   httptest.NewRequest(http.MethodGet, "/users?k1=1&k2=2&k2=3&k2=not-number", nil),
			err: true,
		},
		{
			title: "数组-array-style",
			p: []*ast.Param{
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "k1"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
				},
				{
					Name:       &ast.Attribute{Value: xmlenc.String{Value: "k2"}},
					Type:       &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
					Array:      &ast.BoolAttribute{Value: ast.Bool{Value: true}},
					ArrayStyle: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
				},
			},
			r: httptest.NewRequest(http.MethodGet, "/users?k1=1&k2=2,3,4", nil),
		},
		{
			title: "数组-array-style，格式不正确",
			p: []*ast.Param{
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "k1"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
				},
				{
					Name:       &ast.Attribute{Value: xmlenc.String{Value: "k2"}},
					Type:       &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
					Array:      &ast.BoolAttribute{Value: ast.Bool{Value: true}},
					ArrayStyle: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
				},
			},
			r:   httptest.NewRequest(http.MethodGet, "/users?k1=1&k2=2,3,not-number", nil),
			err: true,
		},
	}

	for _, item := range data {
		err := validQueries(item.p, item.r)
		if item.err {
			a.Error(err, "not error at %s", item.title)
		} else {
			a.NotError("err %s at %s", err, item.title)
		}
	}
}
