// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc 是一个 RESTful API 文档生成工具。
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/output"
	"github.com/issue9/version"
)

const usage = `%v 是一个 RESTful API 文档生成工具。

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
		app.Error(err)
		return
	}

	// 比较版本号兼容问题
	compatible, err := version.SemVerCompatible(app.Version, cfg.Version)
	if err != nil {
		app.Error(err)
		return
	}
	if !compatible {
		app.Error("当前程序与配置文件中指定的版本号不兼容")
		return
	}

	docs, err := input.Parse(cfg.Input)
	if err != nil {
		app.Error(err)
		return
	}

	cfg.Output.Elapsed = time.Now().Sub(start)
	if err = output.Render(docs, cfg.Output); err != nil {
		app.Error(err)
		return
	}

	app.Info("编译完成，总用时：", time.Now().Sub(start))
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
		fmt.Fprintln(os.Stdout, app.Name, app.Version, "build with", runtime.Version())
		return true
	case *l:
		fmt.Fprintln(os.Stdout, "目前支持以下语言：", input.Langs())
		return true
	case *g:
		if err := genConfigFile(); err != nil {
			app.Error(err)
		}
		return true
	}
	return false
}
