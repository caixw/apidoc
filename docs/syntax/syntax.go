// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package syntax 提供简易的词法分析工具
package syntax

import "unicode"

// Enum 分析枚举内容
func Enum() ([]string, error) {
	// TODO

	return nil, nil
}

// 分隔成指定大小的字符串数组
func split(data []byte, size int) [][]byte {
	ret := make([][]byte, 0, size)
	start := 0
	pos := 0
	isspace := true // 前一个字符是否为空白字符

	for ; ; pos++ {
		switch {
		case pos >= len(data): // EOF
			return append(ret, data[start:])
		case unicode.IsSpace(rune(data[pos])):
			if !isspace {
				ret = append(ret, data[start:pos])
				start = pos
				isspace = true
			}
		default:
			if isspace {
				if len(ret) >= size-1 {
					return append(ret, data[pos:])
				}

				start = pos
				isspace = false
			}
		}
	}
}
