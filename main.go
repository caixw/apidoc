// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc 是一个 RESTful API 文档生成工具。
package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/issue9/logs/writers"
	"github.com/issue9/term/colors"

	"github.com/caixw/apidoc/config"
	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/output"
	"github.com/caixw/apidoc/parser"
	"github.com/caixw/apidoc/vars"
)

// 日志信息输出
var (
	info = newLog(os.Stdout, vars.InfoColor, "[INFO] ")
	warn = newLog(os.Stderr, vars.WarnColor, "[WARN] ")
	erro = newLog(os.Stderr, vars.ErroColor, "[ERRO] ")
)

// 确保第一时间初始化本地化信息
func init() {
	if err := locale.Init(); err != nil {
		warn.Println(err)
		return
	}

	info.SetPrefix(locale.Sprintf(locale.InfoPrefix))
	warn.SetPrefix(locale.Sprintf(locale.WarnPrefix))
	erro.SetPrefix(locale.Sprintf(locale.ErrorPrefix))
}

func main() {
	h := flag.Bool("h", false, locale.Sprintf(locale.FlagHUsage))
	v := flag.Bool("v", false, locale.Sprintf(locale.FlagVUsage))
	g := flag.Bool("g", false, locale.Sprintf(locale.FlagGUsage))
	wd := flag.String("wd", "./", locale.Sprintf(locale.FlagWDUsage))
	languages := flag.Bool("languages", false, locale.Sprintf(locale.FlagLanguagesUsage))
	encodings := flag.Bool("encodings", false, locale.Sprintf(locale.FlagEncodingsUsage))
	pprofType := flag.String("pprof", "", locale.Sprintf(locale.FlagPprofUsage, vars.PprofCPU, vars.PprofMem))
	flag.Usage = usage
	flag.Parse()

	switch {
	case *h:
		flag.Usage()
		return
	case *v:
		printVersion()
		return
	case *languages:
		locale.Printf(locale.FlagSupportedLanguages, input.Languages())
		return
	case *encodings:
		locale.Printf(locale.FlagSupportedEncodings, input.Encodings())
		return
	case *g:
		genConfigFile(*wd)
		return
	}

	if len(*pprofType) > 0 {
		buf := new(bytes.Buffer)
		defer func() { // 在程序结束时，将内容写入到文件
			profile := filepath.Join(*wd, *pprofType+".prof")
			if err := ioutil.WriteFile(profile, buf.Bytes(), os.ModePerm); err != nil {
				erro.Println(err)
			}
		}()

		switch strings.ToLower(*pprofType) {
		case vars.PprofMem:
			defer func() {
				if err := pprof.Lookup("heap").WriteTo(buf, 1); err != nil {
					erro.Println(err)
				}
			}()
		case vars.PprofCPU:
			if err := pprof.StartCPUProfile(buf); err != nil {
				erro.Println(err)
			}

			defer pprof.StopCPUProfile()
		default:
			erro.Println(locale.Sprintf(locale.FlagInvalidPprrof))
			return
		}
	}

	parse(*wd)
}

func parse(wd string) {
	cfg, err := config.Load(filepath.Join(wd, vars.ConfigFilename))
	if err != nil {
		erro.Println(err)
		return
	}

	start := time.Now()
	docs, err := parser.Parse(erro, warn, cfg.Inputs...)
	if err != nil {
		erro.Println(err)
	}
	cfg.Output.Elapsed = time.Now().Sub(start)
	if err := output.Render(docs, cfg.Output); err != nil {
		erro.Println(err)
		return
	}

	info.Println(locale.Sprintf(locale.Complete, cfg.Output.Dir, cfg.Output.Elapsed))
}

func usage() {
	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	flag.PrintDefaults()

	locale.Printf(locale.FlagUsage, vars.Name, buf.String(), vars.RepoURL, vars.OfficialURL)
}

// 根据 wd 所在目录的内容生成一个配置文件，并写入到 wd 目录下的 .apidoc.yaml 中
func genConfigFile(wd string) {
	path := filepath.Join(wd, vars.ConfigFilename)
	if err := config.Generate(wd, path); err != nil {
		erro.Println(err)
		return
	}

	info.Println(locale.Sprintf(locale.FlagConfigWritedSuccess, path))
}

func printVersion() {
	locale.Printf(locale.FlagVersionBuildWith, vars.Name, vars.Version(), runtime.Version())
	locale.Printf(locale.FlagVersionCommitHash, vars.CommitHash())
}

func newLog(out *os.File, color colors.Color, prefix string) *log.Logger {
	return log.New(writers.NewConsole(out, color, colors.Default), prefix, 0)
}
