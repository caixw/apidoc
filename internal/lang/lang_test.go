// SPDX-License-Identifier: MIT

package lang

import (
	"strings"
	"testing"
	"unicode"

	"github.com/issue9/assert/v3"
)

func TestLangs(t *testing.T) {
	a := assert.New(t, false)

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
		a.NotEmpty(lang.ID, "语言名称不能为空，在 %d", index)
		a.True(isLower(lang.ID), "名称非小写 %s", lang.ID)

		// 检测 block
		a.NotEmpty(lang.blocks, "blocks 不能为空，在 %s", lang.ID)

		// 检测扩展名
		for _, ext := range lang.Exts {
			a.NotEmpty(ext, "空的扩展名在 %s", lang.ID).
				Equal(ext[0], '.', "扩展名 %s 必须以 . 开头在 %s", ext, lang.ID).
				Equal(strings.TrimSpace(ext), ext, "扩展名 %s 存在首尾空格", ext, lang.ID).
				True(isLower(ext), "非小写的扩展名 %s 在 %s", ext, lang.ID)
		}
	}
}

func TestGet(t *testing.T) {
	a := assert.New(t, false)

	l := Get("go")
	a.NotNil(l).
		Equal(l.ID, "go").
		Equal(l.Exts, []string{".go"})

	// 不比较大小写
	l = Get("Go")
	a.Nil(l)
}

func TestGetByExt(t *testing.T) {
	a := assert.New(t, false)

	l := GetByExt(".go")
	a.NotNil(l).Equal(l.ID, "go")

	l = GetByExt(".cxx")
	a.NotNil(l).Equal(l.ID, "c++")

	// 不存在
	l = GetByExt(".not-exists")
	a.Nil(l)

	a.Panic(func() {
		GetByExt("")
	})

	a.Panic(func() {
		GetByExt("go")
	})
}
