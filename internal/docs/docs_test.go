// SPDX-License-Identifier: MIT

package docs

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"
)

func TestDir(t *testing.T) {
	a := assert.New(t)

	p1, err := filepath.Abs("../../docs")
	a.NotError(err).NotEmpty(p1)

	p2, err := filepath.Abs(Dir())
	a.NotError(err).NotEmpty(p2)

	a.Equal(p1, p2)
}

func TestPath(t *testing.T) {
	a := assert.New(t)

	p1, err := filepath.Abs("../../docs/example")
	a.NotError(err).NotEmpty(p1)

	p2, err := filepath.Abs(Path("example"))
	a.NotError(err).NotEmpty(p2)

	a.Equal(p1, p2)
}

func TestEmbeddedHandler(t *testing.T) {
	a := assert.New(t)

	srv := rest.NewServer(t, Handler("", false), nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/v5/").
		Do().
		Status(http.StatusNotFound)

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

	srv.Get("/v5/").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/v5/apidoc.xsl").
		Do().
		Status(http.StatusOK)

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

	srv.Get("/prefix/v5/").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/prefix/v5/apidoc.xsl").
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

func TestFolderHandler(t *testing.T) {
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

func TestFolderHandler_stylesheet(t *testing.T) {
	a := assert.New(t)

	srv := rest.NewServer(t, Handler(Dir(), true), nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/icon.svg").
		Do().
		Status(http.StatusOK)

	srv.Get("/v5/apidoc.xsl").
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

func TestFolderHandler_prefix(t *testing.T) {
	a := assert.New(t)

	h := http.StripPrefix("/prefix/", folderHandler(Dir(), false))
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

	srv.Get("/prefix/v5/apidoc.xsl").
		Do().
		Status(http.StatusOK)
}
