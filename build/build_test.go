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
	a.NotError(php.sanitize())

	c := &Input{
		Lang:      "c++",
		Dir:       "./testdata",
		Recursive: true,
	}
	a.NotError(c.sanitize())

	rslt := messagetest.NewMessageHandler()
	doc, err := parse(rslt.Handler, php, c)
	a.NotError(err).NotNil(doc)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	a.Equal(2, len(doc.APIs)).
		Equal(doc.Version.V(), "1.1.1")
	api := doc.APIs[0]
	a.Equal(api.Method.V(), "GET")
}
