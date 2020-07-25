// SPDX-License-Identifier: MIT

// Package asttest 提供了一个合法的 ast.APIDoc 对象
package asttest

import (
	"net/http"
	"path/filepath"

	"github.com/issue9/assert"
	"github.com/issue9/source"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

// Filename 文档的文件名
const Filename = "index.xml"

// Get 返回 doc.APIDoc 对象
//
// 同时当前目录下的 index.xml 文件与此返回对象内容是相同的。
func Get() *ast.APIDoc {
	return &ast.APIDoc{
		APIDoc:  &ast.APIDocVersionAttribute{Value: xmlenc.String{Value: ast.Version}},
		Version: &ast.VersionAttribute{Value: xmlenc.String{Value: "1.0.1"}},
		Title:   &ast.Element{Content: ast.Content{Value: "test"}},
		Description: &ast.Richtext{
			Text: &ast.CData{Value: xmlenc.String{Value: "<p>desc</p>"}},
			Type: &ast.Attribute{Value: xmlenc.String{Value: ast.RichtextTypeHTML}},
		},
		Servers: []*ast.Server{
			{
				URL:     &ast.Attribute{Value: xmlenc.String{Value: "https://example.com/admin"}},
				Name:    &ast.Attribute{Value: xmlenc.String{Value: "admin"}},
				Summary: &ast.Attribute{Value: xmlenc.String{Value: "admin"}},
			},
			{
				URL:     &ast.Attribute{Value: xmlenc.String{Value: "https://example.com"}},
				Name:    &ast.Attribute{Value: xmlenc.String{Value: "client"}},
				Summary: &ast.Attribute{Value: xmlenc.String{Value: "client"}},
			},
		},
		Tags: []*ast.Tag{
			{
				Name:  &ast.Attribute{Value: xmlenc.String{Value: "t1"}},
				Title: &ast.Attribute{Value: xmlenc.String{Value: "t1"}},
			},
			{
				Name:  &ast.Attribute{Value: xmlenc.String{Value: "t2"}},
				Title: &ast.Attribute{Value: xmlenc.String{Value: "t2"}},
			},
			{
				Name:  &ast.Attribute{Value: xmlenc.String{Value: "tag1"}},
				Title: &ast.Attribute{Value: xmlenc.String{Value: "tag1"}},
			},
		},
		Mimetypes: []*ast.Element{
			{Content: ast.Content{Value: "application/json"}},
			{Content: ast.Content{Value: "application/xml"}},
		},
		APIs: []*ast.API{
			{
				Method: &ast.MethodAttribute{Value: xmlenc.String{Value: http.MethodGet}},
				Tags: []*ast.Element{
					{Content: ast.Content{Value: "t1"}},
					{Content: ast.Content{Value: "t2"}},
				},
				Path: &ast.Path{Path: &ast.Attribute{Value: xmlenc.String{Value: "/users"}}},
				Servers: []*ast.Element{
					{Content: ast.Content{Value: "admin"}},
				},
				Requests: []*ast.Request{
					{
						Summary: &ast.Attribute{Value: xmlenc.String{Value: "request"}},
						Headers: []*ast.Param{
							{
								Type:    &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
								Name:    &ast.Attribute{Value: xmlenc.String{Value: "authorization"}},
								Summary: &ast.Attribute{Value: xmlenc.String{Value: "authorization"}},
							},
						},
						Examples: []*ast.Example{
							{
								Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "application/json"}},
								Content:  &ast.ExampleValue{Value: xmlenc.String{Value: "xxx"}},
							},
						},
					},
				},
				Responses: []*ast.Request{
					{
						Description: &ast.Richtext{
							Type: &ast.Attribute{Value: xmlenc.String{Value: "html"}},
							Text: &ast.CData{Value: xmlenc.String{Value: "<p>desc</p>"}},
						},
						Type:   &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
						Status: &ast.StatusAttribute{Value: ast.Number{Int: http.StatusOK}},
						Headers: []*ast.Param{
							{
								Type:    &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
								Name:    &ast.Attribute{Value: xmlenc.String{Value: "authorization"}},
								Summary: &ast.Attribute{Value: xmlenc.String{Value: "authorization"}},
							},
						},
						Examples: []*ast.Example{
							{
								Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "application/json"}},
								Content:  &ast.ExampleValue{Value: xmlenc.String{Value: "xxx"}},
							},
						},
						Items: []*ast.Param{
							{
								Type:    &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
								Name:    &ast.Attribute{Value: xmlenc.String{Value: "id"}},
								Summary: &ast.Attribute{Value: xmlenc.String{Value: "ID"}},
							},
							{
								Summary: &ast.Attribute{Value: xmlenc.String{Value: "summary"}},
								Type:    &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
								Name:    &ast.Attribute{Value: xmlenc.String{Value: "name"}},
							},
						},
					},
				},
			},
			{
				Method: &ast.MethodAttribute{Value: xmlenc.String{Value: http.MethodPost}},
				Tags: []*ast.Element{
					{Content: ast.Content{Value: "t1"}},
					{Content: ast.Content{Value: "tag1"}},
				},
				Path:       &ast.Path{Path: &ast.Attribute{Value: xmlenc.String{Value: "/users"}}},
				Deprecated: &ast.VersionAttribute{Value: xmlenc.String{Value: "1.0.1"}},
				Summary:    &ast.Attribute{Value: xmlenc.String{Value: "summary"}},
				Servers: []*ast.Element{
					{Content: ast.Content{Value: "admin"}},
					{Content: ast.Content{Value: "client"}},
				},
				Requests: []*ast.Request{
					{
						Name:    &ast.Attribute{Value: xmlenc.String{Value: "root"}},
						Summary: &ast.Attribute{Value: xmlenc.String{Value: "request"}},
						Headers: []*ast.Param{
							{
								Type:    &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
								Name:    &ast.Attribute{Value: xmlenc.String{Value: "authorization"}},
								Summary: &ast.Attribute{Value: xmlenc.String{Value: "authorization"}},
							},
						},
						Examples: []*ast.Example{
							{
								Mimetype: &ast.Attribute{Value: xmlenc.String{Value: "application/json"}},
								Content:  &ast.ExampleValue{Value: xmlenc.String{Value: "xxx"}},
							},
						},
						Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
						Items: []*ast.Param{
							{
								Type:    &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
								Name:    &ast.Attribute{Value: xmlenc.String{Value: "id"}},
								Summary: &ast.Attribute{Value: xmlenc.String{Value: "ID"}},
							},
							{
								Type:    &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
								Name:    &ast.Attribute{Value: xmlenc.String{Value: "name"}},
								Summary: &ast.Attribute{Value: xmlenc.String{Value: "name summary"}},
							},
						},
					},
				},
				Responses: []*ast.Request{
					{
						Status: &ast.StatusAttribute{Value: ast.Number{Int: http.StatusCreated}},
						Description: &ast.Richtext{
							Type: &ast.Attribute{Value: xmlenc.String{Value: "html"}},
							Text: &ast.CData{Value: xmlenc.String{Value: "<p>desc</p>"}},
						},
						Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNone}},
					},
				},
			},
		},
	}
}

// XML 获取 Get 返回对象的 XML 编码
func XML(a *assert.Assertion) []byte {
	data, err := xmlenc.Encode("", Get(), core.XMLNamespace, "apidoc")
	a.NotError(err).NotNil(data)

	return data
}

// URI 返回测试文件基于 URI 的表示方式
func URI(a *assert.Assertion) core.URI {
	p := core.FileURI(pp(a, Filename))
	a.NotEmpty(p)
	return p
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
	p = source.CurrentPath(p)
	p, err := filepath.Abs(p)
	a.NotError(err).NotEmpty(p)
	return p
}
