// SPDX-License-Identifier: MIT

package docs

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
)

func TestDir(t *testing.T) {
	a := assert.New(t)

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
	a := assert.New(t)

	a.Equal(StylesheetURL(""), ast.MajorVersion+"/apidoc.xsl")
	a.Equal(StylesheetURL("."), "./"+ast.MajorVersion+"/apidoc.xsl")
	a.Equal(StylesheetURL("./"), "./"+ast.MajorVersion+"/apidoc.xsl")
	a.Equal(StylesheetURL("https://apidoc.tools/"), "https://apidoc.tools/"+ast.MajorVersion+"/apidoc.xsl")
	a.Equal(StylesheetURL("https://apidoc.tools"), "https://apidoc.tools/"+ast.MajorVersion+"/apidoc.xsl")
}

func TestEmbeddedHandler(t *testing.T) {
	a := assert.New(t)

	srv := rest.NewServer(t, Handler("", false), nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/v6/").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/v6/apidoc.xsl").
		Do().
		Status(http.StatusOK).
		Header("Content-Type", "application/xml")

	srv.Get("/v5/apidoc.xsl").
		Do().
		Status(http.StatusOK)

	srv.Get("/example").
		Do().
		Status(http.StatusOK)

	srv.Get("/").
		Do().
		Status(http.StatusOK)

	srv.Get("/icon.svg").
		Do().
		Status(http.StatusOK)
}

func TestEmbeddedHandler_stylesheet(t *testing.T) {
	a := assert.New(t)

	srv := rest.NewServer(t, Handler("", true), nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/v6/").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/v6/apidoc.xsl").
		Do().
		Status(http.StatusOK)

	srv.Get("/v5/apidoc.xsl").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/example").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/icon.svg").
		Do().
		Status(http.StatusOK)
}

func TestEmbeddedHandler_prefix(t *testing.T) {
	a := assert.New(t)

	h := http.StripPrefix("/prefix/", Handler("", false))
	srv := rest.NewServer(t, h, nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/prefix/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/prefix/v6/").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/prefix/v6/apidoc.xsl").
		Do().
		Status(http.StatusOK)

	srv.Get("/prefix/example").
		Do().
		Status(http.StatusOK)

	srv.Get("/prefix/").
		Do().
		Status(http.StatusOK)

	srv.Get("/prefix/icon.svg").
		Do().
		Status(http.StatusOK)
}

func TestLocalHandler(t *testing.T) {
	a := assert.New(t)

	srv := rest.NewServer(t, Handler(Dir(), false), nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/").
		Do().
		Status(http.StatusOK)

	srv.Get("/example/").
		Do().
		Status(http.StatusOK)

	srv.Get("/example").
		Do().
		Status(http.StatusOK)

	srv.Get("/index.xml").
		Do().
		Status(http.StatusOK)
}

func TestLocalHandler_stylesheet(t *testing.T) {
	a := assert.New(t)

	srv := rest.NewServer(t, Handler(Dir(), true), nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/icon.svg").
		Do().
		Status(http.StatusOK)

	srv.Get("/v6/apidoc.xsl").
		Do().
		Status(http.StatusOK)

	srv.Get("/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/example/").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/example").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/index.xml").
		Do().
		Status(http.StatusNotFound)
}

func TestLocalHandler_prefix(t *testing.T) {
	a := assert.New(t)

	h := http.StripPrefix("/prefix/", Handler(Dir(), false))
	srv := rest.NewServer(t, h, nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/prefix/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/prefix/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/prefix/").
		Do().
		Status(http.StatusOK)

	srv.Get("/prefix/index.xml").
		Do().
		Status(http.StatusOK)

	srv.Get("/prefix/v6/apidoc.xsl").
		Do().
		Status(http.StatusOK)
}

func TestRemoteHandler(t *testing.T) {
	a := assert.New(t)

	remote := httptest.NewServer(Handler(Dir(), false))
	srv := rest.NewServer(t, Handler(core.URI(remote.URL), false), nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/").
		Do().
		Status(http.StatusOK)

	srv.Get("/example/").
		Do().
		Status(http.StatusOK)

	srv.Get("/example").
		Do().
		Status(http.StatusOK)

	srv.Get("/index.xml").
		Do().
		Status(http.StatusOK)
}

func TestRemoteHandler_stylesheet(t *testing.T) {
	a := assert.New(t)

	remote := httptest.NewServer(Handler(Dir(), false))
	srv := rest.NewServer(t, Handler(core.URI(remote.URL), true), nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/icon.svg").
		Do().
		Status(http.StatusOK)

	srv.Get("/v6/apidoc.xsl").
		Do().
		Status(http.StatusOK)

	srv.Get("/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/example/").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/example").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/index.xml").
		Do().
		Status(http.StatusNotFound)
}

func TestRemoteHandler_prefix(t *testing.T) {
	a := assert.New(t)

	remote := httptest.NewServer(Handler(Dir(), false))
	h := http.StripPrefix("/prefix/", Handler(core.URI(remote.URL), false))
	srv := rest.NewServer(t, h, nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/prefix/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/prefix/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/prefix/").
		Do().
		Status(http.StatusOK)

	srv.Get("/prefix/index.xml").
		Do().
		Status(http.StatusOK)

	srv.Get("/prefix/v6/apidoc.xsl").
		Do().
		Status(http.StatusOK)
}
