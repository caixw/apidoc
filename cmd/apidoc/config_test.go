// SPDX-License-Identifier: MIT

package main

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/input"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
)

func TestGetPath(t *testing.T) {
	a := assert.New(t)
	home, err := user.Current()
	a.NotError(err).NotEmpty(home)
	hd := home.HomeDir

	wd, err := os.Getwd()
	a.NotError(err).NotEmpty(wd)

	// 指定 home，不依赖于 wd
	path, err := getPath("~/path", "")
	a.NotError(err).Equal(path, filepath.Join(hd, "/path"))

	// 绝对路径
	path, err = getPath("/path", "")
	a.NotError(err).Equal(path, "/path")

	path, err = getPath("path", "")
	a.NotError(err).Equal(path, filepath.Join(wd, "/path"))

	path, err = getPath("./path", "")
	a.NotError(err).Equal(path, filepath.Join(wd, "/path"))

	// 以下为 wd= /wd

	// 指定 home，不依赖于 wd
	path, err = getPath("~/path", "/wd")
	a.NotError(err).Equal(path, filepath.Join(hd, "/path"))

	// 绝对路径
	path, err = getPath("/path", "/wd")
	a.NotError(err).Equal(path, "/path")

	path, err = getPath("path", "/wd")
	a.NotError(err).Equal(path, "/wd/path")

	path, err = getPath("./path", "/wd")
	a.NotError(err).Equal(path, "/wd/path")

	// 以下为 wd= ~/wd

	// 指定 home，不依赖于 wd
	path, err = getPath("~/path", "~/wd")
	a.NotError(err).Equal(path, filepath.Join(hd, "/path"))

	// 绝对路径
	path, err = getPath("/path", "~/wd")
	a.NotError(err).Equal(path, "/path")

	path, err = getPath("path", "~/wd")
	a.NotError(err).Equal(path, filepath.Join(hd, "/wd/path"))

	path, err = getPath("./path", "~/wd")
	a.NotError(err).Equal(path, filepath.Join(hd, "/wd/path"))

	// 以下为 wd= ./wd

	// 指定 home，不依赖于 wd
	path, err = getPath("~/path", "./wd")
	a.NotError(err).Equal(path, filepath.Join(hd, "/path"))

	// 绝对路径
	path, err = getPath("/path", "~/wd")
	a.NotError(err).Equal(path, "/path")

	path, err = getPath("path", "./wd")
	a.NotError(err).Equal(path, filepath.Join(wd, "/wd/path"))

	path, err = getPath("./path", "./wd")
	a.NotError(err).Equal(path, filepath.Join(wd, "/wd/path"))
}

func TestConfig_generateConfig_loadConfig(t *testing.T) {
	a := assert.New(t)

	wd, err := filepath.Abs("./")
	a.NotError(err).NotEmpty(wd)

	a.NotError(generateConfig(wd, filepath.Join(wd, configFilename)))
	cfg, err := loadConfig(wd)
	a.NotError(err).NotNil(cfg)

	a.Equal(cfg.Version, vars.Version())
}

func TestConfig_sanitize(t *testing.T) {
	a := assert.New(t)

	conf := &config{}
	err := conf.sanitize()
	a.Error(err)
	a.Equal(err.(*message.SyntaxError).Field, "version")

	// 版本号错误
	conf.Version = "5.0"
	err = conf.sanitize()
	a.Error(err)
	a.Equal(err.(*message.SyntaxError).Field, "version")

	// 未声明 inputs
	conf.Version = "5.0.1"
	err = conf.sanitize()
	a.Error(err)
	a.Equal(err.(*message.SyntaxError).Field, "inputs")

	// 未声明 output
	conf.Inputs = []*input.Options{{}}
	err = conf.sanitize()
	a.Error(err)
	a.Equal(err.(*message.SyntaxError).Field, "output")
}
