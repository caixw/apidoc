// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc 是一个 RESTful api 文档生成工具。
package main

import (
	"flag"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/output"
)

const usage = `%v 是一个 RESTful api 文档生成工具。

参数:
 -h       显示帮助信息；
 -v       显示版本信息；
 -l       显示所有支持的语言类型；
 -g       在当前目录下创建一个默认的配置文件。

源代码采用 MIT 开源许可证，发布于 %v
详细信息，可访问：%v
`

func main() {
	if flags() {
		return
	}

	start := time.Now() // 记录处理开始时间

	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}

	docs, err := input.Parse(cfg.Input)
	if err != nil {
		panic(err)
	}

	cfg.Output.Elapsed = time.Now().Sub(start)
	if err = output.Render(docs, cfg.Output); err != nil {
		panic(err)
	}
}

// 处理命令行参数，若被处理，返回 true，否则返回 false。
func flags() bool {
	flag.Usage = func() { fmt.Printf(usage, app.Name, app.RepoURL, app.OfficialURL) }
	h := flag.Bool("h", false, "显示帮助信息")
	v := flag.Bool("v", false, "显示版本信息")
	l := flag.Bool("l", false, "显示所有支持的语言")
	g := flag.Bool("g", false, "在当前目录下创建一个默认的配置文件")
	flag.Parse()

	switch {
	case *h:
		flag.Usage()
		return true
	case *v:
		fmt.Println("apidoc", app.Version, "build with", runtime.Version())
		return true
	case *l:
		langs := "[" + strings.Join(input.Langs(), ", ") + "]"
		fmt.Println("目前支持以下语言：", langs)
		return true
	case *g:
		if err := genConfigFile(); err != nil {
			panic(err)
		}
		return true
	}
	return false
}
