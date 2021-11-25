// SPDX-License-Identifier: MIT

package apidoc

import (
	"log"
	"net/http"
	"testing"

	"github.com/issue9/assert/v2"
	"github.com/issue9/assert/v2/rest"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/v7/internal/ast/asttest"
	"github.com/caixw/apidoc/v7/internal/docs"
)

func TestVersion(t *testing.T) {
	a := assert.New(t, false)

	a.True(version.SemVerValid(Version(true)))
	a.True(version.SemVerValid(Version(false)))
	a.True(version.SemVerValid(DocVersion))
	a.True(version.SemVerValid(LSPVersion))
}

func TestStatic(t *testing.T) {
	a := assert.New(t, false)
	srv := rest.NewServer(a, Static(docs.Dir(), false, log.Default()), nil)

	srv.Get("/icon.svg").Do(nil).Status(http.StatusOK)
}

func TestView_Buffer(t *testing.T) {
	a := assert.New(t, false)
	data := asttest.XML(a)

	s := &Server{
		Status:      http.StatusCreated,
		Path:        "/test/apidoc.xml",
		ContentType: "text/xml1",
	}
	srv := rest.NewServer(a, s.Buffer(data), nil)
	srv.Get("/test/apidoc.xml").Do(nil).
		Status(http.StatusCreated).
		Header("content-type", "text/xml1")

	srv.Get("/index.xml").Do(nil).Status(http.StatusOK)

	srv.Get("/v6/apidoc.xsl").Do(nil).Status(http.StatusOK)

	// 能正确覆盖 Static 中的 index.xml
	s = &Server{
		Status:      http.StatusCreated,
		Path:        "/index.xml",
		ContentType: "text/css",
	}
	srv = rest.NewServer(a, s.Buffer(data), nil)
	srv.Get("/index.xml").Do(nil).
		Status(http.StatusCreated).
		Header("content-type", "text/css")

	srv.Get("/v6/apidoc.xsl").Do(nil).Status(http.StatusOK)
}

func TestView_File(t *testing.T) {
	a := assert.New(t, false)

	s := &Server{
		Status:      http.StatusAccepted,
		Path:        "/apidoc.xml",
		ContentType: "text/xml",
	}
	h, err := s.File(asttest.URI(a))
	a.NotError(err).NotNil(h)
	srv := rest.NewServer(a, h, nil)
	srv.Get("/apidoc.xml").Do(nil).
		Status(http.StatusAccepted).
		Header("content-type", "text/xml")

	s = &Server{
		Status: http.StatusAccepted,
		Path:   "/apidoc.xml",
	}
	h, err = s.File(asttest.URI(a))
	a.NotError(err).NotNil(h)
	srv = rest.NewServer(a, h, nil)
	srv.Get("/apidoc.xml").Do(nil).
		Status(http.StatusAccepted)

	// 覆盖现有的 index.xml
	s = &Server{
		Status: http.StatusAccepted,
	}
	h, err = s.File(asttest.URI(a))
	a.NotError(err).NotNil(h)
	srv = rest.NewServer(a, h, nil)
	srv.Get("/index.xml").Do(nil).
		Status(http.StatusAccepted)
}

func TestAddStylesheet(t *testing.T) {
	a := assert.New(t, false)

	data := []*struct {
		input  string
		output string
	}{
		{
			input: "",
			output: `
<?xml-stylesheet type="text/xsl" href="./v6/apidoc.xsl"?>`,
		},
		{
			input: `<?xml version="1.0"?>`,
			output: `<?xml version="1.0"?>
<?xml-stylesheet type="text/xsl" href="./v6/apidoc.xsl"?>`,
		},
		{
			input: `<?xml version="1.0"?>
<?xml-stylesheet href="xxx"?>`,
			output: `<?xml version="1.0"?>
<?xml-stylesheet type="text/xsl" href="./v6/apidoc.xsl"?>
<?xml-stylesheet href="xxx"?>`,
		},
	}

	for index, item := range data {
		output := string(addStylesheet([]byte(item.input)))
		a.Equal(output, item.output, "not equal at %d\nv1: %s\nv2:%s", index, item.output, output)
	}
}
