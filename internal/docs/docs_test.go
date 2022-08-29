// SPDX-License-Identifier: MIT

package docs

import (
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/issue9/assert/v3"
	"github.com/issue9/assert/v3/rest"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
)

func TestDir(t *testing.T) {
	a := assert.New(t, false)

	abs, err := filepath.Abs("../../docs")
	a.NotError(err).NotEmpty(abs)
	p1 := core.FileURI(abs)
	a.NotEmpty(p1)

	p2 := Dir()
	a.NotEmpty(p2)
	a.Equal(p1, p2)

	exists, err := Dir().Exists()
	a.NotError(err).True(exists)
}

func TestStylesheetURL(t *testing.T) {
	a := assert.New(t, false)

	a.Equal(StylesheetURL(""), ast.MajorVersion+"/apidoc.xsl")
	a.Equal(StylesheetURL("."), "./"+ast.MajorVersion+"/apidoc.xsl")
	a.Equal(StylesheetURL("./"), "./"+ast.MajorVersion+"/apidoc.xsl")
	a.Equal(StylesheetURL("https://apidoc.tools/"), "https://apidoc.tools/"+ast.MajorVersion+"/apidoc.xsl")
	a.Equal(StylesheetURL("https://apidoc.tools"), "https://apidoc.tools/"+ast.MajorVersion+"/apidoc.xsl")
}

func TestEmbeddedHandler(t *testing.T) {
	a := assert.New(t, false)

	srv := rest.NewServer(a, Handler("", false, log.Default()), nil)
	a.NotNil(srv)

	srv.Get("/not-exists").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/v6/").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/v6/apidoc.xsl").
		Do(nil).
		Status(http.StatusOK)
		//Header("Content-Type", "application/xml") // NOTE: Content-Type 有系统决定，无法通过所有系统的测试

	srv.Get("/v5/apidoc.xsl").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/example").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/icon.svg").
		Do(nil).
		Status(http.StatusOK)
}

func TestEmbeddedHandler_stylesheet(t *testing.T) {
	a := assert.New(t, false)

	srv := rest.NewServer(a, Handler("", true, log.Default()), nil)
	a.NotNil(srv)

	srv.Get("/not-exists").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/v6/").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/v6/apidoc.xsl").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/v5/apidoc.xsl").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/example").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/icon.svg").
		Do(nil).
		Status(http.StatusOK)
}

func TestEmbeddedHandler_prefix(t *testing.T) {
	a := assert.New(t, false)

	h := http.StripPrefix("/prefix/", Handler("", false, log.Default()))
	srv := rest.NewServer(a, h, nil)
	a.NotNil(srv)

	srv.Get("/prefix/not-exists").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/prefix/v6/").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/prefix/v6/apidoc.xsl").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/prefix/example").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/prefix/").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/prefix/icon.svg").
		Do(nil).
		Status(http.StatusOK)
}

func TestLocalHandler(t *testing.T) {
	a := assert.New(t, false)

	srv := rest.NewServer(a, Handler(Dir(), false, log.Default()), nil)
	a.NotNil(srv)

	srv.Get("/not-exists").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/example/").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/example").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/index.xml").
		Do(nil).
		Status(http.StatusOK)
}

func TestLocalHandler_stylesheet(t *testing.T) {
	a := assert.New(t, false)

	srv := rest.NewServer(a, Handler(Dir(), true, log.Default()), nil)
	a.NotNil(srv)

	srv.Get("/icon.svg").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/v6/apidoc.xsl").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/not-exists").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/example/").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/example").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/index.xml").
		Do(nil).
		Status(http.StatusNotFound)
}

func TestLocalHandler_prefix(t *testing.T) {
	a := assert.New(t, false)

	h := http.StripPrefix("/prefix/", Handler(Dir(), false, log.Default()))
	srv := rest.NewServer(a, h, nil)
	a.NotNil(srv)

	srv.Get("/prefix/not-exists").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/prefix/not-exists").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/prefix/").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/prefix/index.xml").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/prefix/v6/apidoc.xsl").
		Do(nil).
		Status(http.StatusOK)
}

func TestRemoteHandler(t *testing.T) {
	a := assert.New(t, false)

	remote := httptest.NewServer(Handler(Dir(), false, log.Default()))
	srv := rest.NewServer(a, Handler(core.URI(remote.URL), false, log.Default()), nil)
	a.NotNil(srv)

	srv.Get("/not-exists").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/example/").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/example").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/index.xml").
		Do(nil).
		Status(http.StatusOK)
}

func TestRemoteHandler_stylesheet(t *testing.T) {
	a := assert.New(t, false)

	remote := httptest.NewServer(Handler(Dir(), false, log.Default()))
	srv := rest.NewServer(a, Handler(core.URI(remote.URL), true, log.Default()), nil)
	a.NotNil(srv)

	srv.Get("/icon.svg").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/v6/apidoc.xsl").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/not-exists").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/example/").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/example").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/index.xml").
		Do(nil).
		Status(http.StatusNotFound)
}

func TestRemoteHandler_prefix(t *testing.T) {
	a := assert.New(t, false)

	remote := httptest.NewServer(Handler(Dir(), false, log.Default()))
	h := http.StripPrefix("/prefix/", Handler(core.URI(remote.URL), false, log.Default()))
	srv := rest.NewServer(a, h, nil)
	a.NotNil(srv)

	srv.Get("/prefix/not-exists").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/prefix/not-exists").
		Do(nil).
		Status(http.StatusNotFound)

	srv.Get("/prefix/").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/prefix/index.xml").
		Do(nil).
		Status(http.StatusOK)

	srv.Get("/prefix/v6/apidoc.xsl").
		Do(nil).
		Status(http.StatusOK)
}
