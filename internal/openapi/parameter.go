// SPDX-License-Identifier: MIT

package openapi

import (
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

// Parameter.IN 的可选值
const (
	ParameterINPath   = "path"
	ParameterINQuery  = "query"
	ParameterINHeader = "header"
	ParameterINCookie = "cookie"
)

// Header 即 Parameter 的别名，但 Name 字段必须不能存在。
type Header Parameter

// Parameter 参数信息
// 可同时作用于路径参数、请求参数、报头内容和 Cookie 值。
type Parameter struct {
	Style
	Name            string                `json:"name,omitempty" yaml:"name,omitempty"`
	IN              string                `json:"in,omitempty" yaml:"in,omitempty"`
	Description     string                `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool                  `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                  `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Schema          *Schema               `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example         ExampleValue          `json:"example,omitempty" yaml:"example,omitempty"`
	Examples        map[string]*Example   `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content         map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`

	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

// Sanitize 对数据进行验证
func (p *Parameter) Sanitize() *message.SyntaxError {
	if err := p.Style.Sanitize(); err != nil {
		return err
	}

	switch p.IN {
	case ParameterINCookie, ParameterINHeader, ParameterINPath, ParameterINQuery:
	default:
		return message.NewLocaleError("", "in", 0, locale.ErrInvalidValue)
	}

	return nil
}

// Sanitize 对数据进行验证
func (h *Header) Sanitize() *message.SyntaxError {
	if err := h.Style.Sanitize(); err != nil {
		return err
	}

	if h.IN != "" {
		return message.NewLocaleError("", "in", 0, locale.ErrInvalidValue)
	}

	if h.Name != "" {
		return message.NewLocaleError("", "name", 0, locale.ErrInvalidValue)
	}

	return nil
}
