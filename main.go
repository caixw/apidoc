// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc 是一个 RESTful api 文档生成工具。
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
	version = "2.2.56.160601"

	// 配置文件名称。
	configFilename = ".apidoc.json"
)

const usage = `apidoc 是一个 RESTful api 文档生成工具。

参数:
 -h       显示帮助信息；
 -v       显示版本信息；
 -l       显示所有支持的语言类型；
 -g       在当前目录下创建一个默认的配置文件。

有关 apidoc 的详细信息，可访问官网：http://apidoc.site`

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

	cfg.Output.AppVersion = version
	cfg.Output.Elapsed = time.Now().UnixNano() - start.UnixNano()
	if err = output.Render(docs, cfg.Output); err != nil {
		panic(err)
	}
}

// 处理命令行参数，若被处理，返回 true，否则返回 false。
func flags() bool {
	flag.Usage = func() { logs.Println(usage) }
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
