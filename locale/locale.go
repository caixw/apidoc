// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Init 初始化 locale 包。
// defaultTag 默认的语言
func Init(defaultLang, lang string) {

	messages, found := locales[defaultLang]
	if !found {
		panic("参数 defaultTag 所指的语言不存在")
	}
	locales["und"] = messages

	for id, messages := range locales {
		tag := language.MustParse(id)
		for key, val := range messages {
			message.SetString(tag, key, val)
		}
	}

	tag := language.MustParse(lang)
	localePrinter = message.NewPrinter(tag)
	if localePrinter == nil {
		panic(fmt.Errorf("无法获取指定语言[%v]的相关翻译内容", tag))
	}
}
