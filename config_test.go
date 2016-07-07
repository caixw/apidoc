// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"testing"

	"github.com/caixw/apidoc/app"
	"github.com/issue9/assert"
)

var _ io.Writer = &syntaxWriter{}

func TestConfig(t *testing.T) {
	a := assert.New(t)

	path := "./" + app.ConfigFilename
	a.NotError(genConfigFile(path))

	conf, err := loadConfig(path)
	a.NotError(err).NotNil(conf)
	a.Equal(conf.Inputs[0].Lang, "go")
}
