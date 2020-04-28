// SPDX-License-Identifier: MIT

// Package mock 根据 doc 生成 mock 数据
package mock

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/issue9/mux/v2"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/ast"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/vars"
)

const (
	disableHead    = false
	disableOptions = false
)

// Mock 管理 mock 数据
type Mock struct {
	h       *core.MessageHandler
	doc     *ast.APIDoc
	mux     *mux.Mux
	servers map[string]string
}

// New 声明 Mock 对象
//
// h 用于处理各类输出消息，仅在 ServeHTTP 中的消息才输出到 h；
// d doc.APIDoc 实例，调用方需要保证该数据类型的正确性；
// servers 用于指定 d.Servers 中每一个服务对应的路由前缀
func New(h *core.MessageHandler, d *ast.APIDoc, servers map[string]string) (http.Handler, error) {
	c, err := version.SemVerCompatible(d.APIDoc.V(), vars.Version())
	if err != nil {
		return nil, err
	}
	if !c {
		return nil, locale.Errorf(locale.VersionInCompatible)
	}

	m := &Mock{
		h:       h,
		doc:     d,
		mux:     mux.New(disableOptions, disableHead, true, nil, nil),
		servers: servers,
	}

	if err := m.parse(); err != nil {
		return nil, err
	}

	return m, nil
}

// Load 从本地或是远程加载文档内容
func Load(h *core.MessageHandler, path core.URI, servers map[string]string) (http.Handler, error) {
	data, err := path.ReadAll(nil)
	if err != nil {
		return nil, err
	}

	p := core.Position{}
	var offset int
	for {
		r, size := utf8.DecodeRune(data[offset:])
		if size == 0 {
			break
		}
		if r == utf8.RuneError && size == 1 {
			loc := core.Location{
				URI: path,
				Range: core.Range{
					End: p,
				},
			}
			return nil, core.NewLocaleError(loc, "", locale.ErrInvalidUTF8Character)
		}

		offset += size
		p.Character++
		if r == '\n' {
			p.Line++
			p.Character = 0
		}
	}

	loc := core.Location{
		URI: path,
		Range: core.Range{
			End: p,
		},
	}

	// 加载并验证
	d := &ast.APIDoc{}
	if err = d.Parse(core.Block{Location: loc, Data: data}); err != nil {
		return nil, err
	}

	return New(h, d, servers)
}

func (m *Mock) parse() error {
	for _, api := range m.doc.Apis {
		handler := m.buildAPI(api)

		if len(api.Servers) == 0 {
			err := m.mux.Handle(api.Path.Path.V(), handler, api.Method.V())
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

			err := prefix.Handle(api.Path.Path.V(), handler, api.Method.V())
			if err != nil {
				return err
			}
		}
	}

	for path, methods := range m.mux.All(disableHead, disableOptions) {
		m.h.Message(core.Info, locale.LoadAPI, path, strings.Join(methods, ","))
	}

	return nil
}

func (m *Mock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func hasServer(tags []*ast.Element, key string) bool {
	for _, tag := range tags {
		if key == tag.Content.Value {
			return true
		}
	}

	return false
}
