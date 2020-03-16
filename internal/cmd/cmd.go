// SPDX-License-Identifier: MIT

// Package cmd 提供子命令的相关功能
package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/issue9/cmdopt"
	"github.com/issue9/term/colors"
	xmessage "golang.org/x/text/message"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/vars"
)

// 命令行输出的表格中，每一列为了对齐填补的空格数量。
const tail = 3

const (
	succColor = colors.Green
	infoColor = colors.Default
	warnColor = colors.Cyan
	erroColor = colors.Red
)

var (
	succOut io.Writer = os.Stdout
	infoOut io.Writer = os.Stdout
	warnOut io.Writer = os.Stderr
	erroOut io.Writer = os.Stderr
)

var command *cmdopt.CmdOpt

func init() {
	command = cmdopt.New(os.Stdout, flag.ContinueOnError, usage, func(name string) string {
		return locale.Sprintf(locale.CmdNotFound, name)
	})

	command.Help("help", buildUsage(locale.CmdHelpUsage))

	initBuild()
	initDetect()
	initLang()
	initLocale()
	initTest()
	initVersion()
	initMock()
	initStatic()
}

// Exec 执行程序
func Exec() {
	if err := command.Exec(os.Args[1:]); err != nil {
		panic(err)
	}
}

func usage(w io.Writer) error {
	cmds := strings.Join(command.Commands(), ",")
	msg := locale.Sprintf(locale.CmdUsage, vars.Name, cmds, vars.RepoURL, vars.OfficialURL)
	_, err := fmt.Fprintln(w, msg)
	return err
}

func newHandlerFunc() core.HandlerFunc {
	erroPrefix := locale.Sprintf(locale.ErrorPrefix)
	warnPrefix := locale.Sprintf(locale.WarnPrefix)
	infoPrefix := locale.Sprintf(locale.InfoPrefix)
	succPrefix := locale.Sprintf(locale.SuccessPrefix)

	return func(msg *core.Message) {
		switch msg.Type {
		case core.Erro:
			printMessage(erroOut, erroColor, erroPrefix, msg.Message)
		case core.Warn:
			printMessage(warnOut, warnColor, warnPrefix, msg.Message)
		case core.Succ:
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

	if _, err := fmt.Fprintln(out, msg); err != nil {
		panic(err)
	}
}

func buildUsage(key xmessage.Reference, v ...interface{}) cmdopt.DoFunc {
	return func(w io.Writer) error {
		_, err := fmt.Fprintln(w, locale.Sprintf(key, v...))
		return err
	}
}

func getFlagSetUsage(fs *flag.FlagSet) string {
	buf := new(bytes.Buffer)
	origin := fs.Output()
	fs.SetOutput(buf)
	fs.PrintDefaults()
	fs.SetOutput(origin)
	return buf.String()
}
