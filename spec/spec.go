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
	File  string
	Range core.Range

	// Raw 表示原始的注释代码内容
	//
	// Data 为处理之后的数据
	// 为一个正常的 XML 格式内容，且长度应该与 Raw 相同。
	Raw  []byte
	Data []byte
}

func (b *Block) localeError(field string, key xmessage.Reference, v ...interface{}) error {
	return core.NewLocaleError(b.File, field, b.Range.Start.Line, key, v...)
}

func fixedSyntaxError(err error, file, field string, line int) error {
	if serr, ok := err.(*core.SyntaxError); ok {
		serr.File = file
		serr.Line = line

		if serr.Field == "" {
			serr.Field = field
		} else {
			serr.Field = field + serr.Field
		}
		return err
	}

	return core.WithError(file, field, line, err)
}

func newSyntaxError(field string, key xmessage.Reference, val ...interface{}) error {
	return core.NewLocaleError("", field, 0, key, val...)
}
