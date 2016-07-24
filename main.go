// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc 是一个 RESTful API 文档生成工具。
package main

import (
	"bytes"
	"flag"
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
	"golang.org/x/text/language"
)

func main() {
	tag, err := locale.GetLocale()
	if err != nil {
		app.Error(err) // 输出错误信息，但不中断执行
	} else {
		app.Info("使用默认的本化语言：", app.DefaultLocale)
		tag, err = language.Parse(app.DefaultLocale)
		if err != nil {
			app.Error(err)
			return
		}
	}
	locale.SetLocale(tag)

	h := flag.Bool("h", false, locale.Sprintf(locale.FlagHUsage))
	v := flag.Bool("v", false, locale.Sprintf(locale.FlagVUsage))
	l := flag.Bool("l", false, locale.Sprintf(locale.FlagLUsage))
	g := flag.Bool("g", false, locale.Sprintf(locale.FlagGUsage))
	pprofType := flag.String("pprof", "", locale.Sprintf(locale.FlagPprofUsage))
	flag.Usage = usage
	flag.Parse()

	switch {
	case *h:
		flag.Usage()
		return
	case *v:
		locale.Fprintf(os.Stdout, locale.FlagVersionBuildWith, app.Name, app.Version, runtime.Version())
		return
	case *l:
		locale.Fprintf(os.Stdout, locale.FlagSupportedLangs, input.Langs())
		return
	case *g:
		path, err := getConfigFile()
		if err != nil {
			app.Error(err)
			return
		}
		if err = genConfigFile(path); err != nil {
			app.Error(err)
			return
		}
		app.Info(locale.Sprintf(locale.FlagConfigWritedSuccess, path))
		return
	}

	// 指定了 pprof 参数
	if len(*pprofType) > 0 {
		profile := filepath.Join("./", app.Profile)
		f, err := os.Create(profile)
		if err != nil {
			app.Error(err)
			return
		}
		defer func() {
			if err = f.Close(); err != nil {
				app.Error(err)
				return
			}
			app.Info(locale.Sprintf(locale.FlagPprofWritedSuccess, profile))
		}()

		switch strings.ToLower(*pprofType) {
		case "mem":
			defer func() {
				if err = pprof.Lookup("heap").WriteTo(f, 1); err != nil {
					app.Error(err)
				}
			}()
		case "cpu":
			if err := pprof.StartCPUProfile(f); err != nil {
				app.Error(err)
			}
			defer pprof.StopCPUProfile()
		default:
			app.Error(locale.Sprintf(locale.FlagInvalidPprrof))
			return
		}
	}

	run()
}

func usage() {
	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	flag.PrintDefaults()

	locale.Fprintf(os.Stdout, locale.FlagUsage, app.Name, buf.String(), app.RepoURL, app.OfficialURL)
}

// 真正的程序入口，main 主要是作参数的处理。
func run() {
	start := time.Now()

	path, err := getConfigFile()
	if err != nil {
		app.Error(err)
		return
	}

	cfg, err := loadConfig(path)
	if err != nil {
		app.Error(err)
		return
	}

	// 比较版本号兼容问题
	compatible, err := version.SemVerCompatible(app.Version, cfg.Version)
	if err != nil {
		app.Error(err)
		return
	}
	if !compatible {
		app.Error(locale.Sprintf(locale.VersionInCompatible))
		return
	}

	// 分析文档内容
	docs := doc.New()
	wg := &sync.WaitGroup{}
	for _, opt := range cfg.Inputs {
		wg.Add(1)
		go func(o *input.Options) {
			if err := input.Parse(docs, o); err != nil {
				app.Error(err)
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
		app.Error(err)
		return
	}

	app.Info(locale.Sprintf(locale.Complete, cfg.Output.Dir, time.Now().Sub(start)))
}

// 获取配置文件路径。目前只支持从工作路径获取。
func getConfigFile() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, app.ConfigFilename), nil
}
