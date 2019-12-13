// SPDX-License-Identifier: MIT

package static

import (
	"testing"

	"github.com/issue9/assert"
)

func TestPack(t *testing.T) {
	a := assert.New(t)
	a.NotError(Pack("./testdir", "testdata", "Data", "./testdata/testdata.go"))
}

func TestGetFileInfos(t *testing.T) {
	a := assert.New(t)

	info, err := getFileInfos("./testdir")
	a.NotError(err).NotNil(info)
	a.Equal(5, len(info))
}

func TestGetPkgPath(t *testing.T) {
	a := assert.New(t)

	p, err := getPkgPath("")
	a.Error(err).Empty(p)

	p, err = getPkgPath("./testdir/go.mod1")
	a.NotError(err).Equal(p, "test/v6")

	p, err = getPkgPath("./testdir/go.mod2")
	a.Error(err).Empty(p)

	p, err = getPkgPath("./testdir/go.mod3")
	a.Error(err).Empty(p)
}
