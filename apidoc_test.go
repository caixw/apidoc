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
	srv := rest.NewServer(t, Static("./docs", false), nil)
	defer srv.Close()

	srv.Get("/icon.svg").Do().Status(http.StatusOK)
}

func TestView(t *testing.T) {
	a := assert.New(t)

	data, err := doctest.XML()
	a.NotError(err).NotNil(data)

	h := View(http.StatusCreated, "/apidoc.xml", data, "text/xml", "", false)
	srv := rest.NewServer(t, h, nil)
	defer srv.Close()
	srv.Get("/apidoc.xml").Do().
		Status(http.StatusCreated).
		Header("content-type", "text/xml")

	srv.Get("/index.xml").Do().
		Status(http.StatusOK)

	srv.Get("/v5/apidoc.xsl").Do().
		Status(http.StatusOK)
}
