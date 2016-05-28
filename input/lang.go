// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

// 所有支持的语言模型定义
var langs = map[string][]*block{
	// golang
	"go": []*block{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "`", End: "`"},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},

	// c/c++
	"c": []*block{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},

	// php
	"php": []*block{
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: `\`},
		&block{Type: blockTypeString, Begin: "'", End: "'"},
		&block{Type: blockTypeSComment, Begin: `//`},
		&block{Type: blockTypeMComment, Begin: `/*`, End: `*/`},
	},
}

// 返回所有支持的语言
func Langs() []string {
	ret := make([]string, 0, len(langs))
	for l := range langs {
		ret = append(ret, l)
	}

	return ret
}

// 是否支持该语言
func LangIsSupported(lang string) bool {
	_, found := langs[lang]
	return found
}
