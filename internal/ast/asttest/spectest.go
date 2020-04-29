// SPDX-License-Identifier: MIT

// Package asttest 提供了一个合法的 ast.APIDoc 对象
package asttest

import (
	"net/http"
	"path/filepath"

	"github.com/issue9/assert"
	"github.com/issue9/utils"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/token"
)

// Filename 文档的文件名
const Filename = "index.xml"

// Get 返回 doc.APIDoc 对象
//
// 同时当前目录下的 index.xml 文件与此返回对象内容是相同的。
func Get() *ast.APIDoc {
	return &ast.APIDoc{
		APIDoc:  &ast.APIDocVersionAttribute{Value: ast.String{Value: ast.Version}},
		Version: &ast.VersionAttribute{Value: ast.String{Value: "1.0.1"}},
		Title:   &ast.Element{Content: ast.String{Value: "test"}},
		Description: &ast.Richtext{
			Text: &ast.CData{Value: ast.String{Value: "<p>desc</p>"}},
			Type: &ast.Attribute{Value: ast.String{Value: ast.RichtextTypeHTML}},
		},
		Servers: []*ast.Server{
			{
				URL:     &ast.Attribute{Value: ast.String{Value: "https://example.com/admin"}},
				Name:    &ast.Attribute{Value: ast.String{Value: "admin"}},
				Summary: &ast.Attribute{Value: ast.String{Value: "admin"}},
			},
			{
				URL:     &ast.Attribute{Value: ast.String{Value: "https://example.com"}},
				Name:    &ast.Attribute{Value: ast.String{Value: "client"}},
				Summary: &ast.Attribute{Value: ast.String{Value: "client"}},
			},
		},
		Tags: []*ast.Tag{
			{
				Name:  &ast.Attribute{Value: ast.String{Value: "t1"}},
				Title: &ast.Attribute{Value: ast.String{Value: "t1"}},
			},
			{
				Name:  &ast.Attribute{Value: ast.String{Value: "t2"}},
				Title: &ast.Attribute{Value: ast.String{Value: "t2"}},
			},
			{
				Name:  &ast.Attribute{Value: ast.String{Value: "tag1"}},
				Title: &ast.Attribute{Value: ast.String{Value: "tag1"}},
			},
		},
		Mimetypes: []*ast.Element{
			{Content: ast.String{Value: "application/json"}},
			{Content: ast.String{Value: "application/xml"}},
		},
		Apis: []*ast.API{
			{
				Method: &ast.MethodAttribute{Value: ast.String{Value: http.MethodGet}},
				Tags: []*ast.Element{
					{Content: ast.String{Value: "t1"}},
					{Content: ast.String{Value: "t2"}},
				},
				Path: &ast.Path{Path: &ast.Attribute{Value: ast.String{Value: "/users"}}},
				Servers: []*ast.Element{
					{Content: ast.String{Value: "admin"}},
					{Content: ast.String{Value: "client"}},
				},
				Requests: []*ast.Request{
					{
						Summary: &ast.Attribute{Value: ast.String{Value: "request"}},
						Headers: []*ast.Param{
							{
								Type:    &ast.TypeAttribute{Value: ast.String{Value: ast.TypeString}},
								Name:    &ast.Attribute{Value: ast.String{Value: "authorization"}},
								Summary: &ast.Attribute{Value: ast.String{Value: "authorization"}},
							},
						},
						Examples: []*ast.Example{
							{
								Mimetype: &ast.Attribute{Value: ast.String{Value: "application/json"}},
								Content:  &ast.CData{Value: ast.String{Value: "xxx"}},
							},
						},
					},
				},
				Responses: []*ast.Request{
					{
						Description: &ast.Richtext{
							Type: &ast.Attribute{Value: ast.String{Value: "html"}},
							Text: &ast.CData{Value: ast.String{Value: "<p>desc</p>"}},
						},
						Type:   &ast.TypeAttribute{Value: ast.String{Value: ast.TypeObject}},
						Status: &ast.StatusAttribute{Value: ast.Number{Value: http.StatusOK}},
						Headers: []*ast.Param{
							{
								Type:    &ast.TypeAttribute{Value: ast.String{Value: ast.TypeString}},
								Name:    &ast.Attribute{Value: ast.String{Value: "authorization"}},
								Summary: &ast.Attribute{Value: ast.String{Value: "authorization"}},
							},
						},
						Examples: []*ast.Example{
							{
								Mimetype: &ast.Attribute{Value: ast.String{Value: "application/json"}},
								Content:  &ast.CData{Value: ast.String{Value: "xxx"}},
							},
						},
						Items: []*ast.Param{
							{
								Summary: &ast.Attribute{Value: ast.String{Value: "summary"}},
								Type:    &ast.TypeAttribute{Value: ast.String{Value: ast.TypeString}},
								Name:    &ast.Attribute{Value: ast.String{Value: "name"}},
							},
						},
					},
				},
			},
			{
				Method: &ast.MethodAttribute{Value: ast.String{Value: http.MethodPost}},
				Tags: []*ast.Element{
					{Content: ast.String{Value: "t1"}},
					{Content: ast.String{Value: "tag1"}},
				},
				Path:       &ast.Path{Path: &ast.Attribute{Value: ast.String{Value: "/users"}}},
				Deprecated: &ast.VersionAttribute{Value: ast.String{Value: "1.0.1"}},
				Summary:    &ast.Attribute{Value: ast.String{Value: "summary"}},
				Servers: []*ast.Element{
					{Content: ast.String{Value: "admin"}},
				},
				Responses: []*ast.Request{
					{
						Description: &ast.Richtext{
							Type: &ast.Attribute{Value: ast.String{Value: "html"}},
							Text: &ast.CData{Value: ast.String{Value: "<p>desc</p>"}},
						},
						Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNone}},
					},
				},
			},
		},
	}
}

// XML 获取 Get 返回对象的 XML 编码
func XML(a *assert.Assertion) []byte {
	data, err := token.Encode("", Get())
	a.NotError(err).NotNil(data)

	return data
}

// URI 返回测试文件基于 URI 的表示方式
func URI(a *assert.Assertion) core.URI {
	p, err := core.FileURI(pp(a, Filename))
	a.NotError(err).NotEmpty(p)
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
	p = utils.CurrentPath(p)
	p, err := filepath.Abs(p)
	a.NotError(err).NotEmpty(p)
	return p
}
