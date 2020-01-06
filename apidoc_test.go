// SPDX-License-Identifier: MIT

package apidoc

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/v5/doc/doctest"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message/messagetest"
)

func TestVersion(t *testing.T) {
	a := assert.New(t)

	a.Equal(Version(), vars.Version())
	a.True(version.SemVerValid(Version()))
}

func TestValid(t *testing.T) {
	a := assert.New(t)

	data, err := ioutil.ReadFile("./docs/example/index.xml")
	a.NotError(err).NotNil(data)
	a.NotError(Valid(data))
}

func TestStatic(t *testing.T) {
	srv := rest.NewServer(t, Static(vars.DocsDir(), false), nil)
	defer srv.Close()

	srv.Get("/icon.svg").Do().Status(http.StatusOK)
}

func TestView(t *testing.T) {
	a := assert.New(t)

	data := doctest.XML(a)
	h := View(http.StatusCreated, "/test/apidoc.xml", data, "text/xml", "", false)
	srv := rest.NewServer(t, h, nil)
	srv.Get("/test/apidoc.xml").Do().
		Status(http.StatusCreated).
		Header("content-type", "text/xml")

	srv.Get("/index.xml").Do().
		Status(http.StatusOK)

	srv.Get("/v5/apidoc.xsl").Do().
		Status(http.StatusOK)

	srv.Close()

	// 能正确覆盖 Static 中的 index.xml
	h = View(http.StatusCreated, "/index.xml", data, "text/css", "", false)
	srv = rest.NewServer(t, h, nil)
	srv.Get("/index.xml").Do().
		Status(http.StatusCreated).
		Header("content-type", "text/css")

	srv.Get("/v5/apidoc.xsl").Do().
		Status(http.StatusOK)

	srv.Close()
}

func TestViewFile(t *testing.T) {
	a := assert.New(t)

	h, err := ViewFile(http.StatusAccepted, "/apidoc.xml", doctest.Path(a), "text/xml", "", false)
	a.NotError(err).NotNil(h)
	srv := rest.NewServer(t, h, nil)
	srv.Get("/apidoc.xml").Do().
		Status(http.StatusAccepted).
		Header("content-type", "text/xml")
	srv.Close()

	h, err = ViewFile(http.StatusAccepted, "/apidoc.xml", doctest.Path(a), "", "", false)
	a.NotError(err).NotNil(h)
	srv = rest.NewServer(t, h, nil)
	srv.Get("/apidoc.xml").Do().
		Status(http.StatusAccepted)
	srv.Close()

	h, err = ViewFile(http.StatusAccepted, "", doctest.Path(a), "", "", false)
	a.NotError(err).NotNil(h)
	srv = rest.NewServer(t, h, nil)
	srv.Get("/index.xml").Do().
		Status(http.StatusAccepted)
	srv.Close()
}

func TestMockFile(t *testing.T) {
	a := assert.New(t)

	_, _, h := messagetest.MessageHandler()
	mock, err := MockFile(h, doctest.Path(a), map[string]string{"admin": "/admin"})
	a.NotError(err).NotNil(h)

	srv := rest.NewServer(t, mock, nil)
	defer srv.Close()

	srv.Get("/admin/users").
		Header("authorization", "xxx").
		Header("content-type", "application/json").
		Header("Accept", "application/json").
		Do().Status(http.StatusOK)
	srv.Post("/admin/users", nil).Do().Status(http.StatusBadRequest)    // 未指定报头
	srv.Delete("/admin/users").Do().Status(http.StatusMethodNotAllowed) // 不存在
}
