// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package locale 提供了一个本地化翻译服务。
//
// NOTE: locale 包作为一个最底层的功能实现，不应该依赖
// 程序中其它任何包，它们都有可能调用 locale 包中的相关内容。
package locale

import (
	"errors"
	"os"
	"strings"

	"github.com/caixw/apidoc/vars"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Init 初始化 locale 包并。
// 无论是否返回错误信息，都会初始一种语言作为其交互语言。
func Init() error {
	tag, err := getTag()
	localePrinter = message.NewPrinter(tag)

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

	// 此条必定成功，因为与 vars.DefaultLocale 相同的值已经在上面的 for 特环中执行过。
	defaultLocaleTag := language.MustParse(vars.DefaultLocale)

	localeName, err := getLocaleName()
	if err != nil {
		return defaultLocaleTag, err
	}

	// 成功获取了用户的语言信息，但无法解析成 language.Tag 类型
	tag, err := language.Parse(localeName)
	if err != nil {
		return defaultLocaleTag, err
	}

	return tag, nil
}

// NewPrinter 根据 tag 生成一个新的语言输出环境
func NewPrinter(tag language.Tag) *message.Printer {
	return message.NewPrinter(tag)
}

// 获取环境变量 LANG
func getEnvLang() string {
	name := os.Getenv("LANG")

	// LANG = zh_CN.UTF-8 过滤掉最后的编码方式
	index := strings.LastIndexByte(name, '.')
	if index > 0 {
		name = name[:index]
	}

	return name
}
