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
	mux     *mux.Mux
	servers map[string]string
}

// New 声明 Mock 对象
//
// h 用于处理各类输出消息，仅在 ServeHTTP 中的消息才输出到 h；
// d doc.Doc 实例，调用方需要保证该数据类型的正确性；
// servers 用于指定 d.Servers 中每一个服务对应的路由前缀
func New(h *message.Handler, d *doc.Doc, servers map[string]string) (http.Handler, error) {
	m := &Mock{
		h:       h,
		doc:     d,
		mux:     mux.New(false, false, true, nil, nil),
		servers: servers,
	}

	if err := m.parse(); err != nil {
		return nil, err
	}

	return m, nil
}

// NewWithPath 加载 XML 文档用以初始化 Mock 对象
func NewWithPath(h *message.Handler, path string, servers map[string]string) (http.Handler, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 加载并验证
	d := doc.New()
	if err = d.FromXML(path, 0, data); err != nil {
		return nil, err
	}

	return New(h, d, servers)
}

// NewWithURL 从远程 URL 加载文档并初始化为 Mock 对象
func NewWithURL(h *message.Handler, url string, servers map[string]string) (http.Handler, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	d := doc.New()
	if err = d.FromXML(url, 0, data); err != nil {
		return nil, err
	}

	return New(h, d, servers)
}

func (m *Mock) parse() error {
	for name, prefix := range m.servers {
		prefix := m.mux.Prefix(prefix)

		for _, api := range m.doc.Apis {
			if !hasServer(api.Servers, name) {
				continue
			}

			handler := m.buildAPI(api)
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
