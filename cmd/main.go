// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc 是一个 RESTful API 文档生成工具。
package main

import (
	"bytes"
	"flag"
	"github.com/caixw/apidoc/internal/errors"
	"path/filepath"
	"runtime"

	"golang.org/x/text/language"

	"github.com/caixw/apidoc"
	"github.com/caixw/apidoc/internal/lang"
	"github.com/caixw/apidoc/internal/locale"
	"github.com/caixw/apidoc/internal/output"
	"github.com/caixw/apidoc/internal/vars"
)

// 确保第一时间初始化本地化信息
func init() {
	if err := apidoc.InitLocale(language.Und); err != nil {
		warn.Println(err)
		return
	}

	initLogsLocale()
}

func main() {
	h := flag.Bool("h", false, locale.Sprintf(locale.FlagHUsage))
	v := flag.Bool("v", false, locale.Sprintf(locale.FlagVUsage))
	g := flag.Bool("g", false, locale.Sprintf(locale.FlagGUsage))
	wd := flag.String("wd", "./", locale.Sprintf(locale.FlagWDUsage))
	languages := flag.Bool("languages", false, locale.Sprintf(locale.FlagLanguagesUsage))
	encodings := flag.Bool("encodings", false, locale.Sprintf(locale.FlagEncodingsUsage))
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
		locale.Printf(locale.FlagSupportedLanguages, lang.Langs())
		return
	case *encodings:
		locale.Printf(locale.FlagSupportedEncodings)
		return
	case *g:
		genConfigFile(*wd)
		return
	}

	parse(*wd)
}

func parse(wd string) {
	cfg, err := loadConfig(filepath.Join(wd, configFilename))
	if err != nil {
		erro.Println(err)
		return
	}

	doc, err := apidoc.Parse(erro, cfg.Inputs...)
	if err != nil {
		erro.Println(err)
		return
	}

	if err = output.Render(doc, cfg.Output); err != nil {
		if ferr, ok := err.(*errors.Error); ok {
			ferr.File = configFilename
			ferr.Field = "output." + ferr.Field
		}
		erro.Println(err)
	}

	info.Println(locale.Sprintf(locale.Complete, cfg.Output.Path, doc.Elapsed))
}

func usage() {
	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	flag.PrintDefaults()

	locale.Printf(locale.FlagUsage, vars.Name, buf.String(), vars.RepoURL, vars.OfficialURL)
}

// 根据 wd 所在目录的内容生成一个配置文件，并写入到 wd 目录下的 .apidoc.yaml 中
func genConfigFile(wd string) {
	path := filepath.Join(wd, configFilename)
	if err := generateConfig(wd, path); err != nil {
		erro.Println(err)
		return
	}

	info.Println(locale.Sprintf(locale.FlagConfigWritedSuccess, path))
}

func printVersion() {
	locale.Printf(locale.FlagVersionBuildWith, vars.Name, vars.Version(), runtime.Version())
	locale.Printf(locale.FlagVersionCommitHash, vars.CommitHash())
}
