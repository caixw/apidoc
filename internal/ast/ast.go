// SPDX-License-Identifier: MIT

// Package ast 为 xml 服务的抽象语法树
package ast

const (
	// Version 文档规范的版本
	Version = "6.1.0"

	// MajorVersion 文档规范的主版本信息
	MajorVersion = "v6"

	// 文档允许的最小长度
	//
	// 文档都是以 <api 或是 <apidoc 开头的，所以最起码要等于此值。
	minSize = len("<api/>")
)
