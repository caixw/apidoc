// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc 是一个 RESTful api 文档生成工具。
//
// 多行注释和单行注释在处理上会有一定区别：
//
// 单行注释，风格相同且相邻的注释会被合并成一个注释块。
// 单行注释，风格不相同且相邻的注释会被按注释风格多个注释块。
// 而多行注释，即使两个注释释块相邻也会被分成两个注释块来处理。
package main

import (
	"flag"
	"runtime"
	"strings"
	"time"

	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/logs"
	"github.com/caixw/apidoc/output"
)

const (
	// 版本号
	//
	// 版本号按照 http://semver.org/lang/zh-CN/ 中的规则，分成以下四个部分：
	// 主版本号.次版本号.修订号.修订日期
	version = "2.0.50.160529"

	// 配置文件名称。
	configFilename = ".apidoc.json"
)

const usage = `apidoc 是一个 RESTful api 文档生成工具。

参数:
 -h       显示当前帮助信息；
 -v       显示apidoc和go程序的版本信息；
 -l       显示所有支持的语言类型；
 -g       在当前目录下创建一个默认的配置文件；

有关 apidoc 的详细信息，可访问官网：http://apidoc.site`

func main() {
	if flags() {
		return
	}

	elapsed := time.Now()

	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}

	docs, err := input.Parse(cfg.Input)
	if err != nil {
		panic(err)
	}

	cfg.Output.AppVersion = version
	cfg.Output.Elapsed = time.Now().UnixNano() - elapsed.UnixNano()
	if err = output.Html(docs, cfg.Output); err != nil {
		panic(err)
	}
}

// 处理命令行参数，若被处理，返回true，否则返回false。
func flags() bool {
	flag.Usage = func() { logs.Println(usage) }
	h := flag.Bool("h", false, "显示帮助信息")
	v := flag.Bool("v", false, "显示帮助信息")
	l := flag.Bool("l", false, "显示所有支持的语言")
	g := flag.Bool("g", false, "在当前目录下创建一个默认的配置文件")
	flag.Parse()

	switch {
	case *h:
		flag.Usage()
		return true
	case *v:
		logs.Info("apidoc ", version, "build with", runtime.Version())
		return true
	case *l:
		langs := "[" + strings.Join(input.Langs(), ", ") + "]"
		logs.Info("目前支持以下语言：", langs)
		return true
	case *g:
		if err := genConfigFile(); err != nil {
			panic(err)
		}
		return true
	}
	return false
}
