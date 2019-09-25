// SPDX-License-Identifier: MIT

// +build !windows

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/issue9/assert"
)

func TestAbs(t *testing.T) {
	a := assert.New(t)
	hd, err := os.UserHomeDir()
	a.NotError(err).NotNil(hd)

	wd, err := os.Getwd()
	a.NotError(err).NotEmpty(wd)

	absPath := func(path string) string {
		abs, err := filepath.Abs(path)
		a.NotError(err)
		return abs
	}

	data := []*struct {
		path, wd, result string
	}{
		{ // 指定 home，不依赖于 wd
			path:   "~/path",
			wd:     "",
			result: absPath(filepath.Join(hd, "/path")),
		},
		{ // 绝对路径
			path:   "/path",
			wd:     "",
			result: absPath("/path"),
		},
		{
			path:   "path",
			wd:     "",
			result: absPath(filepath.Join(wd, "/path")),
		},
		{
			path:   "./path",
			wd:     "",
			result: absPath(filepath.Join(wd, "/path")),
		},

		// 以下为 wd= /wd
		{ // 指定 home，不依赖于 wd
			path:   "~/path",
			wd:     "/wd",
			result: absPath(filepath.Join(hd, "/path")),
		},
		{ // 绝对路径
			path:   "/path",
			wd:     "/wd",
			result: absPath("/path"),
		},
		{
			path:   "path",
			wd:     "/wd",
			result: absPath("/wd/path"),
		},
		{
			path:   "./path",
			wd:     "/wd",
			result: absPath("/wd/path"),
		},

		// 以下为 wd= ~/wd
		{ // 指定 home，不依赖于 wd
			path:   "~/path",
			wd:     "~/wd",
			result: absPath(filepath.Join(hd, "/path")),
		},
		{ // 绝对路径
			path:   "/path",
			wd:     "~/wd",
			result: absPath("/path"),
		},
		{
			path:   "path",
			wd:     "~/wd",
			result: absPath(filepath.Join(hd, "/wd/path")),
		},
		{
			path:   "./path",
			wd:     "~/wd",
			result: absPath(filepath.Join(hd, "/wd/path")),
		},

		// 以下为 wd= ./wd
		{ // 指定 home，不依赖于 wd
			path:   "~/path",
			wd:     "./wd",
			result: absPath(filepath.Join(hd, "/path")),
		},
		{ // 绝对路径
			path:   "/path",
			wd:     "./wd",
			result: absPath("/path"),
		},
		{
			path:   "path",
			wd:     "./wd",
			result: absPath(filepath.Join(wd, "/wd/path")),
		},
		{
			path:   "./path",
			wd:     "./wd",
			result: absPath(filepath.Join(wd, "/wd/path")),
		},
	}

	for index, item := range data {
		result, err := abs(item.path, item.wd)
		a.NotError(err, "err @%d,%s", index, err).
			Equal(result, item.result, "not equal @%d,v1=%s,v2=%s", index, result, item.result)
	}
}
