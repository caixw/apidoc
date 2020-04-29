// SPDX-License-Identifier: MIT

package asttest

import (
	"path/filepath"
	"testing"

	"github.com/issue9/assert"
)

func TestXML(t *testing.T) {
	a := assert.New(t)
	data := XML(a)
	a.NotNil(data)
}

func TestURI(t *testing.T) {
	a := assert.New(t)

	p1, err := filepath.Abs(Filename)
	a.NotError(err).NotEmpty(p1)

	p2, err := URI(a).File()
	a.NotError(err).NotEmpty(p2)

	a.Equal(p1, p2)
}

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
