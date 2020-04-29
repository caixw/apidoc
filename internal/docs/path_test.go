// SPDX-License-Identifier: MIT

package docs

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
)

func TestDir(t *testing.T) {
	a := assert.New(t)

	p1, err := core.FileURI("../../docs")
	a.NotError(err).NotEmpty(p1)
	p2 := Dir()
	a.NotError(err).NotEmpty(p2)
	a.Equal(p1, p2)

	exists, err := Dir().Exists()
	a.NotError(err).True(exists)
}
