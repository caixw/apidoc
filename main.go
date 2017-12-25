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

	"github.com/issue9/logs/writers"
	"github.com/issue9/term/colors"
	"github.com/issue9/version"
	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/locale"
	"github.com/caixw/apidoc/output"
	"github.com/caixw/apidoc/vars"
)

// 日志信息输出
var (
	info = newLog(os.Stdout, vars.InfoColor, "[INFO] ")
	warn = newLog(os.Stderr, vars.WarnColor, "[WARN] ")
	erro = newLog(os.Stderr, vars.ErroColor, "[ERRO] ")
)

func main() {
	initLocale() // 最先初始化本地化信息

	h := flag.Bool("h", false, locale.Sprintf(locale.FlagHUsage))
	v := flag.Bool("v", false, locale.Sprintf(locale.FlagVUsage))
	g := flag.Bool("g", false, locale.Sprintf(locale.FlagGUsage))
	wd := flag.String("wd", "./", locale.Sprintf(locale.FlagWDUsage))
	languages := flag.Bool("languages", false, locale.Sprintf(locale.FlagLanguagesUsage))
	encodings := flag.Bool("encodings", false, locale.Sprintf(locale.FlagEncodingsUsage))
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

	run(*wd)
}

// 真正的程序入口，main 主要是作参数的处理。
func run(wd string) {
	cfg, err := loadConfig(filepath.Join(wd, vars.ConfigFilename))
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
		return
	}

	docs, elapsed := input.Parse(cfg.Inputs...)

	cfg.Output.Elapsed = elapsed
	if err := output.Render(docs, cfg.Output); err != nil {
		erro.Println(err)
		return
	}

	info.Println(locale.Sprintf(locale.Complete, cfg.Output.Dir, elapsed))
}

func usage() {
	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	flag.PrintDefaults()

	locale.Printf(locale.FlagUsage, vars.Name, buf.String(), vars.RepoURL, vars.OfficialURL)
}

// 初始化本地化信息，确定在第一时间调用。
func initLocale() {
	if err := locale.Init(); err != nil {
		warn.Println(err)
	}

	// 本地化环境初始化成功之后，再设置日志前缀
	info.SetPrefix(locale.Sprintf(locale.InfoPrefix))
	warn.SetPrefix(locale.Sprintf(locale.WarnPrefix))
	erro.SetPrefix(locale.Sprintf(locale.ErrorPrefix))
}

// 生成一个默认的配置文件，并写入到文件中。
func genConfigFile(wd string) {
	o, err := input.Detect(wd, true)
	if err != nil {
		erro.Println(err)
		return
	}

	cfg := &config{
		Version: vars.Version(),
		Inputs:  []*input.Options{o},
		Output: &output.Options{
			Dir: filepath.Join(o.Dir, "doc"),
		},
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		erro.Println(err)
		return
	}

	path := filepath.Join(wd, vars.ConfigFilename)
	if err = ioutil.WriteFile(path, data, os.ModePerm); err != nil {
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
