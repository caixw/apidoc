// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package syslocale 获取所在系统的本地化语言信息。
//
// *nix 系统会使用 LANG 环境变量中的值，windows 在 LANG
// 环境变量不存在的情况下，调用 GetUserDefaultLocaleName 函数获取。
package syslocale

import (
	"os"
	"strings"

	"golang.org/x/text/language"
)

// Get 返回当前系统的本地他信息
func Get() (language.Tag, error) {
	name, err := getLocaleName()
	if err != nil {
		return language.Und, err
	}

	return language.Parse(name)
}

// 获取环境变量 LANG 中有关本地化信息的值。
func getEnvLang() string {
	name := os.Getenv("LANG")

	// LANG = zh_CN.UTF-8 过滤掉最后的编码方式
	index := strings.LastIndexByte(name, '.')
	if index > 0 {
		name = name[:index]
	}

	return name
}
