// SPDX-License-Identifier: MIT

// Package doctest 提供了一个合法的 doc.Doc 对象
package doctest

import (
	"net/http"

	"github.com/caixw/apidoc/v5/doc"
)

// Get 返回 doc.Doc 对象
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
		Tags: []*doc.Tag{{Name: "t1"}, {Name: "t2"}, {Name: "tag1"}},
		Apis: []*doc.API{
			{
				Method:  http.MethodGet,
				Tags:    []string{"t1", "t2"},
				Path:    &doc.Path{Path: "/users"},
				Servers: []string{"admin", "client"},
				Requests: []*doc.Request{
					{
						Summary: "request",
						Type:    doc.Object,
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
						Type:        doc.Object,
					},
				},
			},
		},
	}
}
