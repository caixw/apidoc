// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc 是一个 RESTful API 文档生成工具。
package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/output"
	"github.com/caixw/apidoc/types"
	"github.com/caixw/apidoc/vars"

	"github.com/issue9/logs/writers"
	"github.com/issue9/term/colors"
	"github.com/issue9/version"
	yaml "gopkg.in/yaml.v2"
)

// 日志信息输出
var (
	info = log.New(writers.NewConsole(os.Stdout, colors.Green, colors.Default), "[INFO]", 0)
	warn = log.New(writers.NewConsole(os.Stderr, colors.Cyan, colors.Default), "[WARN]", 0)
	erro = log.New(writers.NewConsole(os.Stderr, colors.Red, colors.Default), "[ERRO]", 0)
)

func main() {
	if err := locale.Init(); err != nil {
		warn.Println(err)
	}

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
		printVersion()
		return
	case *l:
		locale.Printf(locale.FlagSupportedLangs, input.Langs())
		return
	case *g:
		genConfigFile()
		return
	}

	/* 对 pprof 的处理，pprof 需要运行程序，所以注意关闭文件的时间。 */
	if len(*pprofType) > 0 {
		profile := filepath.Join("./", *pprofType+".prof")
		f, err := os.Create(profile)
		if err != nil { // 不能创建文件，则忽略 pprof 相关操作
			erro.Println(err)
			return
		}

		defer func() {
			if err = f.Close(); err != nil {
				erro.Println(err)
				return
			}
			info.Println(locale.Sprintf(locale.FlagPprofWritedSuccess, profile))
		}()

		switch strings.ToLower(*pprofType) {
		case "mem":
			defer func() {
				if err = pprof.Lookup("heap").WriteTo(f, 1); err != nil {
					erro.Println(err)
				}
			}()
		case "cpu":
			if err := pprof.StartCPUProfile(f); err != nil {
				erro.Println(err)
			}
			defer pprof.StopCPUProfile()
		default:
			erro.Println(locale.Sprintf(locale.FlagInvalidPprrof))
			return
		}
	}

	run()
}

// 真正的程序入口，main 主要是作参数的处理。
func run() {
	start := time.Now()

	path, err := getConfigFile()
	if err != nil {
		erro.Println(err)
		return
	}

	cfg, err := loadConfig(path)
	if err != nil {
		erro.Println(err)
		return
	}

	// 比较版本号兼容问题
	compatible, err := version.SemVerCompatible(vars.Version(), cfg.Version)
	if err != nil {
		erro.Println(err)
		return
	}
	if !compatible {
		erro.Println(locale.Sprintf(locale.VersionInCompatible))
	}

	docs := parse(cfg.Inputs)

	// 输出内容
	cfg.Output.Elapsed = time.Now().Sub(start)
	if err := output.Render(docs, cfg.Output); err != nil {
		erro.Println(err)
		return
	}

	info.Println(locale.Sprintf(locale.Complete, cfg.Output.Dir, time.Now().Sub(start)))
}

func parse(inputs []*input.Options) *types.Doc {
	docs := types.NewDoc()

	wg := &sync.WaitGroup{}
	for _, opt := range inputs {
		wg.Add(1)
		go func(o *input.Options) {
			if err := input.Parse(docs, o); err != nil {
				erro.Println(err)
			}
			wg.Done()
		}(opt)
	}
	wg.Wait()

	if len(docs.Title) == 0 {
		docs.Title = vars.DefaultTitle
	}

	return docs
}

func usage() {
	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	flag.PrintDefaults()

	locale.Printf(locale.FlagUsage, vars.Name, buf.String(), vars.RepoURL, vars.OfficialURL)
}

// 生成一个默认的配置文件，并写入到文件中。
func genConfigFile() {
	path, err := getConfigFile()
	if err != nil {
		erro.Println(err)
		return
	}

	dir := filepath.Dir(path)
	lang, err := input.DetectDirLang(dir)
	if err != nil { // 不中断，仅作提示用。
		warn.Println(err)
	}

	cfg := &config{
		Version: vars.Version(),
		Inputs: []*input.Options{
			&input.Options{
				Dir:       dir,
				Recursive: true,
				Lang:      lang,
			},
		},
		Output: &output.Options{
			Dir: filepath.Join(dir, "doc"),
		},
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		erro.Println(err)
		return
	}

	fi, err := os.Create(path)
	if err != nil {
		erro.Println(err)
		return
	}
	defer fi.Close()

	if _, err = fi.Write(data); err != nil {
		erro.Println(err)
		return
	}

	info.Println(locale.Sprintf(locale.FlagConfigWritedSuccess, path))
}

func printVersion() {
	locale.Printf(locale.FlagVersionBuildWith, vars.Name, vars.Version(), runtime.Version())
	locale.Printf(locale.FlagVersionCommitHash, vars.CommitHash())
}
