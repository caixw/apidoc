// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openapi

import (
	"strings"

	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/internal/locale"
)

// 去掉 URL 中的 {} 模板参数。使其符合 is.URL 的判断规则
var urlreplace = strings.NewReplacer("{", "", "}", "")

// Server 服务器描述信息
type Server struct {
	URL         string                     `json:"url" yaml:"url"`
	Description Description                `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// ServerVariable Server 中 URL 模板中对应的参数变量值
type ServerVariable struct {
	Enum        []string    `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string      `json:"default" yaml:"default"`
	Description Description `json:"description,omitempty" yaml:"description,omitempty"`
}

func newServer(srv *doc.Server) *Server {
	return &Server{
		URL:         srv.URL,
		Description: srv.Description,
	}
}

// Sanitize 数据检测
func (srv *Server) Sanitize() *Error {
	url := urlreplace.Replace(srv.URL)
	if url == "" { // 可以是 / 未必是一个 URL
		return newError("url", locale.Sprintf(locale.ErrRequired))
	}

	for key, val := range srv.Variables {
		if err := val.Sanitize(); err != nil {
			err.Field = "variables[" + key + "]." + err.Field
			return err
		}

		k := "{" + key + "}"
		if strings.Index(srv.URL, k) < 0 {
			return newError("variables["+key+"]", locale.Sprintf(locale.ErrInvalidValue))
		}
	}

	return nil
}

// Sanitize 数据检测
func (v *ServerVariable) Sanitize() *Error {
	if v.Default == "" {
		return newError("default", locale.Sprintf(locale.ErrRequired))
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
		return newError("default", locale.Sprintf(locale.ErrInvalidValue))
	}

	return nil
}
