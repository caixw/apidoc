// SPDX-License-Identifier: MIT

package protocol

import (
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
)

// APIDocOutline 传递给客户端的文档摘要
//
// 这不是一个标准的 LSP 数据结构，由 apidoc 自定义，
// 用户由服务端向客户端发送当前的文档结构信息。
type APIDocOutline struct {
	WorkspaceFolder

	Location core.Location   `json:"location,omitempty"`
	Title    string          `json:"title,omitempty"`
	Version  string          `json:"version,omitempty"`
	Tags     []*APIDocTag    `json:"tags,omitempty"`
	Servers  []*APIDocServer `json:"servers,omitempty"`

	APIs []*API `json:"apis,omitempty"`
}

// APIDocTag 文档支持的标签属性
type APIDocTag struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// APIDocServer 文档支持的服务器
type APIDocServer struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// API 描述单个 API 的信息
type API struct {
	Location   core.Location `json:"location"`
	Method     string        `json:"method"`
	Path       string        `json:"path"`
	Tags       []string      `json:"tags,omitempty"`
	Servers    []string      `json:"servers,omitempty"`
	Deprecated string        `json:"deprecated,omitempty"`
	Summary    string        `json:"summary,omitempty"`
}

// BuildAPIDocOutline 根据 ast.APIDoc 构建 APIDoc
//
// 如果 doc 不是一个有效的文档内容，比如是零值，则返回 nil。
func BuildAPIDocOutline(f WorkspaceFolder, doc *ast.APIDoc) *APIDocOutline {
	if doc == nil || doc.Title.V() == "" {
		return nil
	}

	tags := make([]*APIDocTag, 0, len(doc.Tags))
	for _, t := range doc.Tags {
		tags = append(tags, &APIDocTag{
			ID:    t.Name.V(),
			Title: t.Title.V(),
		})
	}

	servers := make([]*APIDocServer, 0, len(doc.Servers))
	for _, srv := range doc.Servers {
		servers = append(servers, &APIDocServer{
			ID:  srv.Name.V(),
			URL: srv.URL.V(),
		})
	}

	apis := make([]*API, 0, len(doc.APIs))
	for _, api := range doc.APIs {
		uri := api.URI
		if uri == "" {
			uri = doc.URI
		}

		ts := make([]string, 0, len(api.Tags))
		for _, tag := range api.Tags {
			ts = append(ts, tag.V())
		}

		srvs := make([]string, 0, len(api.Servers))
		for _, srv := range api.Servers {
			srvs = append(srvs, srv.V())
		}

		summary := api.Summary.V()
		if summary == "" {
			summary = api.Description.V()
		}

		apis = append(apis, &API{
			Location: core.Location{
				URI:   uri,
				Range: api.Range,
			},
			Method:     api.Method.V(),
			Path:       api.Path.Path.V(),
			Tags:       ts,
			Servers:    srvs,
			Deprecated: api.Description.V(),
			Summary:    summary,
		})
	}

	return &APIDocOutline{
		WorkspaceFolder: f,
		Location: core.Location{
			URI:   doc.URI,
			Range: doc.Range,
		},
		Title:   doc.Title.V(),
		Version: doc.Version.V(),
		Tags:    tags,
		Servers: servers,
		APIs:    apis,
	}
}
