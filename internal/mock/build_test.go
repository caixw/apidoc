// SPDX-License-Identifier: MIT

package mock

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/qheader"

	"github.com/caixw/apidoc/v6/spec"
)

func TestFindRequestByContentType(t *testing.T) {
	a := assert.New(t)
	data := []*struct {
		// 输入参数
		requests []*spec.Request
		ct       string

		// 返回参数
		index int
	}{
		{
			index: -1,
		},
		{
			requests: []*spec.Request{{Mimetype: "application/json"}},
			ct:       "application/json",
			index:    0,
		},
		{
			requests: []*spec.Request{{Mimetype: "application/json"}},
			ct:       "not/exists",
			index:    -1,
		},
		{
			requests: []*spec.Request{{Mimetype: "application/json"}, {Mimetype: "text/xml"}},
			ct:       "text/xml",
			index:    1,
		},
		{
			requests: []*spec.Request{{}, {Mimetype: "text/xml"}},
			ct:       "text/xml",
			index:    1,
		},
		{ // 没有明确匹配，则匹配 none
			requests: []*spec.Request{{}, {Mimetype: "application/json"}},
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
		requests  []*spec.Request
		accepts   []*qheader.Header

		// 返回参数
		index int
		ct    string
	}{
		{ // 0
			index: -1,
		},
		{
			requests: []*spec.Request{},
			index:    -1,
		},
		{
			requests: []*spec.Request{{Mimetype: "text/xml"}},
			accepts:  []*qheader.Header{{Value: "text/xml"}},
			index:    0,
			ct:       "text/xml",
		},
		{
			requests: []*spec.Request{{Mimetype: "text/xml"}, {Mimetype: "application/json"}},
			accepts:  []*qheader.Header{{Value: "text/xml"}},
			index:    0,
			ct:       "text/xml",
		},
		{
			requests: []*spec.Request{{Mimetype: "text/xml"}, {Mimetype: "application/json"}},
			accepts:  []*qheader.Header{{Value: "text/*"}},
			index:    0,
			ct:       "text/xml",
		},
		{ // 5
			requests: []*spec.Request{{Mimetype: "text/xml"}, {Mimetype: "application/json"}},
			accepts:  []*qheader.Header{{Value: "application/*"}},
			index:    1,
			ct:       "application/json",
		},
		{
			requests: []*spec.Request{{Mimetype: "text/xml"}, {Mimetype: "application/json"}},
			accepts:  []*qheader.Header{{Value: "*/*"}},
			index:    0,
			ct:       "text/xml",
		},
		{
			requests: []*spec.Request{{Mimetype: "text/xml"}, {Mimetype: "application/json"}},
			accepts:  []*qheader.Header{{Value: "*/*"}, {Value: "application/*"}},
			index:    0,
			ct:       "text/xml", // 第一个元素，匹配 */*
		},
		{
			mimetypes: []string{"text/xml"},
			requests:  []*spec.Request{},
			accepts:   []*qheader.Header{{Value: "*/*"}},
			index:     -1,
			ct:        "",
		},
		{
			mimetypes: []string{"text/xml"},
			requests:  []*spec.Request{{}},
			accepts:   []*qheader.Header{{Value: "*/*"}},
			index:     0,
			ct:        "text/xml",
		},
		{
			mimetypes: []string{"text/xml"},
			requests:  []*spec.Request{{}, {}},
			accepts:   []*qheader.Header{{Value: "*/*"}},
			index:     0,
			ct:        "text/xml",
		},
		{
			mimetypes: []string{"text/xml", "application/json"},
			requests:  []*spec.Request{{}},
			accepts:   []*qheader.Header{{Value: "application/*"}},
			index:     0,
			ct:        "application/json",
		},
		{
			mimetypes: []string{"text/xml", "application/json"},
			requests:  []*spec.Request{{}},
			accepts:   []*qheader.Header{{Value: "application/json"}},
			index:     0,
			ct:        "application/json",
		},
		{
			mimetypes: []string{"text/xml", "application/json"},
			requests:  []*spec.Request{{}},
			accepts:   []*qheader.Header{{Value: "application/*"}, {Value: "text/xml"}},
			index:     0,
			ct:        "text/xml",
		},
		{
			mimetypes: []string{"text/xml", "application/json"},
			requests:  []*spec.Request{{Mimetype: "application/json"}},
			accepts:   []*qheader.Header{{Value: "application/*"}},
			index:     0,
			ct:        "application/json",
		},
		{ // 任意值，匹配 mimetypes
			mimetypes: []string{"text/xml", "application/json"},
			requests:  []*spec.Request{{}},
			accepts:   []*qheader.Header{{Value: "font/*"}, {Value: "*/*"}},
			index:     0,
			ct:        "text/xml",
		},
	}

	for index, item := range data {
		req, ct := findResponseByAccept(item.mimetypes, item.requests, item.accepts)

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
	a.Error(validRequest(nil, r))

	item := data[len(data)-2]

	// 匹配 json
	body := bytes.NewBufferString(item.JSON)
	r = httptest.NewRequest(http.MethodGet, "/path", body)
	r.Header.Set("content-type", "application/json")
	r.Header.Set("encoding", "xxx")
	a.NotError(validRequest([]*spec.Request{item.Type}, r))

	// 匹配 xml
	body = bytes.NewBufferString(item.XML)
	r = httptest.NewRequest(http.MethodGet, "/path", body)
	r.Header.Set("content-type", "application/xml")
	r.Header.Set("encoding", "yyy")
	a.NotError(validRequest([]*spec.Request{item.Type}, r))

	// 无法匹配 content-type
	body = bytes.NewBufferString(`{"name":{"last":"l","first":"f"},"age":1}`)
	r = httptest.NewRequest(http.MethodGet, "/path", body)
	r.Header.Set("content-type", "not-exists")
	r.Header.Set("encoding", "xxx")
	a.Error(validRequest([]*spec.Request{item.Type}, r))
}

func TestBuildResponse(t *testing.T) {
	a := assert.New(t)

	resp, err := buildResponse(nil, nil)
	a.NotError(err).Nil(resp)

	item := data[len(data)-2]

	// 匹配 json
	r := httptest.NewRequest(http.MethodGet, "/path", nil)
	r.Header.Set("accept", "application/json")
	r.Header.Set("encoding", "xxx")
	resp, err = buildResponse(item.Type, r)
	a.NotError(err).Equal(string(resp), item.JSON)

	// 匹配 xml
	r = httptest.NewRequest(http.MethodGet, "/path", nil)
	r.Header.Set("accept", "application/xml")
	r.Header.Set("encoding", "yyy")
	resp, err = buildResponse(item.Type, r)
	a.NotError(err).Equal(string(resp), item.XML)

	// 无法匹配 content-type
	r = httptest.NewRequest(http.MethodGet, "/path", nil)
	r.Header.Set("content-type", "not-exists")
	r.Header.Set("encoding", "xxx")
	resp, err = buildResponse(item.Type, r)
	a.Error(err).Nil(resp)
}

func TestValidSimpleParam(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		title string
		p     *spec.Param
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
			p:     &spec.Param{Type: spec.Number},
			v:     "-10.2",
		},
		{
			title: "number failed",
			p:     &spec.Param{Type: spec.Number},
			v:     "-xxx10.2",
			err:   true,
		},
		{
			title: "number with enum",
			p: &spec.Param{
				Type: spec.Number,
				Enums: []*spec.Enum{
					{Value: "1"},
					{Value: "2"},
					{Value: "10"},
				},
			},
			v: "1",
		},
		{
			title: "number with enum failed",
			p: &spec.Param{
				Type: spec.Number,
				Enums: []*spec.Enum{
					{Value: "1"},
					{Value: "2"},
					{Value: "10"},
				},
			},
			v:   "10001",
			err: true,
		},
		{
			title: "bool",
			p:     &spec.Param{Type: spec.Bool},
			v:     "false",
		},
		{
			title: "bool failed",
			p:     &spec.Param{Type: spec.Bool},
			v:     "-xxx-true",
			err:   true,
		},
		{
			title: "bool with optional",
			p:     &spec.Param{Type: spec.Bool, Optional: true},
			v:     "",
		},
		{
			title: "bool with empty",
			p:     &spec.Param{Type: spec.Bool},
			v:     "",
			err:   true,
		},
		{
			title: "string",
			p:     &spec.Param{Type: spec.String},
			v:     "-xxx10.2",
		},
		{
			title: "doc.None",
			p:     &spec.Param{Type: spec.None},
			v:     "-xxx10.2",
			err:   true,
		},
	}

	for _, item := range data {
		err := validSimpleParam(item.p, item.v)
		if item.err {
			a.Error(err, "%s 并未返回错误值", item.title)
		} else {
			a.NotError(err, "%s 返回了非预期的错误值 %s", item.title, err)
		}
	}
}

