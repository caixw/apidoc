// SPDX-License-Identifier: MIT

package main

import (
	"path/filepath"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/input"
	"github.com/caixw/apidoc/v5/message"
)

func TestFixedSyntaxError(t *testing.T) {
	a := assert.New(t)

	err := &message.SyntaxError{
		Field: "f",
	}
	e := fixedSyntaxError(err, "p")
	a.Equal(e.File, configFilename)
	a.Equal(e.Field, "p.f")

	err.Field = ""
	e = fixedSyntaxError(err, "p")
	a.Equal(e.File, configFilename)
	a.Equal(e.Field, "p")
}

func TestConfig_generateConfig_loadConfig(t *testing.T) {
	a := assert.New(t)

	wd, err := filepath.Abs("./")
	a.NotError(err).NotEmpty(wd)

	a.NotError(generateConfig(wd, filepath.Join(wd, configFilename)))
	cfg, err := loadConfig(wd)
	a.NotError(err).
		NotNil(cfg)

	a.Equal(cfg.Version, apidoc.Version())
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
