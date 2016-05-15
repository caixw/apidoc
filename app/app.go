// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 加载程序的一些基本配置内容。
package app

import (
	"time"

	"github.com/caixw/apidoc/core"
	o "github.com/caixw/apidoc/output"
	"github.com/issue9/term/colors"
)

const (
	// 版本号
	Version = "1.0.42.160514"

	// 配置文件名称。
	configFilename = ".apidoc.json"
)

// 输出模块的常量定义
// TODO 去掉
const (
	out          = colors.Stdout
	titleColor   = colors.Green
	contentColor = colors.Default
	errorColor   = colors.Red
	warnColor    = colors.Cyan
)

func Run(srcDir string) error {
	elapsed := time.Now()

	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	paths, err := recursivePath(cfg)
	if err != nil {
		return err
	}

	docs, err := core.ScanFiles(paths, cfg.lang.scan)
	if err != nil {
		return err
	}
	if docs.HasError() { // 语法错误，并不中断程序
		printSyntaxErrors(docs.Errors())
	}

	opt := &o.Options{
		Title:      cfg.Doc.Title,
		Version:    cfg.Doc.Version,
		DocDir:     cfg.Output.Dir,
		AppVersion: Version,
		Elapsed:    time.Now().UnixNano() - elapsed.UnixNano(),
	}
	if err = o.Html(docs.Items(), opt); err != nil {
		return err
	}

	return nil
}

func printSyntaxErrors(errs []error) {
	for _, v := range errs {
		colors.Println(out, warnColor, colors.Default, v)
	}
}
