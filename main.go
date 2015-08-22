// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc从代码注释中提取并生成api的文档。
package main

import (
	"flag"
	"time"

	"github.com/caixw/apidoc/core"
	o "github.com/caixw/apidoc/output"
)

const version = "0.7.33.150822"

func main() {
	var (
		h     bool
		v     bool
		langs bool
		r     bool
	)

	flag.Usage = printUsage
	flag.BoolVar(&h, "h", false, "显示帮助信息")
	flag.BoolVar(&v, "v", false, "显示帮助信息")
	flag.BoolVar(&langs, "langs", false, "显示所有支持的语言")
	flag.BoolVar(&r, "r", true, "搜索子目录，默认为true")
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
	}

	elapsed := time.Now()

	cfg, err := loadConfig()
	if err != nil {
		printError(err)
		return
	}
	paths, err := recursivePath(cfg.Input.Dir, cfg.Input.Recursive, cfg.Input.Exts...)
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
