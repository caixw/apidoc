// SPDX-License-Identifier: MIT

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/issue9/version"
	"gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
	"github.com/caixw/apidoc/v5/options"
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
	Inputs []*options.Input `yaml:"inputs"`

	// 输出配置项。
	Output *options.Output `yaml:"output"`

	// 配置文件所在的目录。
	//
	// 如果 input 和 output 中涉及到地址为非绝对目录，则使用此值作为基地址。
	wd string
}

// 处理 path，如果 path 是相对路径的，则将其设置为相对于 wd 的路径
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
		dir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(dir, path[2:])
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
		return message.NewError(configFilename, "version", 0, locale.ErrInvalidFormat)
	}

	// 比较版本号兼容问题
	compatible, err := version.SemVerCompatible(vars.Version(), cfg.Version)
	if err != nil {
		return message.WithError(configFilename, "version", 0, err)
	}
	if !compatible {
		return message.NewError(configFilename, "version", 0, locale.VersionInCompatible)
	}

	if len(cfg.Inputs) == 0 {
		return message.NewError(configFilename, "inputs", 0, locale.ErrRequired)
	}

	if cfg.Output == nil {
		return message.NewError(configFilename, "output", 0, locale.ErrRequired)
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
		return nil, message.NewError("", "", 0, locale.ErrNotFoundSupportedLang)
	}

	return &config{
		Version: vars.Version(),
		Inputs:  inputs,
		Output: &options.Output{
			Type: options.ApidocXML,
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
