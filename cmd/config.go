// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/issue9/version"
	"golang.org/x/text/message"
	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/internal/locale"
	"github.com/caixw/apidoc/internal/options"
	"github.com/caixw/apidoc/internal/vars"
	"github.com/caixw/apidoc/output"
)

// 配置文件名称。
const configFilename = ".apidoc.yaml"

type configError struct {
	options.FieldError
	File string
}

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

func (err *configError) Error() string {
	msg := locale.Sprintf(err.MessageKey, err.MessageArgs...)
	return locale.Sprintf(locale.ErrConfig, err.File, err.Field, msg)
}

func newConfigError(field string, key message.Reference, args ...interface{}) error {
	err := options.NewFieldError(field, key, args)
	return newConfigFieldError(err)
}

func newConfigFieldError(err *options.FieldError) error {
	return &configError{
		File:       configFilename,
		FieldError: *err,
	}
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
		return newConfigError("version", locale.Sprintf(locale.ErrInvalidFormat))
	}

	// 比较版本号兼容问题
	compatible, err := version.SemVerCompatible(vars.Version(), cfg.Version)
	if err != nil {
		return newConfigError("version", err.Error())
	}
	if !compatible {
		return newConfigError("version", locale.Sprintf(locale.VersionInCompatible))
	}

	if len(cfg.Inputs) == 0 {
		return newConfigError("inputs", locale.Sprintf(locale.ErrRequired))
	}

	if cfg.Output == nil {
		return newConfigError("output", locale.Sprintf(locale.ErrRequired))
	}

	for i, opt := range cfg.Inputs {
		if err := opt.Sanitize(); err != nil {
			if cerr, ok := err.(*options.FieldError); ok {
				index := strconv.Itoa(i)
				cerr.Field = "inputs[" + index + "]." + cerr.Field
				err = newConfigFieldError(cerr)
			}
			return err
		}
	}

	if err := cfg.Output.Sanitize(); err != nil {
		if cerr, ok := err.(*options.FieldError); ok {
			cerr.Field = "outputs." + cerr.Field
			err = newConfigFieldError(cerr)
		}
		return err
	}

	return nil
}

// 根据 wd 所在目录的内容生成一个配置文件，并写入到 path  中
//
// wd 表示当前程序的工作目录，根据此目录的内容检测其语言特性。
// path 表示生成的配置文件存放的路径。
func generateConfig(wd, path string) error {
	o, err := input.Detect(wd, true)
	if err != nil {
		return err
	}

	cfg := &config{
		Version: vars.Version(),
		Inputs:  []*input.Options{o},
		Output: &output.Options{
			Path: filepath.Join(o.Dir, "apidoc.json"),
		},
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, os.ModePerm)
}
