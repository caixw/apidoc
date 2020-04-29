// SPDX-License-Identifier: MIT

// Package spec 对文档规则的定义
package spec

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/core"
)

const (
	// Version 文档规范的版本
	Version = "6.0.0"

	// MajorVersion 文档规范的主版本信息
	MajorVersion = "v6"
)

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

func newSyntaxError(loc core.Location, field string, key message.Reference, val ...interface{}) error {
	return core.NewLocaleError(loc, field, key, val...)
}
