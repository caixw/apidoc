// SPDX-License-Identifier: MIT

// apidoc 是一个 RESTful API 文档生成工具。
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/issue9/term/colors"
	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/errors"
	"github.com/caixw/apidoc/v5/internal/lang"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/output"
	"github.com/caixw/apidoc/v5/internal/vars"
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
	l := flag.Bool("l", false, locale.Sprintf(locale.FlagLanguagesUsage))
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
		printLanguages()
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

	now := time.Now()
	h := errors.NewHandler(newConsoleHandlerFunc())
	doc, err := doc.Parse(context.Background(), h, cfg.Inputs...)
	if err != nil {
		if ferr, ok := err.(*errors.Error); ok {
			ferr.File = configFilename
		}
		printError(err)
		return
	}

	if err = output.Render(doc, cfg.Output); err != nil {
		if ferr, ok := err.(*errors.Error); ok {
			ferr.File = configFilename
		}
		printError(err)
		return
	}
	elapsed := time.Now().Sub(now)

	printInfo(locale.Sprintf(locale.Complete, cfg.Output.Path, elapsed))
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

// 将支持的语言内容以表格的形式输出
func printLanguages() {
	langs := lang.Langs()
	var maxDisplay, maxName int
	for _, l := range langs {
		if len(l.DisplayName) > maxDisplay {
			maxDisplay = len(l.DisplayName)
		}
		if len(l.Name) > maxName {
			maxName = len(l.Name)
		}
	}

	maxDisplay += 3
	maxName += 3

	for _, l := range langs {
		d := strings.Repeat(" ", maxDisplay-len(l.DisplayName))
		n := strings.Repeat(" ", maxName-len(l.Name))
		fmt.Println(l.Name, n, l.DisplayName, d, strings.Join(l.Exts, ", "))
	}
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

func printWarn(val interface{}) {
	println(os.Stderr, locale.Sprintf(locale.WarnPrefix), warnColor, val)
}

func printError(val interface{}) {
	println(os.Stderr, locale.Sprintf(locale.ErrorPrefix), erroColor, val)
}

func printInfo(val interface{}) {
	println(os.Stderr, locale.Sprintf(locale.InfoPrefix), infoColor, val)
}

func println(out *os.File, prefix string, color colors.Color, val interface{}) {
	if out != os.Stderr && out != os.Stdout {
		panic("无效的 out 参数")
	}

	colors.Fprint(out, color, colors.Default, prefix)
	colors.Fprintln(out, colors.Default, colors.Default, val)
}
