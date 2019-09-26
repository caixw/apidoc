// SPDX-License-Identifier: MIT

// apidoc 是一个 RESTful API 文档生成工具。
package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/issue9/term/colors"
	xmessage "golang.org/x/text/message"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/internal/cmd/config"
	"github.com/caixw/apidoc/v5/internal/lang"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
)

// 控制台的输出颜色
const (
	succColor = colors.Green
	infoColor = colors.Default
	warnColor = colors.Cyan
	erroColor = colors.Red
)

var (
	succOut = os.Stdout
	infoOut = os.Stdout
	warnOut = os.Stderr
	erroOut = os.Stderr
)

func main() {
	h := flag.Bool("h", false, locale.Sprintf(locale.FlagHUsage))
	v := flag.Bool("v", false, locale.Sprintf(locale.FlagVUsage))
	d := flag.String("d", "./", locale.Sprintf(locale.FlagDUsage))
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
	case *d != "":
		dir, err := filepath.Abs(*d)
		if err != nil {
			pLine(erroOut, erroColor, err)
		}
		if err := config.Write(dir); err != nil {
			pLine(erroOut, erroColor, err)
		} else {
			pLocale(succOut, succColor, locale.ConfigWriteSuccess, dir)
		}
		return
	}

	paths := flag.Args()
	if len(paths) == 0 {
		paths = append(paths, "./")
	}
	parse(paths)
}

func parse(paths []string) {
	h := message.NewHandler(newConsoleHandlerFunc())

	for _, path := range paths {
		now := time.Now()

		cfg, err := config.Load(path)
		if err != nil {
			h.Error(message.Erro, err)
			return
		}

		if err := apidoc.Do(h, cfg.Output, cfg.Inputs...); err != nil {
			h.Error(message.Erro, err)
			return
		}
		elapsed := time.Now().Sub(now)

		h.Message(message.Info, locale.Complete, cfg.Output.Path, elapsed)
	}
}

func usage() {
	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	flag.PrintDefaults()

	pLocale(infoOut, infoColor, locale.FlagUsage, vars.Name, buf.String(), vars.RepoURL, vars.OfficialURL)
}

func printVersion() {
	pLocale(infoOut, infoColor, locale.FlagVersionBuildWith, vars.Name, vars.Version(), runtime.Version())
	pLocale(infoOut, infoColor, locale.FlagVersionCommitHash, vars.CommitHash())
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
		pLine(infoOut, infoColor, l.Name, n, l.DisplayName, d, strings.Join(l.Exts, ", "))
	}
}

// 向控制台输出一行本地化的内容
func pLocale(out io.Writer, color colors.Color, key xmessage.Reference, v ...interface{}) {
	l := locale.Sprintf(key, v...)
	pLine(out, color, l)
}

// 向控制台输出一行内容
func pLine(out io.Writer, color colors.Color, v ...interface{}) {
	if _, err := colors.Fprintln(out, color, colors.Default, v...); err != nil {
		panic(err)
	}
}

func newConsoleHandlerFunc() message.HandlerFunc {
	erroPrefix := locale.Sprintf(locale.ErrorPrefix)
	warnPrefix := locale.Sprintf(locale.WarnPrefix)
	infoPrefix := locale.Sprintf(locale.InfoPrefix)
	succPrefix := locale.Sprintf(locale.SuccessPrefix)

	return func(err *message.Message) {
		switch err.Type {
		case message.Erro:
			printMessage(erroOut, erroColor, erroPrefix, err.Message)
		case message.Warn:
			printMessage(warnOut, warnColor, warnPrefix, err.Message)
		case message.Succ:
			printMessage(succOut, succColor, succPrefix, err.Message)
		default: // message.Info 采用相同的值
			printMessage(infoOut, infoColor, infoPrefix, err.Message)
		}
	}
}

func printMessage(out *os.File, color colors.Color, prefix, msg string) {
	if _, err := colors.Fprint(out, color, colors.Default, prefix); err != nil {
		panic(err)
	}
	pLine(out, colors.Default, msg)
}
