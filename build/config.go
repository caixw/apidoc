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
	"github.com/caixw/apidoc/v7/internal/vars"
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
	wd core.URI

	h *core.MessageHandler
}

// LoadConfig 加载指定目录下的配置文件
//
// 所有的错误信息会输出到 h，在出错时，会返回 nil
func LoadConfig(h *core.MessageHandler, wd core.URI) *Config {
	for _, filename := range vars.AllowConfigFilenames {
		path := wd.Append(filename)

		exists, err := path.Exists()
		if err != nil {
			h.Error(core.Erro, err)
			continue
		} else if exists {
			cfg, err := loadFile(wd, path)
			if err != nil {
				h.Error(core.Erro, err)
				return nil
			}
			cfg.h = h
			return cfg
		}
	}

	msg := core.NewSyntaxError(core.Location{}, wd.Append(vars.AllowConfigFilenames[0]).String(), locale.ErrRequired)
	h.Error(core.Erro, msg)
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

		if err := i.Sanitize(); err != nil {
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
	return cfg.Output.Sanitize()
}

// SaveToFile 将内容保存至文件
func (cfg *Config) SaveToFile(path core.URI) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return path.WriteAll(data)
}

// Build 解析文档并输出文档内容
//
// 具体信息可参考 Build 函数的相关文档。
func (cfg *Config) Build(start time.Time) {
	if err := Build(cfg.h, cfg.Output, cfg.Inputs...); err != nil {
		cfg.h.Error(core.Erro, err)
		return
	}

	cfg.h.Locale(core.Succ, locale.Complete, cfg.Output.Path, time.Now().Sub(start))
}

// Buffer 根据 wd 目录下的配置文件生成文档内容并保存至内存
//
// 具体信息可参考 Buffer 函数的相关文档。
func (cfg *Config) Buffer() *bytes.Buffer {
	buf, err := Buffer(cfg.h, cfg.Output, cfg.Inputs...)
	if err != nil {
		cfg.h.Error(core.Erro, err)
		return nil
	}

	return buf
}

// Test 执行对语法内容的测试
func (cfg *Config) Test() {
	Test(cfg.h, cfg.Inputs...)
}
