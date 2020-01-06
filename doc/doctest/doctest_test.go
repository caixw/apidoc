// SPDX-License-Identifier: MIT

package doctest

import (
	"path/filepath"
	"testing"

	"github.com/issue9/assert"
)

func TestPath(t *testing.T) {
	a := assert.New(t)

	p1, err := filepath.Abs(Filename)
	a.NotError(err).NotEmpty(p1)

	p2, err := filepath.Abs(Path(a))
	a.NotError(err).NotEmpty(p2)

	a.Equal(p1, p2)
}

func TestDir(t *testing.T) {
	a := assert.New(t)

	p1, err := filepath.Abs("./")
	a.NotError(err).NotEmpty(p1)

	p2, err := filepath.Abs(Dir(a))
	a.NotError(err).NotEmpty(p2)

	a.Equal(p1, p2)
}
