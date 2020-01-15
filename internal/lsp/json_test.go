// SPDX-License-Identifier: MIT

package lsp

import (
	"testing"

	"github.com/issue9/assert"
)

func TestConvertMethod(t *testing.T) {
	a := assert.New(t)

	a.Equal(convertMethod(""), "")
	a.Equal(convertMethod("a.b"), "a.b")
	a.Equal(convertMethod("a.B"), "a.B")
	a.Equal(convertMethod("a/B"), "a.B")
	a.Equal(convertMethod("a/b"), "a.B")
}
