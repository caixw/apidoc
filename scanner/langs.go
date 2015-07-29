// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

// 各编程语言相关的参数
type lang struct {
	exts []string // 扩展名列表
	scan scanFunc // 代码扫描函数，用于从代码中分离注释文本。
}

var langs = map[string]*lang{
	"go": &lang{
		exts: []string{".go"},
		scan: cstyle,
	},

	"cpp": &lang{
		exts: []string{".h", ".cpp", ".cxx", ".c"},
		scan: cstyle,
	},

	"c": &lang{
		exts: []string{".h", ".c"},
		scan: cstyle,
	},

	"php": &lang{
		exts: []string{".php"},
		scan: cstyle,
	},

	"js": &lang{
		exts: []string{".js"},
		scan: cstyle,
	},
}

// 支持的语言列表。
func Langs() (list []string) {
	for l, _ := range langs {
		list = append(list, l)
	}

	return
}

// 各扩展名对应的语言。
// 数据由init函数从上面的langs数据中分析获得。
var extsIndex = map[string]string{}

func init() {
	for k, lang := range langs {
		for _, ext := range lang.exts {
			extsIndex[ext] = k
		}
	}
}
