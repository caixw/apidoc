// SPDX-License-Identifier: MIT

package doc

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/input"
	"github.com/caixw/apidoc/v6/message/messagetest"
)

func TestParse(t *testing.T) {
	a := assert.New(t)

	php := &input.Options{
		Lang:      "php",
		Dir:       "../input/testdata",
		Recursive: true,
		Encoding:  "gbk",
	}
	a.NotError(php.Sanitize())

	c := &input.Options{
		Lang:      "c++",
		Dir:       "../input/testdata",
		Recursive: true,
	}
	a.NotError(c.Sanitize())

	doc := New()
	a.NotNil(doc)
	erro, _, h := messagetest.MessageHandler()
	doc.Parse(h, php, c)
	a.NotError(doc.Sanitize())
	a.Equal(2, len(doc.Apis)).
		Equal(doc.Version, "1.1.1")
	api := doc.Apis[0]
	a.Equal(api.Method, "GET")
	h.Stop()
	a.Empty(erro.String())
}

func TestParseFile(t *testing.T) {
	a := assert.New(t)

	c := &input.Options{
		Lang:      "c++",
		Dir:       "../input/testdata",
		Recursive: true,
	}
	a.NotError(c.Sanitize())

	doc := New()
	a.NotNil(doc)
	erro, _, h := messagetest.MessageHandler()
	doc.ParseFile(h, "../input/testdata/testfile.h", c)
	a.NotError(doc.Sanitize())
	a.Equal(0, len(doc.Apis)).
		Equal(doc.Version, "1.1.1")
	h.Stop()
	a.Empty(erro.String())
}

func TestBegin(t *testing.T) {
	a := assert.New(t)

	a.True(bytes.HasPrefix([]byte("<apidoc "), apidocBegin))
	a.True(bytes.HasPrefix([]byte("<apidoc\t"), apidocBegin))
	a.True(bytes.HasPrefix([]byte("<apidoc\n"), apidocBegin))

	a.True(bytes.HasPrefix([]byte("<api "), apiBegin))
	a.True(bytes.HasPrefix([]byte("<api\t"), apiBegin))
	a.True(bytes.HasPrefix([]byte("<api\n"), apiBegin))
}
