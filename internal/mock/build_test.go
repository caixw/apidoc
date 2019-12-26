// SPDX-License-Identifier: MIT

package mock

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/doc"
)

func TestValidRequest(t *testing.T) {
	a := assert.New(t)
	a.NotError(validRequest(nil, nil, ""))

	item := data[len(data)-2]

	// 匹配 json
	body := bytes.NewBufferString(item.JSON)
	r := httptest.NewRequest(http.MethodGet, "/path", body)
	r.Header.Set("content-type", "application/json")
	r.Header.Set("encoding", "xxx")
	a.NotError(validRequest(item.Type, r, "application/json"))

	// 匹配 xml
	body = bytes.NewBufferString(item.XML)
	r = httptest.NewRequest(http.MethodGet, "/path", body)
	r.Header.Set("content-type", "application/xml")
	r.Header.Set("encoding", "yyy")
	a.NotError(validRequest(item.Type, r, "application/xml"))

	// 无法匹配 content-type
	body = bytes.NewBufferString(`{"name":{"last":"l","first":"f"},"age":1}`)
	r = httptest.NewRequest(http.MethodGet, "/path", body)
	r.Header.Set("content-type", "not-exists")
	r.Header.Set("encoding", "xxx")
	a.Error(validRequest(item.Type, r, "not-exists"))
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

func TestValidParam(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		title string
		p     *doc.Param
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
			p:     &doc.Param{Type: doc.Number},
			v:     "-10.2",
		},
		{
			title: "number failed",
			p:     &doc.Param{Type: doc.Number},
			v:     "-xxx10.2",
			err:   true,
		},
		{
			title: "number with enum",
			p: &doc.Param{
				Type: doc.Number,
				Enums: []*doc.Enum{
					{Value: "1"},
					{Value: "2"},
					{Value: "10"},
				},
			},
			v: "1",
		},
		{
			title: "number with enum failed",
			p: &doc.Param{
				Type: doc.Number,
				Enums: []*doc.Enum{
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
			p:     &doc.Param{Type: doc.Bool},
			v:     "false",
		},
		{
			title: "bool failed",
			p:     &doc.Param{Type: doc.Bool},
			v:     "-xxx-true",
			err:   true,
		},
		{
			title: "bool with optional",
			p:     &doc.Param{Type: doc.Bool, Optional: true},
			v:     "",
		},
		{
			title: "bool with empty",
			p:     &doc.Param{Type: doc.Bool},
			v:     "",
			err:   true,
		},
		{
			title: "string",
			p:     &doc.Param{Type: doc.String},
			v:     "-xxx10.2",
		},
		{
			title: "doc.None",
			p:     &doc.Param{Type: doc.None},
			v:     "-xxx10.2",
			err:   true,
		},
	}

	for _, item := range data {
		err := validParam(item.p, item.v)
		if item.err {
			a.Error(err, "%s 并未返回错误值", item.title)
		} else {
			a.NotError(err, "%s 返回了非预期的错误值 %s", item.title, err)
		}
	}
}
