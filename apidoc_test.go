// SPDX-License-Identifier: MIT

package apidoc

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/v5/internal/vars"
)

func TestVersion(t *testing.T) {
	a := assert.New(t)

	a.Equal(Version(), vars.Version())
	a.True(version.SemVerValid(Version()))
}

func TestHandler(t *testing.T) {
	a := assert.New(t)

	srv := rest.NewServer(t, Handler("./docs/example/index.xml", "", "", nil), nil)
	defer srv.Close()
	srv.Get("").
		Do().
		Status(http.StatusOK).
		Header("content-type", docContentType)

	// 测试 content-type 是否正确
	srv = rest.NewServer(t, Handler("./docs/example/index.xml", "text/xml", "", nil), nil)
	defer srv.Close()
	srv.Get("").
		Do().
		Status(http.StatusOK).
		Header("content-type", "text/xml")

	// 测试 xml-stylesheet 是否正确
	xsl := "https://example.com/xsl.xsl"
	buf := new(bytes.Buffer)
	srv = rest.NewServer(t, Handler("./docs/example/index.xml", "", xsl, nil), nil)
	defer srv.Close()
	srv.Get("").
		Do().
		Status(http.StatusOK).
		ReadBody(buf)
	a.Contains(buf.String(), xsl)
}
