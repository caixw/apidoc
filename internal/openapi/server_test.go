// SPDX-License-Identifier: MIT

package openapi

import (
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

func TestNewServer(t *testing.T) {
	a := assert.New(t, false)

	input := &ast.Server{
		URL:     &ast.Attribute{Value: xmlenc.String{Value: "https://example.com"}},
		Name:    &ast.Attribute{Value: xmlenc.String{Value: "name"}},
		Summary: &ast.Attribute{Value: xmlenc.String{Value: "summary"}},
	}

	output := newServer(input)
	a.NotNil(output).
		Equal(output.URL, "https://example.com").
		Equal(output.Description, "summary")

	input.Description = &ast.Richtext{Text: &ast.CData{Value: xmlenc.String{Value: "desc"}}}
	output = newServer(input)
	a.NotNil(output).
		Equal(output.URL, "https://example.com").
		Equal(output.Description, "desc")
}

func TestServer_sanitize(t *testing.T) {
	a := assert.New(t, false)

	srv := &Server{}
	a.Error(srv.sanitize())

	srv.URL = "https://example.com/{tpl1}/{tpl2}/path3"
	a.NotError(srv.sanitize())

	srv.Variables = map[string]*ServerVariable{
		"tpl1": {Default: "1"},
		"tpl2": {Default: "2", Enum: []string{"1", "2"}},
	}
	a.NotError(srv.sanitize())

	// variable 不在 URL 中
	srv.Variables = map[string]*ServerVariable{
		"tpl3": {Default: "1"},
	}
	a.Error(srv.sanitize())

	// variables 存在错误
	srv.Variables = map[string]*ServerVariable{
		"tpl2": {Default: "not-exists", Enum: []string{"1", "2"}},
	}
	a.Error(srv.sanitize())
}

func TestServerVariable_sanitize(t *testing.T) {
	a := assert.New(t, false)

	sv := &ServerVariable{}
	a.Error(sv.sanitize())

	sv.Enum = []string{"e1", "e2"}
	a.Error(sv.sanitize())

	sv.Default = "not-in-enum"
	a.Error(sv.sanitize())

	sv.Default = "e1"
	a.NotError(sv.sanitize())
}
