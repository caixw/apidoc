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
	"runtime/pprof"
	"sync"
	"time"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/output"
	"github.com/issue9/version"
)

func main() {
	h := flag.Bool("h", false, "显示帮助信息")
	v := flag.Bool("v", false, "显示版本信息")
	l := flag.Bool("l", false, "显示所有支持的语言")
	g := flag.Bool("g", false, "在当前目录下创建一个默认的配置文件")
	pprofType := flag.String("pprof", "", "指定一种调试输出类型，可以为 cpu 或是 mem")
	flag.Usage = usage
	flag.Parse()

	switch {
	case *h:
		flag.Usage()
		return
	case *v:
		fmt.Fprintln(os.Stdout, app.Name, app.Version, "build with", runtime.Version())
		return
	case *l:
		fmt.Fprintln(os.Stdout, "目前支持以下语言", input.Langs())
		return
	case *g:
		path, err := getConfigFile()
		if err != nil {
			app.Errorln(err)
			return
		}
		if err = genConfigFile(path); err != nil {
			app.Errorln(err)
			return
		}
		app.Infoln("配置内容成功写入", path)
		return
	}

	if len(*pprofType) > 0 {
		profile := filepath.Join("./", app.Profile)
		f, err := os.Create(profile)
		if err != nil {
			app.Errorln(err)
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				app.Errorln(err)
				return
			}
			app.Infoln("pprof 的相关参数已经写入到", profile)
		}()

		switch *pprofType {
		case "mem":
			defer func() {
				if err = pprof.Lookup("heap").WriteTo(f, 1); err != nil {
					app.Errorln(err)
				}
			}()
		case "cpu":
			if err := pprof.StartCPUProfile(f); err != nil {
				app.Errorln(err)
			}
			defer pprof.StopCPUProfile()
		default:
			app.Errorln("无效的 pprof 参数")
			return
		}
	}

	run()
}

func usage() {
	fmt.Fprintln(os.Stdout, app.Name, "是一个 RESTful API 文档生成工具。")

	fmt.Fprintln(os.Stdout, "\n参数:")
	flag.CommandLine.SetOutput(os.Stdout)
	flag.PrintDefaults()

	fmt.Fprintln(os.Stdout, "\n源代码采用 MIT 开源许可证，发布于", app.RepoURL)
	fmt.Fprintln(os.Stdout, "详细信息可访问官网", app.OfficialURL)
}

// 真正的程序入口，main 主要是作为一个调试代码的处理。
// path 指定了配置文件的地址
func run() {
	start := time.Now() // 记录处理开始时间

	path, err := getConfigFile()
	if err != nil {
		app.Errorln(err)
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
	for _, opt := range cfg.Inputs {
		wg.Add(1)
		go func(o *input.Options) {
			if err := input.Parse(docs, o); err != nil {
				app.Errorln(err)
			}
			wg.Done()
		}(opt)
	}
	wg.Wait()

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

// 获取配置文件路径。目前只支持从工作路径获取。
func getConfigFile() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, app.ConfigFilename), nil
}
