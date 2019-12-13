// SPDX-License-Identifier: MIT

package static

import (
	"net/http"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"
)

var embedData = []*FileInfo{
	{
		Name:        "icon.svg",
		ContentType: "image/svg+xml",
		Content:     []byte{' ', ' '},
	},
	{
		Name:        "index.xml",
		ContentType: "image/svg+xml",
		Content:     []byte{' ', ' '},
	},

	{
		Name:        "example/index.xml",
		ContentType: "image/svg+xml",
		Content:     []byte{' ', ' '},
	},
	{
		Name:        "example/test.xml",
		ContentType: "image/svg+xml",
		Content:     []byte{' ', ' '},
	},

	{
		Name:        "v5/index.xml",
		ContentType: "image/svg+xml",
		Content:     []byte{' ', ' '},
	},
	{
		Name:        "v5/apidoc.xsl",
		ContentType: "image/svg+xml",
		Content:     []byte{' ', ' '},
	},
}

func TestEmbeddedHandler(t *testing.T) {
	a := assert.New(t)

	srv := rest.NewServer(t, EmbeddedHandler(embedData, false), nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/v5/").
		Do().
		Status(http.StatusOK)

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

	srv := rest.NewServer(t, EmbeddedHandler(embedData, true), nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/v5/").
		Do().
		Status(http.StatusOK)

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

	h := http.StripPrefix("/prefix/", EmbeddedHandler(embedData, true))
	srv := rest.NewServer(t, h, nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/prefix/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/prefix/v5/").
		Do().
		Status(http.StatusOK)

	srv.Get("/prefix/v5/apidoc.xsl").
		Do().
		Status(http.StatusOK)

	srv.Get("/prefix/example").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/prefix/").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/prefix/icon.svg").
		Do().
		Status(http.StatusOK)
}

func TestFolderHandler(t *testing.T) {
	a := assert.New(t)

	srv := rest.NewServer(t, FolderHandler("../docs", false), nil)
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

	srv := rest.NewServer(t, FolderHandler("../docs", true), nil)
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

	h := http.StripPrefix("/prefix/", FolderHandler("../docs", false))
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
