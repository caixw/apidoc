// SPDX-License-Identifier: MIT

// Package ast 为 xml 服务的抽象语法树
package ast

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v6/core"
)

const (
	// Version 文档规范的版本
	Version = "6.0.0"

	// MajorVersion 文档规范的主版本信息
	MajorVersion = "v6"

	// 文档允许的最小长度
	//
	// 文档都是以 <api 或是 <apidoc 开头的，所以最起码要大于 len("<api/>") 的值。
	minSize = 6
)

func newError(r core.Range, key message.Reference, v ...interface{}) error {
	return core.NewLocaleError(core.Location{Range: r}, "", key, v...)
}

func withError(r core.Range, err error) error {
	return core.WithError(core.Location{Range: r}, "", err)
}
