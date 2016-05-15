// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"github.com/caixw/apidoc/core"
	"github.com/caixw/apidoc/scanner"
)

// 每种语言对应的扩展名及用分析函数。
var langs = map[string]*lang{
	"go": &lang{
		Exts: []string{".go"},
		scan: scanner.C,
	},

	"cpp": &lang{
		Exts: []string{".h", ".cpp", ".cxx", ".c"},
		scan: scanner.C,
	},

	"c": &lang{
		Exts: []string{".h", ".c"},
		scan: scanner.C,
	},

	"php": &lang{
		Exts: []string{".php"},
		scan: scanner.C,
	},

	"js": &lang{
		Exts: []string{".js"},
		scan: scanner.C,
	},
	"ruby": &lang{
		Exts: []string{".rb"},
		scan: scanner.Ruby,
	},
}

// 各扩展名对应的语言。
// 数据从上面的langs数据中分析获得。
var extsIndex = map[string]string{}

type lang struct {
	Exts []string // 扩展名列表
	scan core.ScanFunc
}

func init() {
	for k, lang := range langs {
		for _, ext := range lang.Exts {
			extsIndex[ext] = k
		}
	}
}

func Langs() map[string]*lang {
	return langs
}
