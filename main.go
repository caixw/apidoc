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
	"strings"
	"sync"
	"time"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/output"
	"github.com/issue9/version"
)

func main() {
	h := flag.Bool("h", false, app.Sprintf(locale.FlagHUsage))
	v := flag.Bool("v", false, app.Sprintf(locale.FlagVUsage))
	l := flag.Bool("l", false, app.Sprintf(locale.FlagLUsage))
	g := flag.Bool("g", false, app.Sprintf(locale.FlagGUsage))
	pprofType := flag.String("pprof", "", app.Sprintf(locale.FlagPprofUsage))
	flag.Usage = usage
	flag.Parse()

	switch {
	case *h:
		flag.Usage()
		return
	case *v:
		app.Fprintf(os.Stdout, locale.FlagVersionBuildWith, app.Name, app.Version, runtime.Version())
		return
	case *l:
		fmt.Fprintf(os.Stdout, locale.FlagSupportedLangs, input.Langs())
		return
	case *g:
		path, err := getConfigFile()
		if err != nil {
			app.Error().Println(err)
			return
		}
		if err = genConfigFile(path); err != nil {
			app.Error().Println(err)
			return
		}
		app.Info().Printf(locale.FlagConfigWritedSuccess, path)
		return
	}

	// 指定了 pprof 参数
	if len(*pprofType) > 0 {
		profile := filepath.Join("./", app.Profile)
		f, err := os.Create(profile)
		if err != nil {
			app.Error().Println(err)
			return
		}
		defer func() {
			if err = f.Close(); err != nil {
				app.Error().Println(err)
				return
			}
			app.Info().Printf(locale.FlagPprofWritedSuccess, profile)
		}()

		switch strings.ToLower(*pprofType) {
		case "mem":
			defer func() {
				if err = pprof.Lookup("heap").WriteTo(f, 1); err != nil {
					app.Error().Println(err)
				}
			}()
		case "cpu":
			if err := pprof.StartCPUProfile(f); err != nil {
				app.Error().Println(err)
			}
			defer pprof.StopCPUProfile()
		default:
			app.Error().Printf(locale.FlagInvalidPprrof)
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

// 真正的程序入口，main 主要是作参数的处理。
func run() {
	start := time.Now()

	path, err := getConfigFile()
	if err != nil {
		app.Error().Println(err)
		return
	}

	cfg, err := loadConfig(path)
	if err != nil {
		app.Error().Println(err)
		return
	}

	// 比较版本号兼容问题
	compatible, err := version.SemVerCompatible(app.Version, cfg.Version)
	if err != nil {
		app.Error().Println(err)
		return
	}
	if !compatible {
		app.Error().Printf(locale.VersionInCompatible)
		return
	}

	// 分析文档内容
	docs := doc.New()
	wg := &sync.WaitGroup{}
	for _, opt := range cfg.Inputs {
		wg.Add(1)
		go func(o *input.Options) {
			if err := input.Parse(docs, o); err != nil {
				app.Error().Println(err)
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
		app.Error().Println(err)
		return
	}

	app.Info().Printf(locale.Complete, cfg.Output.Dir, time.Now().Sub(start))
}

// 获取配置文件路径。目前只支持从工作路径获取。
func getConfigFile() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, app.ConfigFilename), nil
}
