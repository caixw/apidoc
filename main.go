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

	i "github.com/caixw/apidoc/input"
	o "github.com/caixw/apidoc/output"
	"github.com/issue9/term/colors"
)

const (
	// 版本号
	version = "2.0.49.160529"

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

	err := run("./")
	if err != nil {
		panic(err)
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
		if err := genConfigFile(); err != nil {
			panic(err)
		}
		return true
	}
	return false
}

func printUsage() {
	colors.Print(colors.Stdout, colors.Default, colors.Default, usage)
}

func printLangs() {
	langs := "[" + strings.Join(i.Langs(), ", ") + "]"

	colors.Print(colors.Stdout, colors.Green, colors.Default, "目前支持以下语言：")
	colors.Println(colors.Stdout, colors.Default, colors.Default, langs)
}

func printVersion() {
	colors.Print(colors.Stdout, colors.Green, colors.Default, "apidoc ")
	colors.Println(colors.Stdout, colors.Default, colors.Default, version, "build with", runtime.Version())
}

func run(srcDir string) error {
	elapsed := time.Now()

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	docs, err := i.Parse(cfg.Input)
	if err != nil {
		return err
	}

	opt := &o.Options{
		Title:      cfg.Doc.Title,
		Version:    cfg.Doc.Version,
		DocDir:     cfg.Output.Dir,
		AppVersion: version,
		Elapsed:    time.Now().UnixNano() - elapsed.UnixNano(),
	}
	if err = o.Html(docs, opt); err != nil {
		return err
	}

	return nil
}
