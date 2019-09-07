// SPDX-License-Identifier: MIT

// Package options 输入和输出的配置项
package options

// 文档类型定义
const (
	ApidocXML   = "apidoc+xml"
	OpenapiJSON = "openapi+json"
	OpenapiYAML = "openapi+yaml"
	RAMLYAML    = "raml+yaml"
)

// Output 指定了渲染输出的相关设置项。
type Output struct {
	// 文档的保存路径，建议使用绝对路径。
	Path string `yaml:"path,omitempty"`

	// 输出类型
	Type string `yaml:"type,omitempty"`

	// 只输出该标签的文档，若为空，则表示所有。
	Tags []string `yaml:"tags,omitempty"`
}
