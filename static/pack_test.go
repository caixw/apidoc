// SPDX-License-Identifier: MIT

package static

import (
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/internal/vars"
)

// 保证 modulePath 的值正确性
func TestModulePath(t *testing.T) {
	a := assert.New(t)

	suffix := path.Join(vars.DocVersion(), "static")
	a.True(strings.HasSuffix(modulePath, suffix))
}

func TestPack(t *testing.T) {
	a := assert.New(t)
	a.NotError(Pack("./testdir", "testdata", "Data", "./testdata/testdata.go", TypeAll))
}

func TestGetFileInfos(t *testing.T) {
	a := assert.New(t)

	info, err := getFileInfos("./testdir", TypeAll)
	a.NotError(err).NotNil(info)
	a.Equal(3, len(info))

	// 采用绝对路径
	dir, err := filepath.Abs("./testdir")
	a.NotError(err).NotEmpty(dir)
	info, err = getFileInfos(dir, TypeAll)
	a.NotError(err).NotNil(info)
	a.Equal(3, len(info))
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
