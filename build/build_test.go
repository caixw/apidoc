// SPDX-License-Identifier: MIT

package build

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/ast"
)

func TestParse(t *testing.T) {
	a := assert.New(t)

	php := &Input{
		Lang:      "php",
		Dir:       "./testdata",
		Recursive: true,
		Encoding:  "gbk",
	}
	a.NotError(php.Sanitize())

	c := &Input{
		Lang:      "c++",
		Dir:       "./testdata",
		Recursive: true,
	}
	a.NotError(c.Sanitize())

	doc := &ast.APIDoc{}
	erro, _, h := messagetest.MessageHandler()
	Parse(doc, h, php, c)
	h.Stop()
	a.Empty(erro.String())

	a.Equal(2, len(doc.Apis)).
		Equal(doc.Version.V(), "1.1.1")
	api := doc.Apis[0]
	a.Equal(api.Method.V(), "GET")
}

func TestParseFile(t *testing.T) {
	a := assert.New(t)

	c := &Input{
		Lang:      "c++",
		Dir:       "./testdata",
		Recursive: true,
	}
	a.NotError(c.Sanitize())

	doc := &ast.APIDoc{}
	erro, _, h := messagetest.MessageHandler()
	uri, err := core.FileURI("./testdata/testfile.h")
	a.NotError(err).NotEmpty(uri)
	ParseFile(doc, h, uri, c)
	a.Equal(0, len(doc.Apis)).
		Equal(doc.Version.V(), "1.1.1")
	h.Stop()
	a.Empty(erro.String())
}
