// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLangs(t *testing.T) {
	a := assert.New(t)

	list := Langs()
	a.Contains(list, "go", "php")
}

// 比较 langs 和 langExts 中的语言类型是否都一样。
func TestCompareLangsAndLangExts(t *testing.T) {
	a := assert.New(t)

	// 查询 langs 中的键名是否存在于 langExts
	for lang := range langs {
		exts, found := langExts[lang]
		a.True(found, "未找到与[%v]相对应的扩展名定义", lang).
			True(len(exts) > 0)
	}

	// 查询 langExts 中的键名是否存在于 langs
	for lang := range langExts {
		blocks, found := langs[lang]
		a.True(found, "未找到与[%v]相对应的代码块定义", lang).
			True(len(blocks) > 0)
	}
}

func TestDetectDirLang(t *testing.T) {
	a := assert.New(t)

	lang, err := DetectDirLang("./testdir")
	a.NotError(err).Equal(lang, "c++")

	lang, err = DetectDirLang("./testdir/testdir1")
	a.Error(err).Empty(lang)
}

func TestGetLangByExt(t *testing.T) {
	a := assert.New(t)

	a.Equal(getLangByExt(".C"), "c++")
	a.Equal(getLangByExt(".h"), "c++")
	a.Equal(getLangByExt(".c"), "c++")
	a.Equal(getLangByExt(".php"), "php")

	a.Equal(getLangByExt("php"), "")         // 扩展名不带.符号，查不到
	a.Equal(getLangByExt(".not exists"), "") // 真的不存在此扩展名
}
