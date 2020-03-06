// SPDX-License-Identifier: MIT

// Package spectest 提供了一个合法的 spec.APIDoc 对象
package spectest

import (
	"encoding/xml"
	"net/http"
	"path/filepath"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/internal/path"
	"github.com/caixw/apidoc/v6/spec"
)

// Filename 文档的文件名
const Filename = "index.xml"

// Get 返回 doc.APIDoc 对象
//
// 同时当前目录下的 index.xml 文件与此返回对象内容是相同的。
func Get() *spec.APIDoc {
	return &spec.APIDoc{
		Version: "1.0.1",
		Title:   "test",
		Description: spec.Richtext{
			Text: "<p>desc</p>",
			Type: spec.RichtextTypeHTML,
		},
		Servers: []*spec.Server{
			{
				URL:     "https://example.com/admin",
				Name:    "admin",
				Summary: "admin",
			},
			{
				URL:     "https://example.com",
				Name:    "client",
				Summary: "client",
			},
		},
		Tags: []*spec.Tag{
			{Name: "t1", Title: "t1"},
			{Name: "t2", Title: "t2"},
			{Name: "tag1", Title: "tag1"},
		},
		Mimetypes: []string{"application/json", "application/xml"},
		Apis: []*spec.API{
			{
				Method:  http.MethodGet,
				Tags:    []string{"t1", "t2"},
				Path:    &spec.Path{Path: "/users"},
				Servers: []string{"admin", "client"},
				Requests: []*spec.Request{
					{
						Summary: "request",
						Type:    spec.None,
						Headers: []*spec.Param{
							{
								Type:    spec.String,
								Name:    "authorization",
								Summary: "authorization",
							},
						},
						Examples: []*spec.Example{
							{
								Mimetype: "application/json",
								Content:  "xxx",
							},
						},
					},
				},
				Responses: []*spec.Request{
					{
						Description: spec.Richtext{Type: "html", Text: "<p>desc</p>"},
						Type:        spec.Object,
						Status:      http.StatusOK,
						Headers: []*spec.Param{
							{
								Type:    spec.String,
								Name:    "authorization",
								Summary: "authorization",
							},
						},
						Examples: []*spec.Example{
							{
								Mimetype: "application/json",
								Content:  "xxx",
							},
						},
						Items: []*spec.Param{
							{
								Summary: "summary",
								Type:    spec.String,
								Name:    "name",
							},
						},
					},
				},
			},
			{
				Method:     http.MethodPost,
				Tags:       []string{"t1", "tag1"},
				Path:       &spec.Path{Path: "/users"},
				Deprecated: "1.0.1",
				Summary:    "summary",
				Servers:    []string{"admin"},
				Responses: []*spec.Request{
					{
						Description: spec.Richtext{Type: "html", Text: "<p>desc</p>"},
						Type:        spec.None,
					},
				},
			},
		},
	}
}

// XML 获取 Get 返回对象的 XML 编码
func XML(a *assert.Assertion) []byte {
	data, err := xml.Marshal(Get())
	a.NotError(err).NotNil(data)

	return data
}

// Path 返回测试文件的绝对路径
//
// NOTE: 该文件与 Get() 对象的内容是相同的。
func Path(a *assert.Assertion) string {
	return pp(a, Filename)
}

// Dir 返回测试文件所在的目录
func Dir(a *assert.Assertion) string {
	return pp(a, "")
}

func pp(a *assert.Assertion, p string) string {
	p = path.CurrPath(p)
	p, err := filepath.Abs(p)
	a.NotError(err).NotEmpty(p)
	return p
}
