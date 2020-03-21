// SPDX-License-Identifier: MIT

package path

import (
	"path/filepath"
	"testing"

	"github.com/issue9/assert"
)

func TestCurrPath(t *testing.T) {
	a := assert.New(t)

	dir, err := filepath.Abs("./")
	a.NotError(err).NotEmpty(dir)

	d, err := filepath.Abs(CurrPath("./"))
	a.NotError(err).NotEmpty(d)

	a.Equal(d, dir)
}
