// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openapi

// Style.Style 的可选值
const (
	StyleMatrix         = "matrix"
	StyleLabel          = "label"
	StyleForm           = "form"
	StyleSimple         = "simple"
	StyleSpaceDelimited = "spaceDelimited"
	StylePipeDelimited  = "pipeDelimited"
	StyleDeepObject     = "deepObject"
)

// Style 民法风格的相关定义
//
// 不直接作用于对象，被部分对象包含，比如 Encoding 和 Parameter 等
type Style struct {
	Style         string `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       bool   `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool   `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}

// Sanitize 对数据进行验证
func (style *Style) Sanitize() *Error {
	switch style.Style {
	case StyleMatrix, StyleLabel, StyleForm, StyleSimple, StyleSpaceDelimited, StylePipeDelimited, StyleDeepObject:
	default:
		return newError("style", "无效的值")
	}

	return nil
}
