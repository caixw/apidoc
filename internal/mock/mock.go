// SPDX-License-Identifier: MIT

// Package mock 根据 doc 生成 mock 数据
package mock

import (
	"io/ioutil"
	"net/http"

	"github.com/issue9/mux/v2"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/message"
)

// Mock 管理 mock 数据
type Mock struct {
	h       *message.Handler
	doc     *doc.Doc
	prefix  string // 生成的所有路由中的路径前缀
	mux     *mux.Mux
	servers map[string]string
}

// New 声明 Mock 对象
//
// h 用于处理各类输出消息，仅在 ServeHTTP 中的消息才输出到 h；
// prefix 所有路由的前缀；
// servers 用于指定 d.Server 中每一个服务对应的路由前缀
func New(h *message.Handler, d *doc.Doc, prefix string, servers map[string]string) (*Mock, error) {
	m := &Mock{
		h:       h,
		doc:     d,
		prefix:  prefix,
		mux:     mux.New(false, false, true, nil, nil),
		servers: servers,
	}

	if err := m.parse(); err != nil {
		return nil, err
	}

	return m, nil
}

// NewWithPath 加载 XML 文档用以初始化 Mock 对象
func NewWithPath(h *message.Handler, path, prefix string, servers map[string]string) (*Mock, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	d := doc.New()
	if err = d.FromXML(path, 0, data); err != nil {
		return nil, err
	}

	return New(h, d, prefix, servers)
}

func (m *Mock) parse() error {
	p := m.mux.Prefix(m.prefix)

	for name, prefix := range m.servers {
		prefix := p.Prefix(prefix)

		for _, api := range m.doc.Apis {
			if !hasServer(api.Servers, name) {
				continue
			}

			handler := m.build(api)
			err := prefix.Handle(api.Path.Path, handler, string(api.Method))
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

func hasServer(tags []string, key string) bool {
	for _, tag := range tags {
		if key == tag {
			return true
		}
	}

	return false
}
