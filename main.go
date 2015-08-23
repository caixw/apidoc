// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc是一个RESTful api文档生成工具。
package main

import (
	"flag"
	"runtime"
	"time"

	"github.com/caixw/apidoc/core"
	o "github.com/caixw/apidoc/output"
	"github.com/issue9/term/colors"
)

const version = "0.7.37.150823"

const (
	out          = colors.Stdout
	titleColor   = colors.Green
	contentColor = colors.Default
	errorColor   = colors.Red
	warnColor    = colors.Cyan
)

const usage = `apidoc是一个RESTful api文档生成工具。

命令行语法:
 apidoc [options]

options:
 -h       显示当前帮助信息；
 -v       显示apidoc和go程序的版本信息；
 -l       显示所有支持的语言类型；
 -r       是否搜索子目录，默认为true；
 -g       在当前目录下创建一个默认的配置文件；

有关apidoc的详细信息，可访问官网：https://caixw.github.io/apidoc`

func main() {
	if flags() {
		return
	}

	elapsed := time.Now()

	cfg, err := loadConfig()
	if err != nil {
		printError(err)
		return
	}
	paths, err := recursivePath(cfg)
	if err != nil {
		printError(err)
		return
	}

	docs, err := core.ScanFiles(paths, cfg.lang.scan)
	if err != nil {
		printError(err)
		return
	}
	if docs.HasError() { // 语法错误，并不中断程序
		printSyntaxErrors(docs.Errors())
	}

	opt := &o.Options{
		Title:      cfg.Doc.Title,
		Version:    cfg.Doc.Version,
		DocDir:     cfg.Output.Dir,
		AppVersion: version,
		Elapsed:    time.Now().UnixNano() - elapsed.UnixNano(),
	}
	if err = o.Html(docs.Items(), opt); err != nil {
		printError(err)
		return
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
		err := genConfigFile()
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
