// SPDX-License-Identifier: MIT

package static

import (
	"net/http"
	"strings"
	"testing"

	"github.com/caixw/apidoc/v5/internal/vars"
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

// 保证 styles 中保存着最新的 xml-stylesheet 内容
func TestStyles(t *testing.T) {
	a := assert.New(t)

	v := vars.DocVersion()
	found := false
	for _, file := range styles {
		if strings.HasPrefix(file, v) {
			found = true
		}
	}
	a.True(found)
}

func TestEmbeddedHandler(t *testing.T) {
	a := assert.New(t)

	srv := rest.NewServer(t, EmbeddedHandler(embedData), nil)
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

func TestEmbeddedHandler_prefix(t *testing.T) {
	a := assert.New(t)

	h := http.StripPrefix("/prefix/", EmbeddedHandler(embedData))
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

	srv := rest.NewServer(t, FolderHandler("../docs", TypeAll), nil)
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

func TestFolderHandler_TypeStylesheet(t *testing.T) {
	a := assert.New(t)

	srv := rest.NewServer(t, FolderHandler("../docs", TypeStylesheet), nil)
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

func TestFolderHandler_TypeNone(t *testing.T) {
	a := assert.New(t)

	srv := rest.NewServer(t, FolderHandler("../docs", TypeNone), nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/icon.svg").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/v5/apidoc.xsl").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/").
		Do().
		Status(http.StatusNotFound)
}

func TestFolderHandler_prefix(t *testing.T) {
	a := assert.New(t)

	h := http.StripPrefix("/prefix/", FolderHandler("../docs", TypeAll))
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
