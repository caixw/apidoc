// SPDX-License-Identifier: MIT

// Package build 提供构建文档的相关功能
package build

import (
	"bytes"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// Build 解析文档并输出文档内容
func Build(h *core.MessageHandler, o *Output, i ...*Input) error {
	d, err := parse(h, i...)
	if err != nil {
		return err
	}
	if err = o.sanitize(); err != nil {
		return err
	}

	buf, err := o.buffer(d)
	if err != nil {
		return err
	}

	return o.Path.WriteAll(buf.Bytes())
}

// Buffer 生成文档内容并返回
func Buffer(h *core.MessageHandler, o *Output, i ...*Input) (*bytes.Buffer, error) {
	d, err := parse(h, i...)
	if err != nil {
		return nil, err
	}
	if err = o.sanitize(); err != nil {
		return nil, err
	}

	return o.buffer(d)
}

// CheckSyntax 测试文档语法，并将结果输出到 h
func CheckSyntax(h *core.MessageHandler, i ...*Input) {
	if _, err := parse(h, i...); err != nil {
		h.Error(err)
		return
	}
	h.Locale(core.Succ, locale.TestSuccess)
}

func parse(h *core.MessageHandler, i ...*Input) (*ast.APIDoc, error) {
	for _, item := range i {
		if err := item.sanitize(); err != nil {
			return nil, err
		}
	}

	d := &ast.APIDoc{}
	d.ParseBlocks(h, func(blocks chan core.Block) {
		ParseInputs(blocks, h, i...)
	})

	return d, nil
}
