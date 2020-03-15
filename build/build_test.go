// SPDX-License-Identifier: MIT

package build

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/message/messagetest"
	"github.com/caixw/apidoc/v6/spec"
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

	doc := spec.NewAPIDoc()
	a.NotNil(doc)
	erro, _, h := messagetest.MessageHandler()
	Parse(doc, h, php, c)
	h.Stop()
	a.Empty(erro.String())

	a.NotError(doc.Sanitize())
	a.Equal(2, len(doc.Apis)).
		Equal(doc.Version, "1.1.1")
	api := doc.Apis[0]
	a.Equal(api.Method, "GET")
}

func TestParseFile(t *testing.T) {
	a := assert.New(t)

	c := &Input{
		Lang:      "c++",
		Dir:       "./testdata",
		Recursive: true,
	}
	a.NotError(c.Sanitize())

	doc := spec.NewAPIDoc()
	a.NotNil(doc)
	erro, _, h := messagetest.MessageHandler()
	ParseFile(doc, h, "./testdata/testfile.h", c)
	a.NotError(doc.Sanitize())
	a.Equal(0, len(doc.Apis)).
		Equal(doc.Version, "1.1.1")
	h.Stop()
	a.Empty(erro.String())
}
