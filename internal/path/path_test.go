// SPDX-License-Identifier: MIT

package path

import (
	"path/filepath"
	"testing"

	"github.com/issue9/assert"
)

func TestRel(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		path, wd, result string
	}{
		{
			path:   "./test/data",
			wd:     "./test",
			result: filepath.Clean("data"),
		},
		{
			path:   "./test/data",
			wd:     "./tex/data",
			result: filepath.Clean("../../test/data"),
		},

		{ // 无法计算，返回原值
			path:   "/test/data",
			wd:     "./tex/data",
			result: filepath.Clean("/test/data"),
		},
	}

	for index, item := range data {
		result := filepath.Clean(Rel(item.path, item.wd))
		a.Equal(result, item.result, "not equal @%d,v1=%s,v2=%s", index, result, item.result)
	}
}

func TestCurrPath(t *testing.T) {
	a := assert.New(t)

	dir, err := filepath.Abs("./")
	a.NotError(err).NotEmpty(dir)

	d, err := filepath.Abs(CurrPath("./"))
	a.NotError(err).NotEmpty(d)

	a.Equal(d, dir)
}
