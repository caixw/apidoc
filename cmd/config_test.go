// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"path/filepath"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/internal/errors"
	"github.com/caixw/apidoc/internal/vars"
	"github.com/caixw/apidoc/options"
)

func TestConfig_generateConfig_loadConfig(t *testing.T) {
	a := assert.New(t)

	path, err := filepath.Abs("./.apidoc.yaml")
	a.NotError(err).NotEmpty(path)

	a.NotError(generateConfig("./..", path))
	cfg, err := loadConfig(path)
	a.NotError(err).NotNil(cfg)

	a.Equal(cfg.Version, vars.Version())
}

func TestConfig_sanitize(t *testing.T) {
	a := assert.New(t)

	conf := &config{}
	err := conf.sanitize()
	a.Error(err)
	a.Equal(err.(*errors.Error).Field, "version")

	// 版本号错误
	conf.Version = "5.0"
	err = conf.sanitize()
	a.Error(err)
	a.Equal(err.(*errors.Error).Field, "version")

	// 未声明 inputs
	conf.Version = "5.0.1"
	err = conf.sanitize()
	a.Error(err)
	a.Equal(err.(*errors.Error).Field, "inputs")

	// 未声明 output
	conf.Inputs = []*options.Input{{}}
	err = conf.sanitize()
	a.Error(err)
	a.Equal(err.(*errors.Error).Field, "output")
}
