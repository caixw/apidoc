// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"strings"
	"testing"

	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/output"
	"github.com/issue9/assert"
	"github.com/issue9/term/colors"
)

func TestLogWriter_Write(t *testing.T) {
	a := assert.New(t)

	l := &logWriter{out: os.Stdout, prefix: "[TEST] ", color: colors.Red}
	_, err := l.Write([]byte("这是一行红色前缀的字"))
	a.NotError(err)
}

func TestConfig_sanitize(t *testing.T) {
	a := assert.New(t)

	conf := &config{}
	err := conf.sanitize()
	a.Equal(err.Field, "version")

	// 版本号错误
	conf.Version = "4.0"
	err = conf.sanitize()
	a.Equal(err.Field, "version")

	// 版本号不兼容
	conf.Version = "1000.0.1"
	err = conf.sanitize()
	a.Equal(err.Field, "version")

	// 未声明 inputs
	conf.Version = "4.0.1"
	err = conf.sanitize()
	a.Equal(err.Field, "inputs")

	// 未声明 output
	conf.Inputs = []*input.Options{&input.Options{}}
	err = conf.sanitize()
	a.Equal(err.Field, "output")

	// 查看错误提示格式是否正确
	conf.Output = &output.Options{}
	conf.Inputs = append(conf.Inputs, &input.Options{
		Lang: "123",
	})
	err = conf.sanitize()
	a.True(strings.HasPrefix(err.Field, "inputs[0]"))
}
