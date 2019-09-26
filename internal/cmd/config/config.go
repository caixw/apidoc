// SPDX-License-Identifier: MIT

// Package config 负责命令行工具 apidoc 的配置文件管理
package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/issue9/version"
	"gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/input"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
	"github.com/caixw/apidoc/v5/output"
)

// 采用 detect 检测目录内容时，需要赋予的一些默认值
const (
	configFilename = ".apidoc.yaml"
	docFilename    = "apidoc.xml"
)

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

// Load 加载指定的配置文件
func Load(path string) (*Config, *message.SyntaxError) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, message.WithError(configFilename, "", 0, err)
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, message.WithError(configFilename, "", 0, err)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, message.WithError(configFilename, "", 0, err)
	}
	cfg.wd = filepath.Dir(path)

	if err := cfg.sanitize(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) sanitize() *message.SyntaxError {
	// 比较版本号兼容问题
	compatible, err := version.SemVerCompatible(apidoc.Version(), cfg.Version)
	if err != nil {
		return message.WithError(configFilename, "version", 0, err)
	}
	if !compatible {
		return message.NewLocaleError(configFilename, "version", 0, locale.VersionInCompatible)
	}

	if len(cfg.Inputs) == 0 {
		return message.NewLocaleError(configFilename, "inputs", 0, locale.ErrRequired)
	}

	if cfg.Output == nil {
		return message.NewLocaleError(configFilename, "output", 0, locale.ErrRequired)
	}

	for index, i := range cfg.Inputs {
		field := "inputs[" + strconv.Itoa(index) + "]"

		if i.Dir, err = abs(i.Dir, cfg.wd); err != nil {
			return message.WithError(configFilename, field+".path", 0, err)
		}
	}

	if cfg.Output.Path, err = abs(cfg.Output.Path, cfg.wd); err != nil {
		return message.WithError(configFilename, "output.path", 0, err)
	}

	return nil
}

func fixedSyntaxError(err *message.SyntaxError, field string) *message.SyntaxError {
	err.File = configFilename
	if err.Field == "" {
		err.Field = field
	} else {
		err.Field = field + "." + err.Field
	}
	return err
}

func detectConfig(wd string) (*Config, error) {
	inputs, err := input.Detect(wd, true)
	if err != nil {
		return nil, err
	}
	if len(inputs) == 0 {
		return nil, message.NewLocaleError("", "", 0, locale.ErrNotFoundSupportedLang)
	}

	for _, i := range inputs {
		i.Dir = rel(wd, i.Dir)
	}

	outputFile := rel(wd, filepath.Join(wd, docFilename))
	return &Config{
		Version: apidoc.Version(),
		Inputs:  inputs,
		Output: &output.Options{
			Type: output.ApidocXML,
			Path: outputFile,
		},
	}, nil
}

// Write 根据 wd 所在目录的内容生成一个配置文件，并写入到该目录配置文件中。
//
// wd 表示当前程序的工作目录，根据此目录的内容检测其语言特性。
// path 表示生成的配置文件存放的路径。
func Write(wd string) error {
	cfg, err := detectConfig(wd)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	path, err := abs(configFilename, wd)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, os.ModePerm)
}
