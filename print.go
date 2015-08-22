// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"runtime"

	"github.com/issue9/term/colors"
)

const (
	out          = colors.Stdout
	titleColor   = colors.Green
	contentColor = colors.Default
	errorColor   = colors.Red
	warnColor    = colors.Cyan
)

const usage = `apidoc从代码注释中提取并生成api的文档。

命令行语法:
 apidoc [options] src doc

options:
 -h       显示当前帮助信息；
 -v       显示apidoc和go程序的版本信息；
 -langs   显示所有支持的语言类型。
 -r       是否搜索子目录，默认为true；
 -t       目标文件类型，支持的类型可以通过-langs来查看；
 -version 指定文档的版本号；
 -title   指定文档的标题；
 -ext     需要分析的文件的扩展名，若不指定，则会根据-t参数自动生成相应的扩展名。
          若-t也未指定，则会根据src目录下的文件，自动判断-t的值。

src:
 源文件所在的目录。
doc:
 产生的文档保存的目录。

源代码采用MIT开源许可证，并发布于github:https://github.com/caixw/apidoc`

func printUsage() {
	colors.Println(out, contentColor, colors.Default, usage)
}

func printError(msg ...interface{}) {
	colors.Println(out, errorColor, colors.Default, msg...)
}

func printSyntaxErrors(errs []error) {
	for _, v := range errs {
		colors.Println(out, warnColor, colors.Default, v)
	}
}

func printLangs() {
	colors.Println(out, titleColor, colors.Default, "目前支持以下类型的代码解析:")
	for k, v := range langs {
		colors.Print(out, titleColor, colors.Default, k, ":")
		colors.Println(out, contentColor, colors.Default, v.exts)
	}
}

func printVersion() {
	colors.Print(out, titleColor, colors.Default, "apidoc: ")
	colors.Println(out, contentColor, colors.Default, version)

	colors.Print(out, titleColor, colors.Default, "Go: ")
	goVersion := runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH
	colors.Println(out, contentColor, colors.Default, goVersion)
}
