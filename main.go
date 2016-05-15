// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc 是一个 RESTful api 文档生成工具。
package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/caixw/apidoc/app"
	"github.com/issue9/term/colors"
)

const (
	out          = colors.Stdout
	titleColor   = colors.Green
	contentColor = colors.Default
	errorColor   = colors.Red
	warnColor    = colors.Cyan
)

const usage = `apidoc 是一个 RESTful api 文档生成工具。

参数:
 -h       显示当前帮助信息；
 -v       显示apidoc和go程序的版本信息；
 -l       显示所有支持的语言类型；
 -r       是否搜索子目录，默认为true；
 -g       在当前目录下创建一个默认的配置文件；

有关 apidoc 的详细信息，可访问官网：http://apidoc.site`

func main() {
	if flags() {
		return
	}

	err := app.Run("./")
	if err != nil {
		fmt.Println(err)
	}
}

// 处理命令行参数，若被处理，返回true，否则返回false。
func flags() (ok bool) {
	var h, v, l, g bool

	flag.Usage = printUsage
	flag.BoolVar(&h, "h", false, "显示帮助信息")
	flag.BoolVar(&v, "v", false, "显示帮助信息")
	flag.BoolVar(&l, "l", false, "显示所有支持的语言")
	flag.BoolVar(&g, "g", false, "在当前目录下创建一个默认的配置文件")
	flag.Parse()

	switch {
	case h:
		flag.Usage()
		return true
	case v:
		printVersion()
		return true
	case l:
		printLangs()
		return true
	case g:
		err := app.GenConfigFile()
		if err != nil {
			printError(err)
		}
		return true
	}
	return false
}

func printUsage() {
	colors.Println(out, contentColor, colors.Default, usage)
}

func printError(msg ...interface{}) {
	colors.Println(out, errorColor, colors.Default, msg...)
}

func printLangs() {
	colors.Println(out, titleColor, colors.Default, "目前支持以下类型的代码解析:")
	for k, v := range app.Langs() {
		colors.Print(out, titleColor, colors.Default, k, ":")
		colors.Println(out, contentColor, colors.Default, v.Exts)
	}
}

func printVersion() {
	fmt.Println("apidoc", app.Version, "build with", runtime.Version())
}
