// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package apidoc

import (
	"golang.org/x/text/language"

	"github.com/caixw/apidoc/internal/locale"
)

// InitLocale 初始化语言环境
//
// NOTE: 必须保证在第一时间调用。
func InitLocale(tag language.Tag) error {
	return locale.Init(tag)
}
