// SPDX-License-Identifier: MIT

package docs

import (
	"net/http"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"
)

func TestHandler(t *testing.T) {
	a := assert.New(t)

	srv := rest.NewServer(t, Handler("../../docs"), nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/CNAME").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/").
		Do().
		Status(http.StatusOK).
		Header("content-type", "application/xml")

	srv.Get("/example/").
		Do().
		Status(http.StatusOK).
		Header("content-type", "application/xml")

	srv.Get("/example").
		Do().
		Status(http.StatusOK).
		Header("content-type", "application/xml")

	srv.Get("/index.xml").
		Do().
		Status(http.StatusOK).
		Header("content-type", "application/xml")
}

func TestHandler_prefix(t *testing.T) {
	a := assert.New(t)

	h := http.StripPrefix("/prefix/", Handler("../../docs"))
	srv := rest.NewServer(t, h, nil)
	a.NotNil(srv)
	defer srv.Close()

	srv.Get("/prefix/CNAME").
		Do().
		Status(http.StatusNotFound)

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
		Status(http.StatusOK).
		Header("content-type", "application/xml")

	srv.Get("/prefix/index.xml").
		Do().
		Status(http.StatusOK).
		Header("content-type", "application/xml")

	srv.Get("/prefix/v5/apidoc.xsl").
		Do().
		Status(http.StatusOK).
		Header("content-type", "application/xml")
}
