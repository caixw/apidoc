// SPDX-License-Identifier: MIT

package openapi

import (
	"testing"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/issue9/assert"
)

func TestNewServer(t *testing.T) {
	a := assert.New(t)

	input := &doc.Server{
		URL:     "https://example.com",
		Name:    "name",
		Summary: "summary",
	}

	output := newServer(input)
	a.NotNil(output).
		Equal(output.URL, "https://example.com").
		Equal(output.Description, "summary")

	input.Description = doc.Richtext{Text: "desc"}
	output = newServer(input)
	a.NotNil(output).
		Equal(output.URL, "https://example.com").
		Equal(output.Description, "desc")
}

func TestServer_sanitize(t *testing.T) {
	a := assert.New(t)

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
	a := assert.New(t)

	sv := &ServerVariable{}
	a.Error(sv.sanitize())

	sv.Enum = []string{"e1", "e2"}
	a.Error(sv.sanitize())

	sv.Default = "not-in-enum"
	a.Error(sv.sanitize())

	sv.Default = "e1"
	a.NotError(sv.sanitize())
}
