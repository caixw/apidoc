// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"strconv"

	"github.com/issue9/version"
	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/openapi"
	"github.com/caixw/apidoc/output"
)

// 项目的配置内容
type config struct {
	// 产生此配置文件的程序版本号。
	//
	// 程序会用此来判断程序的兼容性。
	Version string `yaml:"version"`

	// 输入的配置项，可以指定多个项目
	//
	// 多语言项目，可能需要用到多个输入面。
	// 但是输出内容依然会被集中到 Output 一个字段中。
	Inputs []*input.Options `yaml:"inputs"`

	Output *output.Options `yaml:"output"`
}

// 加载 path 所指的文件内容到 *config 实例。
func loadConfig(path string) (*config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if err = cfg.sanitize(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *config) sanitize() error {
	if !version.SemVerValid(cfg.Version) {
		return &openapi.Error{Field: "version", Message: locale.Sprintf(locale.ErrInvalidFormat)}
	}

	if len(cfg.Inputs) == 0 {
		return &openapi.Error{Field: "inputs", Message: locale.Sprintf(locale.ErrRequired)}
	}

	if cfg.Output == nil {
		return &openapi.Error{Field: "output", Message: locale.Sprintf(locale.ErrRequired)}
	}

	for i, opt := range cfg.Inputs {
		if err := opt.Sanitize(); err != nil {
			index := strconv.Itoa(i)
			err.Field = "inputs[" + index + "]." + err.Field
			return err
		}
	}

	if err := cfg.Output.Sanitize(); err != nil {
		err.Field = "outputs." + err.Field
		return err
	}

	return nil
}
