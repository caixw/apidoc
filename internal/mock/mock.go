// SPDX-License-Identifier: MIT

// Package mock 根据 doc 生成 mock 数据
package mock

import (
	"io/ioutil"
	"net/http"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/issue9/mux/v2"
)

// Mock 管理 mock 数据
type Mock struct {
	doc     *doc.Doc
	mux     *mux.Mux
	servers map[string]string
}

// New 声明 Mock 对象
func New(d *doc.Doc, prefix string, servers map[string]string) (*Mock, error) {
	m := &Mock{
		doc:     d,
		mux:     mux.New(false, false, true, nil, nil),
		servers: servers,
	}

	if err := m.parse(); err != nil {
		return nil, err
	}

	return m, nil
}

// NewWithPath 加 xml 文档用以初始化 Mock 对象
func NewWithPath(path, prefix string, servers map[string]string) (*Mock, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	d := doc.New()
	if err = d.FromXML(path, 0, data); err != nil {
		return nil, err
	}

	return New(d, prefix, servers)
}

func (m *Mock) parse() error {
	for name, prefix := range m.servers {
		prefix := m.mux.Prefix(prefix)

		for _, api := range m.doc.Apis {
			if !hasTag(api.Tags, name) {
				continue
			}

			err := prefix.Handle(api.Path.Path, build(api), string(api.Method))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Mock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func hasTag(tags []string, key string) bool {
	for _, tag := range tags {
		if key == tag {
			return true
		}
	}

	return false
}
