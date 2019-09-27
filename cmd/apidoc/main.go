// SPDX-License-Identifier: MIT

// apidoc 是一个 RESTful API 文档生成工具。
package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/issue9/term/colors"
	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/internal/cmd/config"
	"github.com/caixw/apidoc/v5/internal/cmd/term"
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

func main() {
	if err := apidoc.Init(language.Und); err != nil {
		term.Line(warnOut, warnColor, err)
	}

	h := flag.Bool("h", false, locale.Sprintf(locale.FlagHUsage))
	v := flag.Bool("v", false, locale.Sprintf(locale.FlagVUsage))
	d := flag.Bool("d", false, locale.Sprintf(locale.FlagDUsage))
	l := flag.Bool("l", false, locale.Sprintf(locale.FlagLUsage))
	flag.Usage = usage
	flag.Parse()

	switch {
	case *h:
		flag.Usage()
		return
	case *v:
		term.Locale(infoOut, infoColor, locale.FlagVersionBuildWith, vars.Name, apidoc.Version(), runtime.Version())
		term.Locale(infoOut, infoColor, locale.FlagVersionCommitHash, vars.CommitHash())
		return
	case *l:
		for _, l := range term.Langs(3) {
			term.Line(infoOut, infoColor, l)
		}
		return
	case *d:
		write(getPaths())
		return
	}

	parse(getPaths())
}

func write(paths []string) {
	for _, dir := range paths {
		dir, err := filepath.Abs(dir)
		if err != nil {
			term.Line(erroOut, erroColor, err)
			return
		}

		if err := config.Write(dir); err != nil {
			term.Line(erroOut, erroColor, err)
			return
		}
		term.Locale(succOut, succColor, locale.ConfigWriteSuccess, dir)
	}
}

func parse(paths []string) {
	h := message.NewHandler(term.NewHandlerFunc(erroOut, warnOut, infoOut, succOut,
		erroColor, warnColor, infoColor, succColor))

	for _, path := range paths {
		now := time.Now()
		path, err := filepath.Abs(path)
		if err != nil {
			h.Error(message.Erro, err)
			return
		}

		cfg, err := config.Load(path)
		if err != nil {
			h.Error(message.Erro, err)
			break
		}

		if err := apidoc.Do(h, cfg.Output, cfg.Inputs...); err != nil {
			h.Error(message.Erro, err)
			break
		}

		elapsed := time.Now().Sub(now)
		h.Message(message.Succ, locale.Complete, cfg.Output.Path, elapsed)
	}

	h.Stop()
}

func usage() {
	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	flag.PrintDefaults()

	term.Locale(infoOut, infoColor, locale.FlagUsage, vars.Name, buf.String(), vars.RepoURL, vars.OfficialURL)
}

func getPaths() []string {
	paths := flag.Args()
	if len(paths) == 0 {
		paths = append(paths, "./")
	}
	return paths
}
