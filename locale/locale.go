// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package locale 提供了一个本地化翻译服务。
package locale

import (
	"errors"

	"github.com/caixw/apidoc/locale/syslocale"
	"github.com/caixw/apidoc/vars"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// 保证有个初始化的值，部分包的测试功能依赖此变量
var localePrinter *message.Printer = message.NewPrinter(language.MustParse(vars.DefaultLocale))

// Init 初始化 locale 包并。
// 无论是否返回错误信息，都会初始一种语言作为其交互语言。
func Init() error {
	tag, err := getTag()
	localePrinter = NewPrinter(tag)

	return err
}

func getTag() (language.Tag, error) {
	found := false
	for id, messages := range locales { // 保证 locales 已经初始化，即要在 init() 函数之后调用
		tag := language.MustParse(id)
		for key, val := range messages {
			message.SetString(tag, key, val)
		}

		if id == vars.DefaultLocale {
			found = true
		}
	}

	if !found {
		return language.Und, errors.New("vars.DefaultLocale 的值并不存在")
	}

	tag, err := syslocale.Get()
	if err != nil {
		// 此条必定成功，因为与 vars.DefaultLocale 相同的值已经在上面的 for 特环中执行过。
		return language.MustParse(vars.DefaultLocale), err
	}
	return tag, nil
}
