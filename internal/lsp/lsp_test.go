// SPDX-License-Identifier: MIT

package lsp

import (
	"testing"

	"github.com/issue9/assert"
)

func TestServe(t *testing.T) {
	a := assert.New(t)
	a.Error(Serve(true, "not-exists", "", nil, nil))
}
