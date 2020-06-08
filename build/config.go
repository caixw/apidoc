// SPDX-License-Identifier: MIT

package build

import (
	"bytes"
	"strconv"
	"time"

	"github.com/issue9/version"
	"gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// 允许的配置文件名
//
// 会按此列表的顺序在目录中查找配置文件，
// 直到找到第一个相符的文件，或是在没有时出错。
//
// 在生成配置文件时，会直接拿第一个元素的值作为文件名。
var allowConfigFilenames = []string{
	".apidoc.yaml",
	".apidoc.yml",
}

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
	wd core.URI

	h *core.MessageHandler
}

// LoadConfig 加载指定目录下的配置文件
//
// 所有的错误信息会输出到 h，在出错时，会返回 nil
func LoadConfig(h *core.MessageHandler, wd core.URI) *Config {
	for _, filename := range allowConfigFilenames {
		path := wd.Append(filename)
		if exists, err := path.Exists(); err != nil {
			h.Error(err)
			continue
		} else if exists {
			cfg, err := loadFile(wd, path)
			if err != nil {
				h.Error(err)
				return nil
			}
			cfg.h = h
			return cfg
		}
	}

	msg := core.NewSyntaxError(core.Location{}, wd.Append(allowConfigFilenames[0]).String(), locale.ErrRequired)
	h.Error(msg)
	return nil
}

func loadFile(wd, path core.URI) (*Config, error) {
	data, err := path.ReadAll(nil)
	if err != nil {
		return nil, core.NewSyntaxErrorWithError(core.Location{URI: path}, "", err)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, core.NewSyntaxErrorWithError(core.Location{URI: path}, "", err)
	}
	cfg.wd = wd

	if err := cfg.sanitize(path); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) sanitize(file core.URI) error {
	// 比较版本号兼容问题
	compatible, err := version.SemVerCompatible(ast.Version, cfg.Version)
	if err != nil {
		return core.NewSyntaxErrorWithError(core.Location{URI: file}, "version", err)
	}
	if !compatible {
		return core.NewSyntaxError(core.Location{URI: file}, "version", locale.VersionInCompatible)
	}

	if len(cfg.Inputs) == 0 {
		return core.NewSyntaxError(core.Location{URI: file}, "inputs", locale.ErrRequired)
	}

	if cfg.Output == nil {
		return core.NewSyntaxError(core.Location{URI: file}, "output", locale.ErrRequired)
	}

	for index, i := range cfg.Inputs {
		field := "inputs[" + strconv.Itoa(index) + "]"

		if i.Dir, err = abs(i.Dir, cfg.wd); err != nil {
			return core.NewSyntaxErrorWithError(core.Location{URI: file}, field+".path", err)
		}

		if err := i.sanitize(); err != nil {
			if serr, ok := err.(*core.SyntaxError); ok {
				serr.Location.URI = file
				serr.Field = field + serr.Field
			}
			return err
		}
	}

	if cfg.Output.Path, err = abs(cfg.Output.Path, cfg.wd); err != nil {
		return core.NewSyntaxErrorWithError(core.Location{URI: file}, "output.path", err)
	}
	return cfg.Output.sanitize()
}

// Save 将内容保存至 wd 目录下的 .apidoc.yaml 文件
func (cfg *Config) Save(wd core.URI) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return wd.Append(allowConfigFilenames[0]).WriteAll(data)
}

// Build 解析文档并输出文档内容
//
// 具体信息可参考 Build 函数的相关文档。
func (cfg *Config) Build(start time.Time) {
	if err := Build(cfg.h, cfg.Output, cfg.Inputs...); err != nil {
		cfg.h.Error(err)
		return
	}

	// 即使部分解析出错，只要有部分内容保存，就会输出此信息。
	cfg.h.Locale(core.Info, locale.Complete, cfg.Output.Path, time.Now().Sub(start))
}

// Buffer 根据 wd 目录下的配置文件生成文档内容并保存至内存
//
// 具体信息可参考 Buffer 函数的相关文档。
func (cfg *Config) Buffer() *bytes.Buffer {
	buf, err := Buffer(cfg.h, cfg.Output, cfg.Inputs...)
	if err != nil {
		cfg.h.Error(err)
		return nil
	}

	return buf
}

// CheckSyntax 执行对语法内容的测试
func (cfg *Config) CheckSyntax() {
	CheckSyntax(cfg.h, cfg.Inputs...)
}
