// SPDX-License-Identifier: MIT

// apidoc 是一个 RESTful API 文档生成工具
//
// 大致的使用方法为：
//  apidoc [options] [path]
// 具体的参数说明，可以使用 h 参数查看：
//  apidoc -h
// path 表示目录列表，多个目录使用空格分隔。
// 用于在 path 下生成配置文件或是从 path 目录加载配置文件。
package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/issue9/term/colors"
	"golang.org/x/text/language"
	xmessage "golang.org/x/text/message"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/internal/lang"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
)

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

func init() {
	if err := apidoc.Init(language.Und); err != nil {
		printLine(warnOut, warnColor, err)
	}
}

func main() {
	h := flag.Bool("h", false, locale.Sprintf(locale.FlagHUsage))
	v := flag.Bool("v", false, locale.Sprintf(locale.FlagVUsage))
	d := flag.Bool("d", false, locale.Sprintf(locale.FlagDUsage))
	t := flag.Bool("t", false, locale.Sprintf(locale.FlagTUsage))
	l := flag.Bool("l", false, locale.Sprintf(locale.FlagLUsage))
	flag.Usage = usage
	flag.Parse()

	switch {
	case *h:
		flag.Usage()
		return
	case *v:
		goVersion := strings.TrimLeft(runtime.Version(), "go")
		printLocale(infoOut, infoColor, locale.FlagVersion, apidoc.Version(), vars.DocVersion(), vars.CommitHash(), goVersion)
		return
	case *t:
		parse(true)
		return
	case *l:
		langs(infoOut, infoColor, 3)
		return
	case *d:
		detect()
		return
	}

	parse(false)
}

func detect() {
	paths, err := getPaths()
	if err != nil {
		printLine(erroOut, erroColor, err)
		return
	}

	for _, dir := range paths {
		if err := apidoc.Detect(dir, true); err != nil {
			printLine(erroOut, erroColor, err)
			return
		}
		printLocale(succOut, succColor, locale.ConfigWriteSuccess, dir)
	}
}

// 参数 test 表示是否只作语法检测，不输出内容。
func parse(test bool) {
	h := message.NewHandler(newHandlerFunc())

	paths, err := getPaths()
	if err != nil {
		h.Error(message.Erro, err)
		return
	}
	for _, path := range paths {
		apidoc.Make(h, path, test)
	} // end for paths

	h.Stop()
}

func langs(w io.Writer, color colors.Color, tail int) {
	ls := lang.Langs()
	langs := make([]*lang.Language, 1, len(ls)+1)
	langs[0] = &lang.Language{
		Name:        locale.Sprintf(locale.LangID),
		DisplayName: locale.Sprintf(locale.LangName),
		Exts:        []string{locale.Sprintf(locale.LangExts)},
	}
	langs = append(langs, ls...)

	// 计算各列的最大长度值
	var maxDisplay, maxName int
	for _, l := range langs {
		width := len(l.DisplayName)
		if width > maxDisplay {
			maxDisplay = width
		}
		width = len(l.Name)
		if width > maxName {
			maxName = width
		}
	}
	maxDisplay += tail
	maxName += tail

	for _, l := range langs {
		n := l.Name + strings.Repeat(" ", maxName-len(l.Name))
		d := l.DisplayName + strings.Repeat(" ", maxDisplay-len(l.DisplayName))
		printLine(w, color, n, d, strings.Join(l.Exts, " "))
	}
}

func usage() {
	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	flag.PrintDefaults()

	printLocale(infoOut, infoColor, locale.FlagUsage, vars.Name, buf.String(), vars.RepoURL, vars.OfficialURL)
}

func getPaths() ([]string, error) {
	paths := flag.Args()
	if len(paths) == 0 {
		paths = append(paths, "./")
	}

	for index, path := range paths {
		path, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}

		paths[index] = path
	}

	return paths, nil
}

func newHandlerFunc() message.HandlerFunc {
	erroPrefix := locale.Sprintf(locale.ErrorPrefix)
	warnPrefix := locale.Sprintf(locale.WarnPrefix)
	infoPrefix := locale.Sprintf(locale.InfoPrefix)
	succPrefix := locale.Sprintf(locale.SuccessPrefix)

	return func(msg *message.Message) {
		switch msg.Type {
		case message.Erro:
			printMessage(erroOut, erroColor, erroPrefix, msg.Message)
		case message.Warn:
			printMessage(warnOut, warnColor, warnPrefix, msg.Message)
		case message.Succ:
			printMessage(succOut, succColor, succPrefix, msg.Message)
		default: // message.Info 采用相同的值
			printMessage(infoOut, infoColor, infoPrefix, msg.Message)
		}
	}
}

func printMessage(out io.Writer, color colors.Color, prefix, msg string) {
	if _, err := colors.Fprint(out, color, colors.Default, prefix); err != nil {
		panic(err)
	}
	printLine(out, colors.Default, msg)
}

// 向控制台输出一行本地化的内容
func printLocale(out io.Writer, color colors.Color, key xmessage.Reference, v ...interface{}) {
	l := locale.Sprintf(key, v...)
	printLine(out, color, l)
}

// 向控制台输出一行内容
func printLine(out io.Writer, color colors.Color, v ...interface{}) {
	if _, err := colors.Fprintln(out, color, colors.Default, v...); err != nil {
		panic(err)
	}
}
