// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/issue9/version"
	"golang.org/x/text/message"
	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/internal/errors"
	"github.com/caixw/apidoc/internal/locale"
	"github.com/caixw/apidoc/internal/vars"
	"github.com/caixw/apidoc/options"
)

// 配置文件名称。
const configFilename = ".apidoc.yaml"

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
	Inputs []*options.Input `yaml:"inputs"`

	Output *options.Output `yaml:"output"`

	// 配置文件所在的目录。
	//
	// 如果 input 和 output 中涉及到地址为非绝对目录，则使用此值作为基地址。
	wd string
}

func newConfigError(field string, key message.Reference, args ...interface{}) error {
	return &errors.Error{
		Field:       field,
		File:        configFilename,
		MessageKey:  key,
		MessageArgs: args,
	}
}

// 处理 path，如果 path 是相对路径的，则将其依附于 wd
//
// wd 表示工作目录；
// path 表示需要处理的路径。
func getPath(path, wd string) (p string, err error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}

	if !strings.HasPrefix(path, "~/") {
		path = filepath.Join(wd, path)
	}

	// 有可能 wd 是 ~/ 开头的
	if strings.HasPrefix(path, "~/") { // 非 home 路开头的相对路径，需要将其定位到 wd 目录之下
		u, err := user.Current()
		if err != nil {
			return "", err
		}

		path = filepath.Join(u.HomeDir, path[2:])
	}

	if !filepath.IsAbs(path) {
		if path, err = filepath.Abs(path); err != nil {
			return "", err
		}
	}

	return filepath.Clean(path), nil
}

// 加载 wd 目录下的配置文件到 *config 实例。
func loadConfig(wd string) (*config, error) {
	data, err := ioutil.ReadFile(filepath.Join(wd, configFilename))
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

	cfg.wd = wd

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

	for _, input := range cfg.Inputs {
		if input.Dir, err = getPath(cfg.wd, input.Dir); err != nil {
			return err
		}
	}

	if !filepath.IsAbs(cfg.Output.Path) {
		if cfg.Output.Path, err = getPath(cfg.wd, cfg.Output.Path); err != nil {
			return err
		}
	}

	return nil
}

func getConfig(wd string) (*config, error) {
	inputs, err := options.Detect(wd, true)
	if err != nil {
		return nil, err
	}
	if len(inputs) == 0 {
		return nil, locale.Errorf(locale.ErrNotFoundSupportedLang)
	}

	return &config{
		Version: vars.Version(),
		Inputs:  inputs,
		Output: &options.Output{
			Type: options.ApidocJSON,
			Path: filepath.Join(wd, "apidoc.json"),
		},
	}, nil
}

// 根据 wd 所在目录的内容生成一个配置文件，并写入到 path  中
//
// wd 表示当前程序的工作目录，根据此目录的内容检测其语言特性。
// path 表示生成的配置文件存放的路径。
func generateConfig(wd, path string) error {
	cfg, err := getConfig(wd)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, os.ModePerm)
}
