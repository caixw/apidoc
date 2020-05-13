// SPDX-License-Identifier: MIT

// Package vars 提供了一些公共的函数、结构体及代码级别的设置项。
package vars

// 一些公用的常量
const (
	// 程序的正式名称
	Name = "apidoc"

	// 源码仓库地址
	RepoURL = "https://github.com/caixw/apidoc"

	// 官网
	OfficialURL = "https://apidoc.tools"

	// 默认的本地化语言 ID
	//
	// 当未调用相关函数设置 ID，或是设置为一个不支持的 ID 时，
	// 系统最终会采用此 ID。
	//
	// NOTE: 注意大小写需要与 internal/locale 的相同。
	DefaultLocaleID = "cmn-Hans"
)

// AllowConfigFilenames 允许的配置文件名
//
// apidoc 会按此列表的顺序在目录中查找配置文件，
// 直到找到第一个相符的文件，或是在没有时出错。
//
// 在生成配置文件时，会直接拿第一个元素的值作为文件名。
var AllowConfigFilenames = []string{
	".apidoc.yaml",
	".apidoc.yml",
}
