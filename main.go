// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc 是一个 RESTful API 文档生成工具。
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/output"
	"github.com/issue9/version"
)

func main() {
	start := time.Now() // 记录处理开始时间

	wd, err := os.Getwd()
	if err != nil {
		app.Errorln(err)
		return
	}
	path := filepath.Join(wd, app.ConfigFilename)

	if flags(path) {
		return
	}

	cfg, err := loadConfig(path)
	if err != nil {
		app.Errorln(err)
		return
	}

	// 比较版本号兼容问题
	compatible, err := version.SemVerCompatible(app.Version, cfg.Version)
	if err != nil {
		app.Errorln(err)
		return
	}
	if !compatible {
		app.Errorln("当前程序与配置文件中指定的版本号不兼容")
		return
	}

	// 分析文档内容
	docs := doc.New()
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	for _, opt := range cfg.Inputs {
		wg.Add(1)
		go func() {
			if err := input.Parse(docs, opt); err != nil {
				app.Errorln(err)
			}
			wg.Done()
		}()
	}
	if len(docs.Title) == 0 {
		docs.Title = app.DefaultTitle
	}

	// 输出内容
	cfg.Output.Elapsed = time.Now().Sub(start)
	if err := output.Render(docs, cfg.Output); err != nil {
		app.Errorln(err)
		return
	}

	app.Infoln("完成！文档保存在", cfg.Output.Dir, "总用时", time.Now().Sub(start))
}

// 处理命令行参数，若被处理，返回 true，否则返回 false。
// path 配置文件的路径。
func flags(path string) bool {
	out := os.Stdout

	h := flag.Bool("h", false, "显示帮助信息")
	v := flag.Bool("v", false, "显示版本信息")
	l := flag.Bool("l", false, "显示所有支持的语言")
	g := flag.Bool("g", false, "在当前目录下创建一个默认的配置文件")
	flag.Usage = func() {
		fmt.Fprintln(out, app.Name, "是一个 RESTful API 文档生成工具。")
		fmt.Fprintln(out, "\n参数:")
		flag.CommandLine.SetOutput(out)
		flag.PrintDefaults()
		fmt.Fprintln(out, "\n源代码采用 MIT 开源许可证，发布于", app.RepoURL)
		fmt.Fprintln(out, "详细信息可访问官网", app.OfficialURL)
	}
	flag.Parse()

	switch {
	case *h:
		flag.Usage()
		return true
	case *v:
		fmt.Fprintln(out, app.Name, app.Version, "build with", runtime.Version())
		return true
	case *l:
		fmt.Fprintln(out, "目前支持以下语言", input.Langs())
		return true
	case *g:
		if err := genConfigFile(path); err != nil {
			app.Errorln(err)
			return true
		}
		app.Infoln("配置内容成功写入", path)
		return true
	}
	return false
}
