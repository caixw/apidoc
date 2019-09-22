// SPDX-License-Identifier: MIT

// apidoc 是一个 RESTful API 文档生成工具。
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/issue9/term/colors"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/internal/lang"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
)

// 控制台的输出颜色
const (
	infoColor = colors.Green
	warnColor = colors.Cyan
	erroColor = colors.Red
)

var (
	infoOut = os.Stdout
	warnOut = os.Stderr
	erroOut = os.Stderr
)

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
	h := message.NewHandler(newConsoleHandlerFunc())

	cfg, err := loadConfig(wd)
	if err != nil {
		h.Error(message.Erro, err)
		return
	}

	now := time.Now()
	if err := apidoc.Do(h, cfg.Output, cfg.Inputs...); err != nil {
		h.Error(message.Erro, err)
		return
	}
	elapsed := time.Now().Sub(now)

	h.Message(message.Info, locale.Complete, cfg.Output.Path, elapsed)
}

func usage() {
	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	flag.PrintDefaults()

	fmt.Fprintln(infoOut, locale.Sprintf(locale.FlagUsage, vars.Name, buf.String(), vars.RepoURL, vars.OfficialURL))
}

// 根据 wd 所在目录的内容生成一个配置文件，并写入到 wd 目录下的 .apidoc.yaml 中
func genConfigFile(wd string) {
	path := filepath.Join(wd, configFilename)
	if err := generateConfig(wd, path); err != nil {
		fmt.Fprintln(infoOut, err)
		return
	}

	fmt.Fprintln(infoOut, locale.Sprintf(locale.FlagConfigWritedSuccess, path))
}

func printVersion() {
	fmt.Fprintln(infoOut, locale.Sprintf(locale.FlagVersionBuildWith, vars.Name, vars.Version(), runtime.Version()))
	fmt.Fprintln(infoOut, locale.Sprintf(locale.FlagVersionCommitHash, vars.CommitHash()))
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
		fmt.Fprintln(infoOut, l.Name, n, l.DisplayName, d, strings.Join(l.Exts, ", "))
	}
}

func newConsoleHandlerFunc() message.HandlerFunc {
	return func(err *message.Message) {
		switch err.Type {
		case message.Erro:
			printMessage(erroOut, erroColor, locale.Sprintf(locale.ErrorPrefix), err.Message)
		case message.Warn:
			printMessage(warnOut, warnColor, locale.Sprintf(locale.WarnPrefix), err.Message)
		default: // message.Info 采用相同的值
			printMessage(infoOut, infoColor, locale.Sprintf(locale.InfoPrefix), err.Message)
		}
	}
}

func printMessage(out *os.File, color colors.Color, prefix, msg string) {
	colors.Fprint(out, color, colors.Default, prefix)
	colors.Fprintln(out, colors.Default, colors.Default, msg)
}
