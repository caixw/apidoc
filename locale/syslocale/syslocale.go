// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package syslocale 获取所在系统的本地化信息
package syslocale

import (
	"os"
	"strings"
)

// Get 返回当前系统的本地他信息
func Get() (string, error) {
	return getLocaleName()
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
