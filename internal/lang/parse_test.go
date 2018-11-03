// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lang

import (
	"strings"
	"testing"

	"github.com/issue9/assert"
)

var (
	code1 = `
# include <stdio.h>
printf("hello world!")
/**
 * 注释代码
 */
`

	code2 = `
# include <stdio.h>
printf("hello world!")
/**
 * 注释代码
	`
)

func TestParse(t *testing.T) {
	a := assert.New(t)

	ret := Parse(nil, nil, nil)
	a.NotNil(ret).
		Equal(0, len(ret))

	ret = Parse(nil, cStyle, nil)
	a.NotNil(ret).
		Equal(0, len(ret))

	ret = Parse([]byte(code1), cStyle, nil)
	a.NotNil(ret).
		Equal(1, len(ret)). // 字符串直接被过滤，不再返回
		True(strings.Contains(string(ret[4]), "注释代码"))

	// 注释缺少结束符
	//
	// 但依然会返回内容
	ret = Parse([]byte(code2), cStyle, nil)
	a.NotNil(ret).
		Equal(0, len(ret))
}

func TestMergeLines(t *testing.T) {
	a := assert.New(t)

	lines := [][]byte{
		[]byte("   l1\n"),
		[]byte("  l2\n"),
		[]byte("   l3"),
	}
	a.Equal(string(mergeLines(lines)), `l1
l2
 l3`)

	// 包含空格行
	lines = [][]byte{
		[]byte("   l1\n"),
		[]byte("    \n"),
		[]byte("  l2\n"),
		[]byte("   l3"),
	}
	a.Equal(string(mergeLines(lines)), `l1
    
l2
 l3`)

	// 包含空行
	lines = [][]byte{
		[]byte("   l1\n"),
		[]byte("\n"),
		[]byte("  l2\n"),
		[]byte("   l3"),
	}
	a.Equal(string(mergeLines(lines)), `l1

l2
 l3`)
}
