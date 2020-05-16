// SPDX-License-Identifier: MIT

package build

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core/messagetest"
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

	erro, _, h := messagetest.MessageHandler()
	doc, err := parse(h, php, c)
	a.NotError(err).NotNil(doc)
	h.Stop()
	a.Empty(erro.String())

	a.Equal(2, len(doc.Apis)).
		Equal(doc.Version.V(), "1.1.1")
	api := doc.Apis[0]
	a.Equal(api.Method.V(), "GET")
}
