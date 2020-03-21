// SPDX-License-Identifier: MIT

// Package spec 对文档规则的定义
package spec

import (
	xmessage "golang.org/x/text/message"

	"github.com/caixw/apidoc/v6/core"
)

const (
	// Version 文档规范的版本
	Version = "6.0.0"

	// MajorVersion 文档规范的主版本信息
	MajorVersion = "v6"
)

// Block 表示原始的注释代码块
type Block struct {
	Location core.Location

	// Raw 表示原始的注释代码内容
	//
	// Data 为处理之后的数据
	// 为一个正常的 XML 格式内容，且长度应该与 Raw 相同。
	Raw  []byte
	Data []byte
}

// 返回基于当前范围的错误信息
func (b *Block) localeError(field string, key xmessage.Reference, v ...interface{}) error {
	return core.NewLocaleError(b.Location, field, key, v...)
}

func fixedSyntaxError(loc core.Location, err error, field string) error {
	if serr, ok := err.(*core.SyntaxError); ok {
		serr.Location = loc

		if serr.Field == "" {
			serr.Field = field
		} else {
			serr.Field = field + serr.Field
		}
		return err
	}

	return core.WithError(loc, field, err)
}

func newSyntaxError(loc core.Location, field string, key xmessage.Reference, val ...interface{}) error {
	return core.NewLocaleError(loc, field, key, val...)
}
