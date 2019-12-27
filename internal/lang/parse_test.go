// SPDX-License-Identifier: MIT

package lang

import (
	"strings"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/message/messagetest"
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

	erro, _, h := messagetest.MessageHandler()
	ret := Parse("", nil, nil, nil)
	a.NotNil(ret).
		Equal(0, len(ret))

	ret = Parse("", nil, cStyle, h)
	a.NotNil(ret).
		Equal(0, len(ret))
	h.Stop()
	a.Empty(erro.String())

	erro, _, h = messagetest.MessageHandler()
	ret = Parse("", []byte(code1), cStyle, h)
	a.NotNil(ret).
		Equal(1, len(ret)). // 字符串直接被过滤，不再返回
		True(strings.Contains(string(ret[4]), "注释代码"))
	h.Stop()
	a.Empty(erro.String())

	// 注释缺少结束符
	//
	// 但依然会返回内容
	erro, _, h = messagetest.MessageHandler()
	ret = Parse("", []byte(code2), cStyle, h)
	a.NotNil(ret).
		Equal(0, len(ret))
	a.Empty(erro.String())
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
