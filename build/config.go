// SPDX-License-Identifier: MIT

package build

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/issue9/utils"
	"github.com/issue9/version"
	"gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/path"
	"github.com/caixw/apidoc/v6/internal/vars"
	"github.com/caixw/apidoc/v6/message"
	"github.com/caixw/apidoc/v6/spec"
)

// Config 配置文件映身的结构
type Config struct {
	// 文档的版本信息
	//
	// 程序会用此来判断程序的兼容性。
	Version string `yaml:"version"`

	// 输入的配置项，可以指定多个项目
	//
	// 多语言项目，可能需要用到多个输入面。
	Inputs []*Input `yaml:"inputs"`

	// 输出配置项
	Output *Output `yaml:"output"`

	// 配置文件所在的目录
	//
	// 如果 input 和 output 中涉及到地址为非绝对目录，则使用此值作为基地址。
	wd string

	h *message.Handler
}

// LoadConfig 加载指定目录下的配置文件
//
// 所有的错误信息会输出到 h，在出错时，会返回 nil
func LoadConfig(h *message.Handler, wd string) *Config {
	for _, filename := range vars.AllowConfigFilenames {
		p := filepath.Join(wd, filename)
		if utils.FileExists(p) {
			cfg, err := loadFile(wd, p)
			if err != nil {
				h.Error(message.Erro, err)
				return nil
			}
			cfg.h = h
			return cfg
		}
	}

	msg := message.NewLocaleError("", filepath.Join(wd, vars.AllowConfigFilenames[0]), 0, locale.ErrRequired)
	h.Error(message.Erro, msg)
	return nil
}

// DetectConfig 检测 wd 内容并生成 Config 实例
func DetectConfig(wd string, recursive bool) (*Config, error) {
	inputs, err := detectInput(wd, recursive)
	if err != nil {
		return nil, err
	}
	if len(inputs) == 0 {
		return nil, message.NewLocaleError("", "", 0, locale.ErrNotFoundSupportedLang)
	}

	for _, i := range inputs {
		i.Dir = path.Rel(i.Dir, wd)
	}

	return &Config{
		Version: spec.Version,
		Inputs:  inputs,
		Output: &Output{
			Path: path.Rel(filepath.Join(wd, "apidoc.xml"), wd),
		},
		wd: wd,
	}, nil
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

func (cfg *Config) sanitize(file string) error {
	// 比较版本号兼容问题
	compatible, err := version.SemVerCompatible(spec.Version, cfg.Version)
	if err != nil {
		return message.WithError(file, "version", 0, err)
	}
	if !compatible {
		return message.NewLocaleError(file, "version", 0, locale.VersionInCompatible)
	}

	if len(cfg.Inputs) == 0 {
		return message.NewLocaleError(file, "inputs", 0, locale.ErrRequired)
	}

	if cfg.Output == nil {
		return message.NewLocaleError(file, "output", 0, locale.ErrRequired)
	}

	for index, i := range cfg.Inputs {
		field := "inputs[" + strconv.Itoa(index) + "]"

		if i.Dir, err = path.Abs(i.Dir, cfg.wd); err != nil {
			return message.WithError(file, field+".path", 0, err)
		}

		if err := i.Sanitize(); err != nil {
			if serr, ok := err.(*message.SyntaxError); ok {
				serr.File = file
				serr.Line = 0
				serr.Field = field + serr.Field
			}
			return err
		}
	}

	if cfg.Output.Path, err = path.Abs(cfg.Output.Path, cfg.wd); err != nil {
		return message.WithError(file, "output.path", 0, err)
	}
	return cfg.Output.Sanitize()
}

// Detect 根据 wd 所在目录的内容生成一个配置文件，并写入到该目录配置文件中。
//
// wd 表示当前程序的工作目录，根据此目录的内容检测其语言特性。
func Detect(wd string, recursive bool) error {
	cfg, err := DetectConfig(wd, recursive)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	p, err := path.Abs(vars.AllowConfigFilenames[0], wd)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(p, data, os.ModePerm)
}

// Build 解析文档并输出文档内容
//
// 具体信息可参考 Build 函数的相关文档。
func (cfg *Config) Build(start time.Time) {
	if err := Build(cfg.h, cfg.Output, cfg.Inputs...); err != nil {
		cfg.h.Error(message.Erro, err)
		return
	}

	cfg.h.Message(message.Succ, locale.Complete, cfg.Output.Path, time.Now().Sub(start))
}

// Buffer 根据 wd 目录下的配置文件生成文档内容并保存至内存
//
// 具体信息可参考 Buffer 函数的相关文档。
func (cfg *Config) Buffer() *bytes.Buffer {
	buf, err := Buffer(cfg.h, cfg.Output, cfg.Inputs...)
	if err != nil {
		cfg.h.Error(message.Erro, err)
		return nil
	}

	return buf
}

// Test 执行对语法内容的测试
func (cfg *Config) Test() {
	Test(cfg.h, cfg.Inputs...)
}
