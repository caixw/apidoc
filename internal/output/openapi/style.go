// SPDX-License-Identifier: MIT

package openapi

import (
	"github.com/caixw/apidoc/v5/errors"
	"github.com/caixw/apidoc/v5/internal/locale"
)

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
func (style *Style) Sanitize() *errors.Error {
	switch style.Style {
	case StyleMatrix, StyleLabel, StyleForm, StyleSimple, StyleSpaceDelimited, StylePipeDelimited, StyleDeepObject:
	default:
		return errors.New("", "style", 0, locale.ErrInvalidValue)
	}

	return nil
}
