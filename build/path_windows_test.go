// SPDX-License-Identifier: MIT

package build

import (
	"os"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
)

func TestAbs_windows(t *testing.T) {
	a := assert.New(t)
	hd, err := os.UserHomeDir()
	a.NotError(err).NotNil(hd)
	hdURI := core.FileURI(hd)
	a.NotEmpty(hdURI)

	data := []*pathTester{
		{
			path:   "",
			wd:     "file:///wd/",
			result: "file://C:\\wd",
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
			wd:     "file:///wd/",
			result: "file://C:\\wd\\path",
		},
		{
			path:   "../../path",
			wd:     "file:///wd/dir", // 相当于 c:\wd\dir
			result: "file://C:\\path",
		},
		{
			path:   "./path",
			wd:     "file:///wd/",
			result: "file://C:\\wd\\path",
		},

		{
			path:   ".\\path",
			wd:     "file:///wd/",
			result: "file://C:\\wd\\path",
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
	a := assert.New(t)
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
			path:   "C:/wd/path",
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
