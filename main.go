// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc从代码注释中提取并生成api的文档。
package main

import (
	"flag"
	"strings"
	"time"

	"github.com/caixw/apidoc/core"
	"github.com/caixw/apidoc/output"
)

const version = "0.6.32.150822"

func main() {
	var (
		h      bool
		v      bool
		langs  bool
		r      bool
		t      string
		ext    string
		docVer string
		title  string
	)

	flag.Usage = printUsage
	flag.BoolVar(&h, "h", false, "显示帮助信息")
	flag.BoolVar(&v, "v", false, "显示帮助信息")
	flag.BoolVar(&langs, "langs", false, "显示所有支持的语言")
	flag.BoolVar(&r, "r", true, "搜索子目录，默认为true")
	flag.StringVar(&t, "t", "", "指定源文件的类型，若不指定，系统会自行判断")
	flag.StringVar(&ext, "ext", "", "匹配的扩展名，若不指定，会根据-t的指定，自行判断")
	flag.StringVar(&docVer, "version", "", "指定文档版本号")
	flag.StringVar(&title, "title", "apidoc", "指定文档标题")
	flag.Parse()

	switch {
	case h:
		flag.Usage()
		return
	case v:
		printVersion()
		return
	case langs:
		printLangs()
		return
	case flag.NArg() != 2:
		printError("请同时指定src和dest参数")
		return
	}

	var exts []string
	if len(ext) > 0 {
		exts = strings.Split(strings.TrimSpace(ext), ",")
	}

	elapsed := time.Now()

	inputOpt := &Options{
		SrcDir:    flag.Arg(0),
		Recursive: r,
		Type:      t,
		Exts:      exts,
	}
	f, paths, err := getArgs(inputOpt)
	if err != nil {
		printError(err)
		return
	}

	docs, err := core.ScanFiles(paths, f)
	if err != nil {
		printError(err)
		return
	}
	if docs.HasError() { // 语法错误，并不中断程序
		printSyntaxErrors(docs.Errors())
	}

	opt := &output.Options{
		Title:      title,
		Version:    docVer,
		DocDir:     flag.Arg(1),
		AppVersion: version,
		Elapsed:    time.Now().UnixNano() - elapsed.UnixNano(),
	}
	if err = output.Html(docs.Items(), opt); err != nil {
		printError(err)
		return
	}
}
