// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package options 输入和输出的配置项
package options

// Input 指定输入内容的相关信息。
type Input struct {
	Lang      string   `yaml:"lang"`               // 输入的目标语言
	Dir       string   `yaml:"dir"`                // 源代码目录
	Exts      []string `yaml:"exts,omitempty"`     // 需要扫描的文件扩展名，若未指定，则使用默认值
	Recursive bool     `yaml:"recursive"`          // 是否查找 Dir 的子目录
	Encoding  string   `yaml:"encoding,omitempty"` // 文件的编码，为空表示 utf-8
}

// 文档类型定义
const (
	ApidocJSON  = "apidoc+json"
	ApidocYAML  = "apidoc+yaml"
	OpenapiJSON = "openapi+json"
	OpenapiYAML = "openapi+yaml"
	RamlJSON    = "raml+json"
)

// Output 指定了渲染输出的相关设置项。
type Output struct {
	// 文档的保存路径，包含目录和文件名，若为空，则为当前目录下的
	Path string `yaml:"path,omitempty"`

	// 输出类型
	Type string `yaml:"type,omitempty"`

	// 只输出该标签的文档，若为空，则表示所有。
	Tags []string `yaml:"tags,omitempty"`
}