func TestValidQueryArrayParam(t *testing.T) {
	a := assert.New(t)
	data := []*struct {
		title string
		p     []*spec.Param
		r     *http.Request
		err   bool
	}{
		{
			title: "nil",
		},
		{
			title: "非数组",
			p: []*spec.Param{
				{Name: "k1", Type: spec.String},
				{Name: "k2", Type: spec.Number},
			},
			r: httptest.NewRequest(http.MethodGet, "/users?k1=1&k2=2", nil),
		},
		{
			title: "非数组，格式不正确",
			p: []*spec.Param{
				{Name: "k1", Type: spec.String},
				{Name: "k2", Type: spec.Number},
			},
			r:   httptest.NewRequest(http.MethodGet, "/users?k1=1&k2=not-number", nil),
			err: true,
		},
		{
			title: "数组-form",
			p: []*spec.Param{
				{Name: "k1", Type: spec.String},
				{Name: "k2", Type: spec.Number, Array: true},
			},
			r: httptest.NewRequest(http.MethodGet, "/users?k1=1&k2=2&k2=3&k2=4", nil),
		},
		{
			title: "数组-form，格式不正确",
			p: []*spec.Param{
				{Name: "k1", Type: spec.String},
				{Name: "k2", Type: spec.Number, Array: true},
			},
			r:   httptest.NewRequest(http.MethodGet, "/users?k1=1&k2=2&k2=3&k2=not-number", nil),
			err: true,
		},
		{
			title: "数组-array-style",
			p: []*spec.Param{
				{Name: "k1", Type: spec.String},
				{Name: "k2", Type: spec.Number, Array: true, ArrayStyle: true},
			},
			r: httptest.NewRequest(http.MethodGet, "/users?k1=1&k2=2,3,4", nil),
		},
		{
			title: "数组-array-style，格式不正确",
			p: []*spec.Param{
				{Name: "k1", Type: spec.String},
				{Name: "k2", Type: spec.Number, Array: true, ArrayStyle: true},
			},
			r:   httptest.NewRequest(http.MethodGet, "/users?k1=1&k2=2,3,not-number", nil),
			err: true,
		},
	}

	for _, item := range data {
		err := validQueryArrayParam(item.p, item.r)
		if item.err {
			a.Error(err, "not error at %s", item.title)
		} else {
			a.NotError("err %s at %s", err, item.title)
		}
	}
}
