// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"flag"
	"io"
	"os"

	"github.com/issue9/cmdopt"
	"github.com/issue9/term/colors"
	xmessage "golang.org/x/text/message"

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

var command = cmdopt.New(os.Stdout, flag.ContinueOnError, usage, func(name string) string {
	return locale.Sprintf(locale.CmdNotFound, name)
})

func init() {
	command.Help("help", locale.Sprintf(locale.CmdHelpUsage))
}

// Exec 执行程序
func Exec() {
	if err := command.Exec(os.Args[1:]); err != nil {
		printLine(erroOut, erroColor, err)
	}
}

func usage(w io.Writer) error {
	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	flag.PrintDefaults()

	printLocale(w, infoColor, locale.FlagUsage, vars.Name, buf.String(), vars.RepoURL, vars.OfficialURL)

	return nil
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
