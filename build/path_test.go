// SPDX-License-Identifier: MIT

package build

import (
	"os"
	"runtime"
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/caixw/apidoc/v7/core"
)

type pathTester struct {
	path, wd, result string
}

func TestAbs(t *testing.T) {
	if runtime.GOOS == "windows" { // windows 由 path_windows_test.go 程序测试
		return
	}

	a := assert.New(t, false)
	hd, err := os.UserHomeDir()
	a.NotError(err).NotNil(hd)
	hdURI := core.FileURI(hd)
	a.NotEmpty(hdURI)

	data := []*pathTester{
		{
			path:   "",
			wd:     "file:///wd/",
			result: "file:///wd/",
		},
		{
			path:   "~/path",
			wd:     "file:///wd/",
			result: hdURI.Append("/path").String(),
		},
		{
			path:   "/path",
			wd:     "file:///wd/",
			result: "file:///path",
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
			path:   "../path",
			wd:     "file:///wd/",
			result: "file:///path",
		},
		{
			path:   "../../path",
			wd:     "file:///wd/dir",
			result: "file:///path",
		},
		{
			path:   "./path",
			wd:     "file:///wd/",
			result: "file:///wd/path",
		},
		{
			path:   "path",
			wd:     "/wd/",
			result: "file:///wd/path",
		},
		{
			path:   "./path",
			wd:     "/wd/",
			result: "file:///wd/path",
		},
	}

	for index, item := range data {
		result, err := abs(core.URI(item.path), core.URI(item.wd))

		a.NotError(err, "err @%d,%s", index, err).
			Equal(result, item.result, "not equal @%d,v1=%s,v2=%s", index, result, item.result)
	}
}

func TestRel(t *testing.T) {
	if runtime.GOOS == "windows" { // windows 由 path_windows_test.go 程序测试
		return
	}

	a := assert.New(t, false)
	hd, err := os.UserHomeDir()
	a.NotError(err).NotNil(hd)
	hdURI := core.FileURI(hd)
	a.NotEmpty(hdURI)

	data := []*pathTester{
		{
			path:   "",
			wd:     "file:///wd/",
			result: "",
		},
		{
			path:   "/wd/path",
			wd:     "file:///wd/",
			result: "path",
		},
		{
			path:   "/wd/path",
			wd:     "file:///wd1/",
			result: "/wd/path",
		},
		{
			path:   "wd/path",
			wd:     "file:///wd/",
			result: "wd/path",
		},
		{
			path:   "path",
			wd:     "file:///wd/",
			result: "path",
		},
		{
			path:   "path",
			wd:     "/wd/",
			result: "path",
		},
		{
			path:   "file:///wd/path",
			wd:     "/wd/",
			result: "path",
		},
	}

	for index, item := range data {
		result, err := rel(core.URI(item.path), core.URI(item.wd))

		a.NotError(err, "err @%d,%s", index, err).
			Equal(result, item.result, "not equal @%d,v1=%s,v2=%s", index, result, item.result)
	}
}
