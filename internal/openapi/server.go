// SPDX-License-Identifier: MIT

package openapi

import (
	"strings"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// 去掉 URL 中的 {} 模板参数。使其符合 is.URL 的判断规则
var urlreplace = strings.NewReplacer("{", "", "}", "")

// Server 服务器描述信息
type Server struct {
	URL         string                     `json:"url" yaml:"url"`
	Description string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// ServerVariable Server 中 URL 模板中对应的参数变量值
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `json:"default" yaml:"default"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}

func newServer(srv *ast.Server) *Server {
	desc := srv.Summary.V()
	if srv.Description != nil && srv.Description.Text != nil {
		desc = srv.Description.V()
	}

	return &Server{
		URL:         srv.URL.V(),
		Description: desc,
	}
}

func (srv *Server) sanitize() *core.Error {
	url := urlreplace.Replace(srv.URL)
	if url == "" { // 可以是 / 未必是一个 URL
		return core.NewError(locale.ErrIsEmpty, "url").WithField("url")
	}

	for key, val := range srv.Variables {
		if err := val.sanitize(); err != nil {
			err.Field = "variables[" + key + "]." + err.Field
			return err
		}

		k := "{" + key + "}"
		if !strings.Contains(srv.URL, k) {
			return core.NewError(locale.ErrInvalidValue).WithField("variables[" + key + "]")
		}
	}

	return nil
}

func (v *ServerVariable) sanitize() *core.Error {
	if v.Default == "" {
		return core.NewError(locale.ErrIsEmpty, "default").WithField("default")
	}

	if len(v.Enum) == 0 {
		return nil
	}

	found := false
	for _, item := range v.Enum {
		if item == v.Default {
			found = true
			break
		}
	}

	if !found {
		return core.NewError(locale.ErrInvalidValue).WithField("default")
	}

	return nil
}
