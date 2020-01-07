// SPDX-License-Identifier: MIT

package openapi

import (
	"strings"

	"github.com/caixw/apidoc/v6/doc"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/message"
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

func newServer(srv *doc.Server) *Server {
	desc := srv.Summary
	if srv.Description.Text != "" {
		desc = srv.Description.Text
	}

	return &Server{
		URL:         srv.URL,
		Description: desc,
	}
}

func (srv *Server) sanitize() *message.SyntaxError {
	url := urlreplace.Replace(srv.URL)
	if url == "" { // 可以是 / 未必是一个 URL
		return message.NewLocaleError("", "url", 0, locale.ErrRequired)
	}

	for key, val := range srv.Variables {
		if err := val.sanitize(); err != nil {
			err.Field = "variables[" + key + "]." + err.Field
			return err
		}

		k := "{" + key + "}"
		if strings.Index(srv.URL, k) < 0 {
			return message.NewLocaleError("", "variables["+key+"]", 0, locale.ErrInvalidValue)
		}
	}

	return nil
}

func (v *ServerVariable) sanitize() *message.SyntaxError {
	if v.Default == "" {
		return message.NewLocaleError("", "default", 0, locale.ErrRequired)
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
		return message.NewLocaleError("", "default", 0, locale.ErrInvalidValue)
	}

	return nil
}
