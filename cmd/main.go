// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// apidoc 是一个 RESTful API 文档生成工具。
package main

import (
	"bytes"
	"context"
	"flag"
	"os"
	"path/filepath"
	"runtime"

	"github.com/issue9/term/colors"
	"golang.org/x/text/language"

	"github.com/caixw/apidoc"
	"github.com/caixw/apidoc/doc"
	"github.com/caixw/apidoc/errors"
	"github.com/caixw/apidoc/internal/lang"
	"github.com/caixw/apidoc/internal/locale"
	"github.com/caixw/apidoc/internal/output"
	"github.com/caixw/apidoc/internal/vars"
)

// 控制台的输出颜色
const (
	infoColor = colors.Green
	warnColor = colors.Cyan
	erroColor = colors.Red
)

// 确保第一时间初始化本地化信息
func init() {
	if err := apidoc.InitLocale(language.Und); err != nil {
		printWarn(err)
		return
	}
}

func main() {
	h := flag.Bool("h", false, locale.Sprintf(locale.FlagHUsage))
	v := flag.Bool("v", false, locale.Sprintf(locale.FlagVUsage))
	g := flag.Bool("g", false, locale.Sprintf(locale.FlagGUsage))
	wd := flag.String("wd", "./", locale.Sprintf(locale.FlagWDUsage))
	languages := flag.Bool("l", false, locale.Sprintf(locale.FlagLanguagesUsage))
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
	case *g:
		genConfigFile(*wd)
		return
	}

	parse(*wd)
}

func parse(wd string) {
	cfg, err := loadConfig(wd)
	if err != nil {
		printError(err)
		return
	}

	h := errors.NewHandler(newConsoleHandlerFunc())
	doc, err := doc.Parse(context.Background(), h, cfg.Inputs...)
	if err != nil {
		if ferr, ok := err.(*errors.Error); ok {
			ferr.File = configFilename
			ferr.Field = "inputs." + ferr.Field
		}
		printError(err)
		return
	}

	if err = output.Render(doc, cfg.Output); err != nil {
		if ferr, ok := err.(*errors.Error); ok {
			ferr.File = configFilename
			ferr.Field = "output." + ferr.Field
		}
		printError(err)
	}

	printInfo(locale.Sprintf(locale.Complete, cfg.Output.Path, doc.Elapsed))
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
		printError(err)
		return
	}

	printInfo(locale.Sprintf(locale.FlagConfigWritedSuccess, path))
}

func printVersion() {
	locale.Printf(locale.FlagVersionBuildWith, vars.Name, vars.Version(), runtime.Version())
	locale.Printf(locale.FlagVersionCommitHash, vars.CommitHash())
}

func newConsoleHandlerFunc() errors.HandlerFunc {
	return func(err *errors.Error) {
		switch err.Type {
		case errors.SyntaxError:
			printError(err)
		case errors.SyntaxWarn:
			printWarn(err)
		default:
			printError(err)
		}
	}
}

func print(out *os.File, prefix string, color colors.Color, val interface{}) {
	if out != os.Stderr && out != os.Stdout {
		panic("无效的 out 参数")
	}

	colors.Fprint(out, color, colors.Default, prefix)
	colors.Fprintln(out, colors.Default, colors.Default, val)
}

func printWarn(val interface{}) {
	print(os.Stderr, locale.Sprintf(locale.WarnPrefix), warnColor, val)
}

func printError(val interface{}) {
	print(os.Stderr, locale.Sprintf(locale.ErrorPrefix), erroColor, val)
}

func printInfo(val interface{}) {
	print(os.Stderr, locale.Sprintf(locale.InfoPrefix), infoColor, val)
}
