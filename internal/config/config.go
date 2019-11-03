// SPDX-License-Identifier: MIT

// Package config 管理配置文件的相关功能
package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/issue9/utils"
	"github.com/issue9/version"
	"gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/v5/input"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
	"github.com/caixw/apidoc/v5/output"
)

// 允许的配置文件名称
var configFilenames = []string{".apidoc.yaml", ".apidoc.yml"}

// Config 项目的配置内容
type Config struct {
	// 产生此配置文件的程序版本号
	//
	// 程序会用此来判断程序的兼容性。
	Version string `yaml:"version"`

	// 输入的配置项，可以指定多个项目
	//
	// 多语言项目，可能需要用到多个输入面。
	Inputs []*input.Options `yaml:"inputs"`

	// 输出配置项
	Output *output.Options `yaml:"output"`

	// 配置文件所在的目录
	//
	// 如果 input 和 output 中涉及到地址为非绝对目录，则使用此值作为基地址。
	wd string
}

// Load 加载指定目录下的配置文件
func Load(wd string) (*Config, error) {
	for _, filename := range configFilenames {
		path := filepath.Join(wd, filename)
		if utils.FileExists(path) {
			return loadFile(wd, path)
		}
	}

	return nil, message.NewLocaleError("", filepath.Join(wd, configFilenames[0]), 0, locale.ErrRequired)
}

func loadFile(wd, path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, message.WithError(path, "", 0, err)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, message.WithError(path, "", 0, err)
	}
	cfg.wd = wd

	if err := cfg.sanitize(path); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) sanitize(path string) *message.SyntaxError {
	// 比较版本号兼容问题
	compatible, err := version.SemVerCompatible(vars.Version(), cfg.Version)
	if err != nil {
		return message.WithError(path, "version", 0, err)
	}
	if !compatible {
		return message.NewLocaleError(path, "version", 0, locale.VersionInCompatible)
	}

	if len(cfg.Inputs) == 0 {
		return message.NewLocaleError(path, "inputs", 0, locale.ErrRequired)
	}

	if cfg.Output == nil {
		return message.NewLocaleError(path, "output", 0, locale.ErrRequired)
	}

	for index, i := range cfg.Inputs {
		field := "inputs[" + strconv.Itoa(index) + "]"

		if i.Dir, err = abs(i.Dir, cfg.wd); err != nil {
			return message.WithError(path, field+".path", 0, err)
		}
	}

	if cfg.Output.Path, err = abs(cfg.Output.Path, cfg.wd); err != nil {
		return message.WithError(path, "output.path", 0, err)
	}

	return nil
}

func detectConfig(wd string, recursive bool) (*Config, error) {
	inputs, err := input.Detect(wd, recursive)
	if err != nil {
		return nil, err
	}
	if len(inputs) == 0 {
		return nil, message.NewLocaleError("", "", 0, locale.ErrNotFoundSupportedLang)
	}

	for _, i := range inputs {
		i.Dir = rel(i.Dir, wd)
	}

	return &Config{
		Version: vars.Version(),
		Inputs:  inputs,
		Output: &output.Options{
			Path: rel(filepath.Join(wd, "apidoc.xml"), wd),
		},
	}, nil
}

// Write 根据 wd 所在目录的内容生成一个配置文件，并写入到该目录配置文件中。
//
// wd 表示当前程序的工作目录，根据此目录的内容检测其语言特性。
func Write(wd string, recursive bool) error {
	cfg, err := detectConfig(wd, recursive)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	path, err := abs(configFilenames[0], wd)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, os.ModePerm)
}
