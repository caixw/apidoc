// SPDX-License-Identifier: MIT

// Package mock 根据 doc 生成 mock 数据
package mock

import (
	"io/ioutil"

	"github.com/caixw/apidoc/v5/doc"
)

// Mock 管理 mock 数据
type Mock struct {
	doc *doc.Doc
}

// New 声明 Mock 对象
func New(d *doc.Doc) *Mock {
	return &Mock{
		doc: d,
	}
}

// NewWithPath 加 xml 文档用以初始化 Mock 对象
func NewWithPath(path string) (*Mock, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	d := doc.New()
	if err = d.FromXML(path, 0, data); err != nil {
		return nil, err
	}

	return New(d), nil
}

// Serve 执行 HTTP 服务
func (m *Mock) Serve() {
	// TODO
}
