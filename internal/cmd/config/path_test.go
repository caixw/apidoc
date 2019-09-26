// SPDX-License-Identifier: MIT

package config

import (
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
			result: "data",
		},
		{
			path:   "./test/data",
			wd:     "./tex/data",
			result: "../../test/data",
		},

		{ // 无法计算，返回原值
			path:   "/test/data",
			wd:     "./tex/data",
			result: "/test/data",
		},
	}

	for index, item := range data {
		result := rel(item.path, item.wd)
		a.Equal(result, item.result, "not equal @%d,v1=%s,v2=%s", index, result, item.result)
	}
}
