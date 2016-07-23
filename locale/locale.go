// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Init 初始化 locale 包。
// defaultTag 默认的语言
func Init(defaultTag string) {
	tag, found := locales[defaultTag]
	if !found {
		panic("参数 defaultTag 所指的语言不存在")
	}
	locales["und"] = tag

	for id, messages := range locales {
		tag := language.MustParse(id)
		for key, val := range messages {
			message.SetString(tag, key, val)
		}
	}
}
