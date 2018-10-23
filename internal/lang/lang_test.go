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
