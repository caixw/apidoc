// SPDX-License-Identifier: MIT

// Package build 提供构建文档的相关功能
package build

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/message"
	"github.com/caixw/apidoc/v6/spec"
)

// Build 解析文档并输出文档内容
func Build(h *message.Handler, o *Output, i ...*Input) error {
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

	return ioutil.WriteFile(o.Path, buf.Bytes(), os.ModePerm)
}

// Buffer 生成文档内容并返回
func Buffer(h *message.Handler, o *Output, i ...*Input) (*bytes.Buffer, error) {
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
func Test(h *message.Handler, i ...*Input) {
	if _, err := parse(h, i...); err != nil {
		h.Error(message.Erro, err)
		return
	}
	h.Message(message.Succ, locale.TestSuccess)
}

func parse(h *message.Handler, i ...*Input) (*spec.APIDoc, error) {
	for _, item := range i {
		if err := item.Sanitize(); err != nil {
			return nil, err
		}
	}

	d := spec.New()
	Parse(d, h, i...)

	if err := d.Sanitize(); err != nil {
		h.Error(message.Erro, err)
		return nil, err
	}

	return d, nil
}

// Parse 分析从 o 中获取的代码块
//
// 所有与解析有关的错误均通过 h 输出。
// 如果是配置文件的错误，则通过 error 返回
func Parse(doc *spec.APIDoc, h *message.Handler, o ...*Input) {
	parse2APIDoc(doc, h, func(blocks chan spec.Block) {
		parseInputs(blocks, h, o...)
	})
}

// ParseFile 分析 path 的内容，并将其中的文档解析至 doc
func ParseFile(doc *spec.APIDoc, h *message.Handler, path string, o *Input) {
	parse2APIDoc(doc, h, func(blocks chan spec.Block) {
		o.parseFile(blocks, h, path)
	})
}

func parse2APIDoc(doc *spec.APIDoc, h *message.Handler, g func(chan spec.Block)) {
	done := make(chan struct{})
	blocks := make(chan spec.Block, 50)

	go func() {
		for block := range blocks {
			if err := doc.ParseBlock(&block); err != nil {
				h.Error(message.Erro, err)
			}
		}
		done <- struct{}{}
	}()

	g(blocks)
	close(blocks)
	<-done
}
