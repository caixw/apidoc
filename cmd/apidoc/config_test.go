// SPDX-License-Identifier: MIT

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/input"
	"github.com/caixw/apidoc/v5/internal/vars"
)

func TestGetPath(t *testing.T) {
	a := assert.New(t)
	hd, err := os.UserHomeDir()
	a.NotError(err).NotNil(hd)

	wd, err := os.Getwd()
	a.NotError(err).NotEmpty(wd)

	abs := func(path string) string {
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
			result: abs(filepath.Join(hd, "/path")),
		},
		{ // 绝对路径
			path:   "/path",
			wd:     "",
			result: abs("/path"),
		},
		{
			path:   "path",
			wd:     "",
			result: abs(filepath.Join(wd, "/path")),
		},
		{
			path:   "./path",
			wd:     "",
			result: abs(filepath.Join(wd, "/path")),
		},

		// 以下为 wd= /wd
		{ // 指定 home，不依赖于 wd
			path:   "~/path",
			wd:     "/wd",
			result: abs(filepath.Join(hd, "/path")),
		},
		{ // 绝对路径
			path:   "/path",
			wd:     "/wd",
			result: abs("/path"),
		},
		{
			path:   "path",
			wd:     "/wd",
			result: abs("/wd/path"),
		},
		{
			path:   "./path",
			wd:     "/wd",
			result: abs("/wd/path"),
		},

		// 以下为 wd= ~/wd
		{ // 指定 home，不依赖于 wd
			path:   "~/path",
			wd:     "~/wd",
			result: abs(filepath.Join(hd, "/path")),
		},
		{ // 绝对路径
			path:   "/path",
			wd:     "~/wd",
			result: abs("/path"),
		},
		{
			path:   "path",
			wd:     "~/wd",
			result: abs(filepath.Join(hd, "/wd/path")),
		},
		{
			path:   "./path",
			wd:     "~/wd",
			result: abs(filepath.Join(hd, "/wd/path")),
		},

		// 以下为 wd= ./wd
		{ // 指定 home，不依赖于 wd
			path:   "~/path",
			wd:     "./wd",
			result: abs(filepath.Join(hd, "/path")),
		},
		{ // 绝对路径
			path:   "/path",
			wd:     "./wd",
			result: abs("/path"),
		},
		{
			path:   "path",
			wd:     "./wd",
			result: abs(filepath.Join(wd, "/wd/path")),
		},
		{
			path:   "./path",
			wd:     "./wd",
			result: abs(filepath.Join(wd, "/wd/path")),
		},
	}

	for index, item := range data {
		result, err := getPath(item.path, item.wd)
		a.NotError(err, "err @%d,%s", index, err).
			Equal(result, item.result, "not equal @%d,v1=%s,v2=%s", index, result, item.result)
	}
}

func TestConfig_generateConfig_loadConfig(t *testing.T) {
	a := assert.New(t)

	wd, err := filepath.Abs("./")
	a.NotError(err).NotEmpty(wd)

	a.NotError(generateConfig(wd, filepath.Join(wd, configFilename)))
	cfg, err := loadConfig(wd)
	a.NotError(err).
		NotNil(cfg)

	a.Equal(cfg.Version, vars.Version())
}

func TestConfig_sanitize(t *testing.T) {
	a := assert.New(t)

	conf := &config{}
	err := conf.sanitize()
	a.Error(err).
		Equal(err.Field, "version")

	// 版本号错误
	conf.Version = "5.0"
	err = conf.sanitize()
	a.Error(err).
		Equal(err.Field, "version")

	// 未声明 inputs
	conf.Version = "5.0.1"
	err = conf.sanitize()
	a.Error(err).
		Equal(err.Field, "inputs")

	// 未声明 output
	conf.Inputs = []*input.Options{{}}
	err = conf.sanitize()
	a.Error(err).
		Equal(err.Field, "output")
}
