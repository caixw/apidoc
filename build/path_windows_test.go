// SPDX-License-Identifier: MIT

package build

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/caixw/apidoc/v7/core"
)

func TestAbs_windows(t *testing.T) {
	a := assert.New(t, false)
	hd, err := os.UserHomeDir()
	a.NotError(err).NotNil(hd)
	hdURI := core.FileURI(hd)
	a.NotEmpty(hdURI)

	root, err := filepath.Abs("/wd/")
	a.NotError(err)

	data := []*pathTester{
		{
			path:   "",
			wd:     core.FileURI(root).String(),
			result: core.FileURI(root).String(),
		},
		{
			path:   "~/path",
			wd:     "file:///wd/",
			result: hdURI.Append("/path").String(),
		},

		{
			path:   "~\\path",
			wd:     "file:///wd/",
			result: hdURI.Append("\\path").String(),
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
			wd:     core.FileURI(root).String(),
			result: core.FileURI(root).Append("path").String(),
		},
		{
			path:   "./path",
			wd:     core.FileURI(root).String(),
			result: core.FileURI(root).Append("path").String(),
		},

		{
			path:   ".\\path",
			wd:     core.FileURI(root).String(),
			result: core.FileURI(root).Append("path").String(),
		},

		{
			path:   "c:\\path",
			wd:     "file:///wd/",
			result: "file://c:\\path",
		},
	}

	for index, item := range data {
		result, err := abs(core.URI(item.path), core.URI(item.wd))

		a.NotError(err, "err @%d,%s", index, err).
			Equal(result, item.result, "not equal @%d,v1=%s,v2=%s", index, result, item.result)
	}
}

func TestRel_windows(t *testing.T) {
	a := assert.New(t, false)
	hd, err := os.UserHomeDir()
	a.NotError(err).NotNil(hd)
	hdURI := core.FileURI(hd)
	a.NotEmpty(hdURI)

	root, err := filepath.Abs("/wd/")
	a.NotError(err)

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
			path:   root + "/path",
			wd:     core.FileURI(root).String(),
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
		{
			path:   "c:/wd/path",
			wd:     "file:///wd1/",
			result: "c:/wd/path",
		},
		{
			path:   "c:\\wd\\path",
			wd:     "file://c:\\",
			result: "wd\\path",
		},
	}

	for index, item := range data {
		result, err := rel(core.URI(item.path), core.URI(item.wd))

		a.NotError(err, "err @%d,%s", index, err).
			Equal(result, item.result, "not equal @%d,v1=%s,v2=%s", index, result, item.result)
	}
}
