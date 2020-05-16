// SPDX-License-Identifier: MIT

// Package build 提供构建文档的相关功能
//
// 多行注释和单行注释在处理上会有一定区别：
//  - 单行注释，风格相同且相邻的注释会被合并成一个注释块；
//  - 单行注释，风格不相同且相邻的注释会被按注释风格多个注释块；
//  - 多行注释，即使两个注释释块相邻也会被分成两个注释块来处理。
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
	if err = o.Sanitize(); err != nil {
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
	if err = o.Sanitize(); err != nil {
		return nil, err
	}

	return o.buffer(d)
}

// Test 测试文档语法，并将结果输出到 h
func Test(h *core.MessageHandler, i ...*Input) {
	if _, err := parse(h, i...); err != nil {
		h.Error(core.Erro, err)
		return
	}
	h.Locale(core.Succ, locale.TestSuccess)
}

func parse(h *core.MessageHandler, i ...*Input) (*ast.APIDoc, error) {
	for _, item := range i {
		if err := item.Sanitize(); err != nil {
			return nil, err
		}
	}

	d := &ast.APIDoc{}
	d.ParseBlocks(h, func(blocks chan core.Block) {
		parseInputs(blocks, h, i...)
	})

	return d, nil
}
