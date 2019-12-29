// SPDX-License-Identifier: MIT

// Package mock 根据 doc 生成 mock 数据
package mock

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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

// Load 从本地或是远程加载文档内容
func Load(h *message.Handler, path string, servers map[string]string) (http.Handler, error) {
	isURL := strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")

	if isURL {
		return LoadFromURL(h, path, servers)
	}
	return LoadFromPath(h, path, servers)
}

// LoadFromPath 加载 XML 文档用以初始化 Mock 对象
func LoadFromPath(h *message.Handler, path string, servers map[string]string) (http.Handler, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return loadContent(h, path, r, servers)
}

// LoadFromURL 从远程 URL 加载文档并初始化为 Mock 对象
func LoadFromURL(h *message.Handler, url string, servers map[string]string) (http.Handler, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return loadContent(h, url, resp.Body, servers)
}

// path 仅用于定位错误，内容存在于 r 中。
func loadContent(h *message.Handler, path string, r io.ReadCloser, servers map[string]string) (http.Handler, error) {
	defer r.Close()

	data, err := ioutil.ReadAll(r)
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

func (m *Mock) parse() error {
	for _, api := range m.doc.Apis {
		handler := m.buildAPI(api)

		if len(api.Servers) == 0 {
			err := m.mux.Handle(api.Path.Path, handler, string(api.Method))
			if err != nil {
				return err
			}
			continue
		}

		for name, prefix := range m.servers {
			prefix := m.mux.Prefix(prefix)

			if !hasServer(api.Servers, name) {
				continue
			}

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
