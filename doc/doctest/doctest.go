// SPDX-License-Identifier: MIT

// Package doctest 提供了一个合法的 doc.Doc 对象
package doctest

import (
	"encoding/xml"
	"net/http"
	"path/filepath"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/path"
)

// Filename 文档的文件名
const Filename = "index.xml"

// Get 返回 doc.Doc 对象
//
// 同时当前目录下的 index.xml 文件与此返回对象内容是相同的。
func Get() *doc.Doc {
	return &doc.Doc{
		Version: "1.0.1",
		Title:   "test",
		Description: doc.Richtext{
			Text: "<p>desc</p>",
			Type: doc.RichtextTypeHTML,
		},
		Servers: []*doc.Server{
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
		Tags: []*doc.Tag{
			{Name: "t1", Title: "t1"},
			{Name: "t2", Title: "t2"},
			{Name: "tag1", Title: "tag1"},
		},
		Mimetypes: []string{"application/json", "application/xml"},
		Apis: []*doc.API{
			{
				Method:  http.MethodGet,
				Tags:    []string{"t1", "t2"},
				Path:    &doc.Path{Path: "/users"},
				Servers: []string{"admin", "client"},
				Requests: []*doc.Request{
					{
						Summary: "request",
						Type:    doc.None,
						Headers: []*doc.Param{
							{
								Type:    doc.String,
								Name:    "authorization",
								Summary: "authorization",
							},
						},
						Examples: []*doc.Example{
							{
								Mimetype: "application/json",
								Content:  "xxx",
							},
						},
					},
				},
				Responses: []*doc.Request{
					{
						Description: doc.Richtext{Type: "html", Text: "<p>desc</p>"},
						Type:        doc.Object,
						Status:      http.StatusOK,
						Headers: []*doc.Param{
							{
								Type:    doc.String,
								Name:    "authorization",
								Summary: "authorization",
							},
						},
						Examples: []*doc.Example{
							{
								Mimetype: "application/json",
								Content:  "xxx",
							},
						},
						Items: []*doc.Param{
							{
								Summary: "summary",
								Type:    doc.String,
								Name:    "name",
							},
						},
					},
				},
			},
			{
				Method:     http.MethodPost,
				Tags:       []string{"t1", "tag1"},
				Path:       &doc.Path{Path: "/users"},
				Deprecated: "1.0.1",
				Summary:    "summary",
				Servers:    []string{"admin"},
				Responses: []*doc.Request{
					{
						Description: doc.Richtext{Type: "html", Text: "<p>desc</p>"},
						Type:        doc.None,
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
	p := path.CurrPath(Filename)
	p, err := filepath.Abs(p)
	a.NotError(err).NotEmpty(p)
	return p
}
