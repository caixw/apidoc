// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openapi

// Parameter.IN 的可选值
const (
	ParameterINPath   = "path"
	ParameterINQuery  = "query"
	ParameterINHeader = "header"
	ParameterINcookie = "cookie"
)

// Header 即 Parameter 的别名，但 Name 字段必须存在。
type Header Parameter

// Parameter 参数信息
// 可同时作用于路径参数、请求参数、报头内容和 Cookie 值。
type Parameter struct {
	Style
	Name            string                `json:"name,omitempty" yaml:"name,omitempty"`
	IN              string                `json:"in,omitempty" yaml:"in,omitempty"`
	Description     Description           `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool                  `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                  `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Schema          *Schema               `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example         ExampleValue          `json:"example,omitempty" yaml:"example,omitempty"`
	Examples        map[string]*Example   `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content         map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}

// Sanitize 对数据进行验证
func (p *Parameter) Sanitize() *Error {
	if err := p.Style.Sanitize(); err != nil {
		return err
	}

	switch p.IN {
	case ParameterINcookie, ParameterINHeader, ParameterINPath, ParameterINQuery:
	default:
		return newError("in", "无效的值")
	}

	// TODO 其它字段检测

	return nil
}

// Sanitize 对数据进行验证
func (h *Header) Sanitize() *Error {
	if err := h.Style.Sanitize(); err != nil {
		return err
	}

	if h.IN != "" {
		return newError("in", "只能为空")
	}

	if h.Name != "" {
		return newError("name", "只能为空")
	}

	return nil
}
