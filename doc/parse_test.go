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

	erro, _, h := messagetest.MessageHandler()

	php := &input.Options{
		Lang:      "php",
		Dir:       "../input/testdata",
		Recursive: true,
		Encoding:  "gbk",
	}

	c := &input.Options{
		Lang:      "c++",
		Dir:       "../input/testdata",
		Recursive: true,
	}

	doc := New()
	err := doc.Parse(h, php, c)
	a.NotError(err).NotNil(doc).
		Equal(1, len(doc.Apis)).
		Equal(doc.Version, "1.1.1")
	api := doc.Apis[0]
	a.Equal(api.Method, "GET")
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
