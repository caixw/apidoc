// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc从代码注释中提取并生成api的文档。
package main

import (
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/caixw/apidoc/scanner"
	"github.com/issue9/term/colors"
)

const version = "0.1.0.150730"

var usage = `apidoc从代码注释中提取并生成api的文档。

命令行语法:
 apidoc [options] src doc

options:
 -h     显示当前帮助信息；
 -v     显示apidoc和go程序的版本信息；
 -langs 显示所有支持的语言类型。
 -r     是否搜索子目录，默认为true；
 -t     目标文件类型，支持的类型可以通过-langs来查看；
 -ext   需要分析的文件的扩展名，若不指定，则会根据-t参数自动生成相应的扩展名。
        若-t也未指定，则会根据src目录下的文件，自动判断-t的值。

src:
 源文件所在的目录。
doc:
 产生的文档保存的目录。


源代码采用MIT开源许可证，并发布于github:https://github.com/caixw/apidoc
`

func main() {
	var h bool
	var v bool
	var l bool
	var r bool
	var t string
	var ext string

	flag.Usage = func() {
		fmt.Println(usage)
	}
	flag.BoolVar(&h, "h", false, "显示帮助信息")
	flag.BoolVar(&v, "v", false, "显示帮助信息")
	flag.BoolVar(&l, "langs", false, "显示所有支持的语言")
	flag.BoolVar(&r, "r", true, "搜索子目录，默认为true")
	flag.StringVar(&t, "t", "", "指定源文件的类型，若不指定，系统会自行判断")
	flag.StringVar(&ext, "ext", "", "匹配的扩展名，若不指定，会根据-t的指定，自行判断")
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	if v {
		colors.Print(colors.Stdout, colors.Green, colors.Default, "apidoc: ")
		colors.Println(colors.Stdout, colors.Default, colors.Default, version)
		colors.Print(colors.Stdout, colors.Green, colors.Default, "Go: ")
		goVersion := runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH
		colors.Println(colors.Stdout, colors.Default, colors.Default, goVersion)
		return
	}

	if l {
		fmt.Println(scanner.Langs())
		return
	}

	if flag.NArg() != 2 {
		colors.Println(colors.Stderr, colors.Red, colors.Default, "请同时指定src和dest参数")
		return
	}

	var exts []string
	if len(ext) > 0 {
		exts = strings.Split(strings.TrimSpace(ext), ",")
	}

	tree, err := scanner.Scan(flag.Arg(0), r, t, exts)
	if err != nil {
		panic(err)
	}

	tree.Version = version
	if err = tree.OutputHtml(flag.Arg(1)); err != nil {
		panic(err)
	}
}
