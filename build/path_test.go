// SPDX-License-Identifier: MIT

package build

import (
	"os"
	"testing"

	"github.com/caixw/apidoc/v7/core"
	"github.com/issue9/assert"
)

func TestAbs(t *testing.T) {
	a := assert.New(t)
	hd, err := os.UserHomeDir()
	a.NotError(err).NotNil(hd)
	hdURI := core.FileURI(hd)
	a.NotEmpty(hdURI)

	data := []*struct {
		path, wd, result string
	}{
		{
			path:   "~/path",
			wd:     "file:///wd/",
			result: hdURI.Append("/path").String(),
		},
		{
			path:   "/path",
			wd:     "file:///wd/",
			result: "/path",
		},
		{
			path:   "file:///path",
			wd:     "file:///wd/",
			result: "file:///path",
		},
		{
			path:   "path",
			wd:     "file:///wd/",
			result: "file:///wd/path",
		},
		{
			path:   "./path",
			wd:     "file:///wd/",
			result: "file:///wd/path",
		},
	}

	for index, item := range data {
		result, err := abs(core.URI(item.path), core.URI(item.wd))
		a.NotError(err, "err @%d,%s", index, err).
			Equal(result, item.result, "not equal @%d,v1=%s,v2=%s", index, result, item.result)
	}
}
