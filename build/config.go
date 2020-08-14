// SPDX-License-Identifier: MIT

package build

import (
	"bytes"
	"os"
	"strconv"

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
}

// LoadConfig 加载指定目录下的配置文件
func LoadConfig(wd core.URI) (*Config, error) {
	// 禁用加载远程配置文件
	//
	// 如果远程配置文件中的目录指向本地，则用户执行相关操作时，
	// 会映射他本地的系统上，这未必是用户想要的结果，且对用户数据不安全。
	if scheme, _ := wd.Parse(); scheme != "" && scheme != core.SchemeFile {
		return nil, locale.NewError(locale.ErrInvalidURIScheme, scheme)
	}

	for _, filename := range allowConfigFilenames {
		path := wd.Append(filename)
		if exists, err := path.Exists(); err != nil {
			continue
		} else if exists {
			cfg, err := loadFile(wd, path)
			if err != nil {
				return nil, err
			}
			return cfg, nil
		}
	}

	field := wd.Append(allowConfigFilenames[0]).String()
	return nil, core.WithError(os.ErrNotExist).WithField(field)
}

func loadFile(wd, path core.URI) (*Config, error) {
	data, err := path.ReadAll(nil)
	if err != nil {
		return nil, (core.Location{URI: path}).WithError(err)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, (core.Location{URI: path}).WithError(err)
	}

	if err := cfg.sanitize(wd); err != nil {
		return nil, err
	}

	return cfg, nil
}

// file 表示出错时的文件定位
func (cfg *Config) sanitize(wd core.URI) error {
	file := wd.Append(allowConfigFilenames[0])

	// 比较版本号兼容问题
	compatible, err := version.SemVerCompatible(ast.Version, cfg.Version)
	if err != nil {
		return (core.Location{URI: file}).WithError(err).WithField("version")
	}
	if !compatible {
		return (core.Location{URI: file}).NewError(locale.VersionInCompatible).WithField("version")
	}

	if len(cfg.Inputs) == 0 {
		return (core.Location{URI: file}).NewError(locale.ErrIsEmpty, "inputs").WithField("inputs")
	}

	if cfg.Output == nil {
		return (core.Location{URI: file}).NewError(locale.ErrIsEmpty, "output").WithField("output")
	}

	for index, i := range cfg.Inputs {
		field := "inputs[" + strconv.Itoa(index) + "]"

		if i == nil {
			return (core.Location{URI: file}).NewError(locale.ErrIsEmpty, field).WithField(field)
		}

		if i.Dir, err = abs(i.Dir, wd); err != nil {
			return (core.Location{URI: file}).WithError(err).WithField(field + ".path")
		}

		if err := i.sanitize(); err != nil {
			if serr, ok := err.(*core.Error); ok {
				serr.Location.URI = file
				serr.Field = field + serr.Field
			}
			return err
		}
	}

	if cfg.Output.Path, err = abs(cfg.Output.Path, wd); err != nil {
		return (core.Location{URI: file}).WithError(err).WithField("output.path")
	}
	return cfg.Output.sanitize()
}

// Save 将内容保存至 wd 目录下的 .apidoc.yaml 文件
//
// 保存时会将各个与路径相关的字段尽量改成与 wd 相关的相对路径。
func (cfg *Config) Save(wd core.URI) (err error) {
	for _, input := range cfg.Inputs { // 调整成相对路径
		if input.Dir, err = rel(input.Dir, wd); err != nil {
			return err
		}
	}

	if cfg.Output.Path != "" { // 调整成相对路径
		if cfg.Output.Path, err = rel(cfg.Output.Path, wd); err != nil {
			return err
		}
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return wd.Append(allowConfigFilenames[0]).WriteAll(data)
}

// Build 解析文档并输出文档内容
//
// 具体信息可参考 Build 函数的相关文档。
func (cfg *Config) Build(h *core.MessageHandler) {
	if err := Build(h, cfg.Output, cfg.Inputs...); err != nil {
		panic(err) // 由 loadConfig 保证配置项的正确，如果还出错则直接 panic
	}
}

// Buffer 根据 wd 目录下的配置文件生成文档内容并保存至内存
//
// 具体信息可参考 Buffer 函数的相关文档。
func (cfg *Config) Buffer(h *core.MessageHandler) *bytes.Buffer {
	buf, err := Buffer(h, cfg.Output, cfg.Inputs...)
	if err != nil {
		panic(err) // 由 loadConfig 保证配置项的正确，如果还出错则直接 panic
	}

	return buf
}

// CheckSyntax 执行对语法内容的测试
func (cfg *Config) CheckSyntax(h *core.MessageHandler) {
	if err := CheckSyntax(h, cfg.Inputs...); err != nil {
		panic(err) // 由 loadConfig 保证配置项的正确，如果还出错则直接 panic
	}
}
