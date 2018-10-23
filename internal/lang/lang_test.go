// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lang

import (
	"strings"
	"testing"
	"unicode"

	"github.com/issue9/assert"
)

func TestLangs(t *testing.T) {
	a := assert.New(t)

	isLower := func(str string) bool {
		for _, r := range str {
			if unicode.IsUpper(r) {
				return false
			}
		}
		return true
	}

	// Langs() 返回的应该和 langs 有相同的长度
	a.Equal(len(Langs()), len(langs))

	for index, lang := range langs {
		a.NotEmpty(lang.Name, "语言名称不能为空，在 %d", index)
		a.True(isLower(lang.Name), "名称非小写 %s", lang.Name)

		// 检测 block
		a.NotEmpty(lang.Blocks, "blocks 不能为空，在 %s", lang.Name)

		for index, blk := range lang.Blocks {
			b, ok := blk.(*block)
			if !ok {
				continue
			}
			v := (b.Type == blockTypeString || b.Type == blockTypeMComment || b.Type == blockTypeSComment)
			a.True(v, "langs[%v].[%v].Type 值为非法值", lang.Name, index)
		}

		// 检测扩展名
		for _, ext := range lang.Exts {
			a.NotEmpty(ext, "空的扩展名在 %s", lang.Name).
				Equal(ext[0], '.', "扩展名 %s 必须以 . 开头在 %s", ext, lang.Name).
				Equal(strings.TrimSpace(ext), ext, "扩展名 %s 存在首尾空格", ext, lang.Name).
				True(isLower(ext), "非小写的扩展名 %s 在 %s", ext, lang.Name)
		}
	}
}

func TestGet(t *testing.T) {
	a := assert.New(t)

	l := Get("go")
	a.NotNil(l).
		Equal(l.Name, "go").
		Equal(l.Exts, []string{".go"})

	// 不比较大小写
	l = Get("Go")
	a.Nil(l)
}

func TestGetByExt(t *testing.T) {
	a := assert.New(t)

	l := GetByExt(".go")
	a.NotNil(l).Equal(l.Name, "go")

	l = GetByExt(".cxx")
	a.NotNil(l).Equal(l.Name, "c++")

	// 不存在，不以 . 开头
	l = GetByExt("go")
	a.Nil(l)

	// 不存在
	l = GetByExt(".not-exists")
	a.Nil(l)
}
