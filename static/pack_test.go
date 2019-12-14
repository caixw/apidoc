// SPDX-License-Identifier: MIT

package static

import (
	"path/filepath"
	"testing"

	"github.com/issue9/assert"
)

func TestPack(t *testing.T) {
	a := assert.New(t)
	a.NotError(Pack("./testdir", "testdata", "Data", "./testdata/testdata.go", TypeAll))
}

func TestGetFileInfos(t *testing.T) {
	a := assert.New(t)

	info, err := getFileInfos("./testdir", TypeAll)
	a.NotError(err).NotNil(info)
	a.Equal(6, len(info))

	// 采有绝对路径
	dir, err := filepath.Abs("./testdir")
	a.NotError(err).NotEmpty(dir)
	info, err = getFileInfos(dir, TypeAll)
	a.NotError(err).NotNil(info)
	a.Equal(6, len(info))
}

func TestGetFileInfos_TypeStylesheet(t *testing.T) {
	a := assert.New(t)

	info, err := getFileInfos("./testdir", TypeStylesheet)
	a.NotError(err).NotNil(info)
	a.Equal(1, len(info)).
		Equal(info[0].Name, "icon.svg")
}

func TestGetFileInfos_TypeNone(t *testing.T) {
	a := assert.New(t)

	info, err := getFileInfos("./testdir", TypeNone)
	a.NotError(err).Nil(info)
	a.Equal(0, len(info))
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
