// SPDX-License-Identifier: MIT

package html

import (
	"testing"

	"github.com/issue9/assert"
)

func TestGet(t *testing.T) {
	a := assert.New(t)

	data, ct := Get("apidoc.xsl")
	a.NotNil(data).NotEmpty(ct)

	data, ct = Get("not-exists")
	a.Nil(data).Empty(ct)
}
